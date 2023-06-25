package administration

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

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
