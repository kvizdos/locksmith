package register

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"text/template"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/users"
	// "kv.codes/locksmith/database"
)

type registrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

func (r registrationRequest) HasRequiredFields() bool {
	return !(r.Username == "" || r.Password == "" || r.Email == "")
}

type RegisterCustomUserFunc func(users.LocksmithUser) users.LocksmithUserInterface

type RegistrationHandler struct {
	DefaultRoleName           string
	DisablePublicRegistration bool
	ConfigureCustomUser       RegisterCustomUserFunc
}

func (rr RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if rr.DefaultRoleName == "" {
		fmt.Println("Registration role name must be set!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !roles.RoleExists(rr.DefaultRoleName) {
		fmt.Println("Registration role name is invalid!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var registrationReq registrationRequest
	err = json.Unmarshal(body, &registrationReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if rr.DisablePublicRegistration && len(registrationReq.Code) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !registrationReq.HasRequiredFields() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pattern := "^[a-zA-Z0-9]+$"
	onlyContainsAlphanumericalCharacters, _ := regexp.MatchString(pattern, registrationReq.Username)

	if !onlyContainsAlphanumericalCharacters {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emailPattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	isValidemail, _ := regexp.MatchString(emailPattern, registrationReq.Email)

	if !isValidemail {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	useRole := rr.DefaultRoleName

	var invite invitations.Invitation
	if len(registrationReq.Code) > 0 {
		if len(registrationReq.Code) != 96 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		invite, err = invitations.GetInviteFromCode(db, registrationReq.Code)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if invite.Email != registrationReq.Email {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		useRole = invite.Role
	}

	usernameAndEmailCheck, _ := db.Find("users", map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"username": registrationReq.Username,
			},
			{
				"email": registrationReq.Email,
			},
		},
	})

	if len(usernameAndEmailCheck) != 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}

	password, err := authentication.CompileLocksmithPassword(registrationReq.Password)

	if err != nil {
		fmt.Println("Error compiling password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var lsu users.LocksmithUserInterface
	lsu = users.LocksmithUser{
		ID:               uuid.New().String(),
		Username:         registrationReq.Username,
		Email:            registrationReq.Email,
		PasswordInfo:     password,
		WebAuthnSessions: []webauthn.SessionData{},
		PasswordSessions: []authentication.PasswordSession{},
		Role:             useRole,
	}

	if rr.ConfigureCustomUser != nil {
		lsu = rr.ConfigureCustomUser(lsu.(users.LocksmithUser))
	}

	_, err = db.InsertOne("users", lsu.ToMap())

	if err != nil {
		fmt.Println("Error adding new user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(registrationReq.Code) > 0 {
		invite.Expire(db)
	}

	w.WriteHeader(http.StatusOK)
}

type RegistrationPageHandler struct {
	AppName string
	// Only allow users with an invite code to register
	DisablePublicRegistration bool
	Styling                   pages.LocksmithPageStyling
}

func (rr RegistrationPageHandler) servePublicHTML(w http.ResponseWriter, r *http.Request, invite ...invitations.Invitation) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("register.html").Parse(string(pages.RegisterPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type TemplateData struct {
		HasInvite  bool
		Invitation invitations.Invitation
		Title      string
		Styling    pages.LocksmithPageStyling
	}
	inv := TemplateData{
		Title:   rr.AppName,
		Styling: rr.Styling,
	}

	if inv.Styling.SubmitColor == "" {
		inv.Styling.SubmitColor = "#476ade"
	}

	if inv.Styling.StartGradient == "" {
		inv.Styling.StartGradient = "#476ade"
	}

	if inv.Styling.EndGradient == "" {
		inv.Styling.EndGradient = "#2744a3"
	}

	if inv.Title == "" {
		inv.Title = "Locksmith"
	}

	if len(invite) > 0 {
		inv.HasInvite = true
		inv.Invitation = invite[0]
	}

	err = tmpl.Execute(w, inv)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func (rr RegistrationPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("database").(database.DatabaseAccessor)

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	inviteCode := params.Get("invite")

	if rr.DisablePublicRegistration && len(inviteCode) == 0 {
		w.Write([]byte("public registrations are not allowed."))
		return
	}

	if inviteCode != "" {
		invite, err := invitations.GetInviteFromCode(db, inviteCode)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("invalid invitation code."))
			return
		}

		invite.Code = inviteCode

		rr.servePublicHTML(w, r, invite)
	} else {
		rr.servePublicHTML(w, r)
	}
}
