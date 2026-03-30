package verificationcodes

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/users"
)

type VerificationPageHandler struct {
	AppName string
	Styling pages.LocksmithPageStyling
}

func (rr VerificationPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !r.URL.Query().Has("code") {
		rr.serveMain(w, r)
		return
	}
	rr.serveExchange(w, r)
}

func (rr VerificationPageHandler) serveMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("verification.html").Parse(string(pages.VerificationPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type TemplateData struct {
		Title   string
		Styling pages.LocksmithPageStyling
		Email   string
	}

	authUser, _ := r.Context().Value("authUser").(users.LocksmithUser)

	if !authUser.RequiresEmailVerification() {
		http.Redirect(w, r, "/app", http.StatusSeeOther)
		return
	}

	inv := TemplateData{
		Title:   rr.AppName,
		Styling: rr.Styling,
		Email:   authUser.GetEmail(),
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
		log.Println("Error executing template :", err)
		return
	}
}

func (rr VerificationPageHandler) serveExchange(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("verification-exchange.html").Parse(string(pages.VerificationExchangePageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type TemplateData struct {
		Title   string
		Styling pages.LocksmithPageStyling
		Code    string
	}

	authUser, _ := r.Context().Value("authUser").(users.LocksmithUser)

	if !authUser.RequiresEmailVerification() {
		http.Redirect(w, r, "/app", http.StatusSeeOther)
		return
	}

	inv := TemplateData{
		Title:   rr.AppName,
		Styling: rr.Styling,
		Code:    r.URL.Query().Get("code"),
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
		log.Println("Error executing template :", err)
		return
	}
}
