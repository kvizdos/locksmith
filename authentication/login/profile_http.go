package login

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kvizdos/locksmith/pages"
)

type ProfileHTTP struct {
	AppName string
	Styling pages.LocksmithPageStyling
}

func (lr ProfileHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("profile.html").Parse(string(pages.ProfilePageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type PageData struct {
		Title   string
		Styling pages.LocksmithPageStyling
	}

	data := PageData{
		Title:   lr.AppName,
		Styling: lr.Styling,
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
