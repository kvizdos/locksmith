package administration

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

type AdministrationListUsersHandler struct{}

func (h AdministrationListUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	var lsu users.LocksmithUser
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

func ServeAdminPage(w http.ResponseWriter, r *http.Request) {
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
