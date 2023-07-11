package login

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/kvizdos/locksmith/database"
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
		w.WriteHeader(http.StatusNotAcceptable)
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

	var loginReq loginRequest
	err = json.Unmarshal(body, &loginReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !loginReq.HasRequiredFields() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": loginReq.Username,
	})

	if !usernameExists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var tmpUser users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

	passwordValidated, err := user.ValidatePassword(loginReq.Password)

	if err != nil {
		fmt.Println("Error validating user password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !passwordValidated {
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

	cookieValue := user.GenerateCookieValueFromSession(session)

	cookie := http.Cookie{Name: "token", Value: cookieValue, Expires: time.Unix(session.ExpiresAt, 0), HttpOnly: true, Secure: true, Path: "/"}
	http.SetCookie(w, &cookie)
}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("login.html").Parse(string(pages.LoginPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
