package register

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/google/uuid"
	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
	// "kv.codes/locksmith/database"
)

type registrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r registrationRequest) HasRequiredFields() bool {
	return !(r.Username == "" || r.Password == "" || r.Email == "")
}

type RegistrationHandler struct {
	DefaultRoleName string
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

	_, usernameOrEmailTaken := db.Find("users", map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"username": registrationReq.Username,
			},
			{
				"email": registrationReq.Email,
			},
		},
	})

	if usernameOrEmailTaken {
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
		"role":        rr.DefaultRoleName,
	})

	if err != nil {
		fmt.Println("Error adding new user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fp := filepath.Join("pages", "register.html")

	tmpl, err := template.New("register.html").ParseFiles(fp)

	if err != nil {
		fmt.Println(err)
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
