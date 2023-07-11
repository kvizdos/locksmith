package administration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type AdministrationListUsersHandler struct {
	UserInterface users.LocksmithUserInterface
}

func (h AdministrationListUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.UserInterface == nil {
		h.UserInterface = users.LocksmithUser{}
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	lsu := reflect.Zero(reflect.TypeOf(h.UserInterface)).Interface().(users.LocksmithUserInterface)
	users, err := ListUsers(db, lsu)

	if err != nil {
		fmt.Println("failed to serve listing:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsond, err := json.Marshal(users)

	if err != nil {
		fmt.Println("failed to marshal user json:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsond)
}

type AdministrationDeleteUsersHandler struct {
	UserInterface users.LocksmithUserInterface
}

func (h AdministrationDeleteUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type deleteRequest struct {
		Username string `json:"username"`
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("authUser").(users.LocksmithUser)

	if !ok {
		fmt.Println("Delete users endpoint is required to be wrapped in SecureEndpointHTTPMiddleware()")
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

	var delReq deleteRequest
	err = json.Unmarshal(body, &delReq)

	if err != nil || (err == nil && delReq.Username == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := user.GetRole()

	if err != nil {
		// handle the error
		fmt.Println("Error parsing role:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if delReq.Username != user.GetUsername() {
		if !role.HasPermission("user.delete.other") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else if !role.HasPermission("user.delete.self") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	deleted, err := DeleteUser(db, delReq.Username)

	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type ServeAdminPage struct{}

func (p ServeAdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fp := filepath.Join("pages", "admin.html")

	tmpl, err := template.New("admin.html").ParseFiles(fp)

	if err != nil {
		fmt.Println(err)
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
