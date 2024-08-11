package register

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/logger"
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
	PwnOK    bool   `json:"pwnok,omitempty"`
}

func (r registrationRequest) HasRequiredFields() bool {
	return !(r.Username == "" || r.Password == "")
}

type RegisterCustomUserFunc func(users.LocksmithUser, database.DatabaseAccessor) users.LocksmithUserInterface

type RegistrationHandler struct {
	DefaultRoleName                string
	DisablePublicRegistration      bool
	ConfigureCustomUser            RegisterCustomUserFunc
	EmailAsUsername                bool
	MinimumLengthRequirement       int
	HIBP                           hibp.HIBPSettings
	DefaultRegistrationEntitlement string
	NewRegistrationEvent           func(users.LocksmithUserInterface)
}

type registrationResponse struct {
	Error     string `json:"error,omitempty"`
	PwnStatus bool   `json:"pwned,omitempty"`
}

func (r registrationResponse) Marshal() []byte {
	js, _ := json.Marshal(r)
	return js
}

func (r *registrationResponse) Unmarshal(err []byte) {
	json.Unmarshal(err, r)
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
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var registrationReq registrationRequest
	err = json.Unmarshal(body, &registrationReq)

	if err != nil {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{
			Error: "could not unmarshal",
		}.Marshal())
		return
	}

	if rr.DisablePublicRegistration && len(registrationReq.Code) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if rr.EmailAsUsername {
		registrationReq.Email = registrationReq.Username
	}

	if !registrationReq.HasRequiredFields() {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{
			Error: "missing fields",
		}.Marshal())
		return
	}

	// Start HIBP Check
	hibpIsPwnedChan := make(chan bool)
	if rr.HIBP.Enabled && !(rr.HIBP.Enforcement == hibp.LOOSE && registrationReq.PwnOK) {
		httpClient := &http.Client{}
		if rr.HIBP.HTTPClient != nil {
			httpClient = rr.HIBP.HTTPClient
		}

		go hibp.CheckPassword(rr.HIBP.AppName, registrationReq.Password, hibpIsPwnedChan, httpClient)
	}

	// Confirm Password Length Requirements
	if rr.MinimumLengthRequirement != 0 && rr.MinimumLengthRequirement > len(registrationReq.Password) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{
			Error: "password too short",
		}.Marshal())
		return
	}

	// Default restrictive username
	pattern := "^[a-zA-Z0-9]+$"

	// Allow slightly more open if emails-as-username
	if rr.EmailAsUsername {
		pattern = `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	}

	onlyContainsAlphanumericalCharacters, _ := regexp.MatchString(pattern, registrationReq.Username)

	if !onlyContainsAlphanumericalCharacters {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{
			Error: "illegal username characters",
		}.Marshal())
		return
	}

	emailPattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	isValidemail, _ := regexp.MatchString(emailPattern, registrationReq.Email)

	if !isValidemail {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{
			Error: "invalid email",
		}.Marshal())
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	useID := uuid.New().String()

	var invite invitations.Invitation
	useRole := rr.DefaultRoleName
	if len(registrationReq.Code) > 0 {
		if len(registrationReq.Code) != 96 {
			logger.LOGGER.Log(logger.INVITE_CODE_MALFORMED, logger.GetIPFromRequest(*r), registrationReq.Code)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(registrationResponse{
				Error: "bad invite code",
			}.Marshal())
			return
		}

		invite, err = invitations.GetInviteFromCode(db, registrationReq.Code)

		if err != nil {
			logger.LOGGER.Log(logger.INVITE_CODE_FAKE, logger.GetIPFromRequest(*r), registrationReq.Code)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(registrationResponse{
				Error: "invalid code",
			}.Marshal())
			return
		}

		if invite.Email != registrationReq.Email {
			logger.LOGGER.Log(logger.INVITE_CODE_INCORRECT_EMAIL, logger.GetIPFromRequest(*r), registrationReq.Code)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(registrationResponse{
				Error: "invalid email",
			}.Marshal())
			return
		}

		useRole = invite.Role
		useID = invite.AttachUserID
	}

	usernameAndEmailCheck, _ := db.Find("users", map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"username": strings.ToLower(registrationReq.Username),
			},
			{
				"email": strings.ToLower(registrationReq.Email),
			},
		},
	})

	if len(usernameAndEmailCheck) != 0 {
		w.WriteHeader(http.StatusConflict)
		w.Write(registrationResponse{
			Error: "taken",
		}.Marshal())
		return
	}

	password, err := authentication.CompileLocksmithPassword(registrationReq.Password)

	if err != nil {
		fmt.Println("Error compiling password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Confirm HIBP stuff
	if rr.HIBP.Enabled && !(rr.HIBP.Enforcement == hibp.LOOSE && registrationReq.PwnOK) {
		passwordIsPwned := <-hibpIsPwnedChan
		if passwordIsPwned && (rr.HIBP.Enforcement == hibp.STRICT || (rr.HIBP.Enforcement == hibp.LOOSE && !registrationReq.PwnOK)) {
			w.WriteHeader(http.StatusConflict)
			w.Write(registrationResponse{
				Error:     "password pwned",
				PwnStatus: true,
			}.Marshal())
			return
		}
	}

	var lsu users.LocksmithUserInterface
	lsu = users.LocksmithUser{
		ID:               useID,
		Username:         strings.ToLower(registrationReq.Username),
		Email:            strings.ToLower(registrationReq.Email),
		PasswordInfo:     password,
		WebAuthnSessions: []webauthn.SessionData{},
		PasswordSessions: []authentication.PasswordSession{},
		Roles:            []string{useRole},
	}

	if rr.ConfigureCustomUser != nil {
		lsu = rr.ConfigureCustomUser(lsu.(users.LocksmithUser), db)
	}

	_, err = db.InsertOne("users", lsu.ToMap())

	if err != nil {
		fmt.Println("Error adding new user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(registrationReq.Code) > 0 {
		logger.LOGGER.Log(logger.INVITE_CODE_USED, logger.GetIPFromRequest(*r), registrationReq.Username, registrationReq.Code)
		invite.Expire(db)
	} else {
		logger.LOGGER.Log(logger.REGISTRATION_SUCCESS, logger.GetIPFromRequest(*r), registrationReq.Username)
	}

	if rr.NewRegistrationEvent != nil {
		go rr.NewRegistrationEvent(lsu)
	}

	w.WriteHeader(http.StatusOK)
}

type RegistrationPageHandler struct {
	AppName string
	// Only allow users with an invite code to register
	DisablePublicRegistration bool
	Styling                   pages.LocksmithPageStyling
	EmailAsUsername           bool
	HasOnboarding             bool
	InviteUsedRedirect        string
	HIBPIntegrationOptions    hibp.HIBPSettings
	MinimumLengthRequirement  int
}

func (rr RegistrationPageHandler) servePublicHTML(w http.ResponseWriter, r *http.Request, invite ...invitations.Invitation) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("register.html").Parse(string(pages.RegisterPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type TemplateData struct {
		HasInvite             bool
		Invitation            invitations.Invitation
		Title                 string
		Styling               pages.LocksmithPageStyling
		EmailAsUsername       bool
		HasOnboarding         bool
		HIBPEnforcement       hibp.HIBPEnforcement
		MinimumPasswordLength int
		PasswordSecurityLink  string
	}
	inv := TemplateData{
		Title:                 rr.AppName,
		Styling:               rr.Styling,
		EmailAsUsername:       rr.EmailAsUsername,
		HasOnboarding:         rr.HasOnboarding,
		HIBPEnforcement:       rr.HIBPIntegrationOptions.Enforcement,
		MinimumPasswordLength: rr.MinimumLengthRequirement,
		PasswordSecurityLink:  rr.HIBPIntegrationOptions.PasswordSecurityInfoLink,
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
			logger.LOGGER.Log(logger.INVITE_CODE_FAKE_VIEW, logger.GetIPFromRequest(*r), inviteCode)
			http.Redirect(w, r, rr.InviteUsedRedirect, http.StatusTemporaryRedirect)
			return
		}

		logger.LOGGER.Log(logger.INVITE_CODE_LOADED, logger.GetIPFromRequest(*r), inviteCode, invite.AttachUserID)

		invite.Code = inviteCode

		rr.servePublicHTML(w, r, invite)
	} else {
		rr.servePublicHTML(w, r)
	}
}
