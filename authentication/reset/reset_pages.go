package reset

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/pages"
)

type ResetPasswordPageHandler struct {
	AppName               string
	Styling               pages.LocksmithPageStyling
	EmailAsUsername       bool
	ShowResetStage        bool
	HIBP                  hibp.HIBPSettings
	MinimumPasswordLength int
}

func (h ResetPasswordPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("reset-password.html").Parse(string(pages.ResetPasswordPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type TemplateData struct {
		Title                 string
		Styling               pages.LocksmithPageStyling
		EmailAsUsername       bool
		ShowResetStage        bool
		PasswordSecurityLink  string
		MinimumPasswordLength int
		HIBPEnforcement       hibp.HIBPEnforcement
	}
	inv := TemplateData{
		Title:                 h.AppName,
		Styling:               h.Styling,
		EmailAsUsername:       h.EmailAsUsername,
		ShowResetStage:        h.ShowResetStage,
		PasswordSecurityLink:  h.HIBP.PasswordSecurityInfoLink,
		MinimumPasswordLength: h.MinimumPasswordLength,
		HIBPEnforcement:       h.HIBP.Enforcement,
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
