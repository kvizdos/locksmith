package reset

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/kvizdos/locksmith/pages"
)

type ResetPasswordPageHandler struct {
	AppName         string
	Styling         pages.LocksmithPageStyling
	EmailAsUsername bool
	ShowResetStage  bool
}

func (h ResetPasswordPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("reset-password.html").Parse(string(pages.ResetPasswordPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type TemplateData struct {
		Title           string
		Styling         pages.LocksmithPageStyling
		EmailAsUsername bool
		ShowResetStage  bool
	}
	inv := TemplateData{
		Title:           h.AppName,
		Styling:         h.Styling,
		EmailAsUsername: h.EmailAsUsername,
		ShowResetStage:  h.ShowResetStage,
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

	err = tmpl.Execute(w, inv)

	if err != nil {
		fmt.Println("Error executing template :", err)
		return
	}
}
