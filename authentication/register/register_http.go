package register

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/google/uuid"
	"kv.codes/locksmith/administration/invitations"
	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
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

type RegistrationHandler struct {
	DefaultRoleName           string
	DisablePublicRegistration bool
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
			fmt.Println("here")
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

	_, err = db.InsertOne("users", map[string]interface{}{
		"id":          uuid.New().String(),
		"username":    registrationReq.Username,
		"password":    password,
		"email":       registrationReq.Email,
		"sessions":    []interface{}{},
		"websessions": []interface{}{},
		"role":        useRole,
	})

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
	// Only allow users with an invite code to register
	DisablePublicRegistration bool
}

func (rr RegistrationPageHandler) servePublicHTML(w http.ResponseWriter, r *http.Request, invite ...invitations.Invitation) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fp := filepath.Join("pages", "register.html")

	tmpl, err := template.New("register.html").ParseFiles(fp)

	if err != nil {
		fmt.Println(err)
	}

	type TemplateData struct {
		HasInvite  bool
		Invitation invitations.Invitation
	}
	inv := TemplateData{}

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
