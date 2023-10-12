package login

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/users"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r loginRequest) HasRequiredFields() bool {
	return !(r.Username == "" || r.Password == "")
}

type LoginHandler struct{}

func (lh LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		logger.LOGGER.Log(logger.INVALID_METHOD, logger.GetIPFromRequest(*r), r.URL.Path, "POST", r.Method)
		w.WriteHeader(http.StatusNotAcceptable)
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

	var loginReq loginRequest
	err = json.Unmarshal(body, &loginReq)

	if err != nil {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !loginReq.HasRequiredFields() {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": strings.ToLower(loginReq.Username),
	})

	if !usernameExists {
		logger.LOGGER.Log(logger.LOGIN_INVALID_USERNAME, logger.GetIPFromRequest(*r), loginReq.Username)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var tmpUser users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

	passwordValidated, err := user.ValidatePassword(loginReq.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !passwordValidated {
		logger.LOGGER.Log(logger.LOGIN_FAIL_BAD_PASSWORD, loginReq.Username, logger.GetIPFromRequest(*r))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := user.GeneratePasswordSession()

	if err != nil {
		fmt.Println("Error generating session token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		fmt.Println("Error marshalling session token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = user.SavePasswordSession(session, db)

	if err != nil {
		fmt.Println("Error saving session token to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.LOGGER.Log(logger.LOGIN, loginReq.Username, logger.GetIPFromRequest(*r))

	cookieValue := user.GenerateCookieValueFromSession(session)

	cookie := http.Cookie{Name: "token", Value: cookieValue, Expires: time.Unix(session.ExpiresAt, 0), HttpOnly: true, Secure: true, Path: "/"}
	http.SetCookie(w, &cookie)
}

type LoginPageHandler struct {
	AppName string
	// Only allow users with an invite code to register
	DisablePublicRegistration bool
	Styling                   pages.LocksmithPageStyling
	EmailAsUsername           bool
	OnboardingPath            string
}

func (lr LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("login.html").Parse(string(pages.LoginPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type PageData struct {
		Title           string
		Styling         pages.LocksmithPageStyling
		EmailAsUsername bool
		OnboardingPath  string
	}

	data := PageData{
		Title:           lr.AppName,
		Styling:         lr.Styling,
		EmailAsUsername: lr.EmailAsUsername,
		OnboardingPath:  lr.OnboardingPath,
	}

	if data.Styling.SubmitColor == "" {
		data.Styling.SubmitColor = "#476ade"
	}

	if data.Styling.StartGradient == "" {
		data.Styling.StartGradient = "#476ade"
	}

	if data.Styling.EndGradient == "" {
		data.Styling.EndGradient = "#2744a3"
	}

	if data.Title == "" {
		data.Title = "Locksmith"
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
