package httphandlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	captchaproviders "github.com/kvizdos/locksmith/captcha-providers"
	"github.com/kvizdos/locksmith/pages"
)

type LoginPageMiddleware struct {
	Next http.Handler
}

func (h LoginPageMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Next.ServeHTTP(w, r)
}

type LoginPageHandler struct {
	AppName                   string
	DisablePublicRegistration bool
	Styling                   pages.LocksmithPageStyling
	EmailAsUsername           bool
	OnboardingPath            string
	CaptchaProvider           captchaproviders.CAPTCHAProvider
	OAuthProviders            []string
}

func (lr LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("login.html").Parse(string(pages.LoginPageHTML))
	if err != nil {
		fmt.Println(err)
	}

	type PageData struct {
		Title                     string
		Styling                   pages.LocksmithPageStyling
		EmailAsUsername           bool
		OnboardingPath            string
		OAuthProviders            string
		CaptchaProvider           captchaproviders.CAPTCHAProvider
		DisablePublicRegistration bool
	}

	providers := ""
	if lr.OAuthProviders != nil {
		js, _ := json.Marshal(lr.OAuthProviders)
		providers = string(js)
	}

	data := PageData{
		Title:                     lr.AppName,
		Styling:                   lr.Styling,
		EmailAsUsername:           lr.EmailAsUsername,
		OnboardingPath:            lr.OnboardingPath,
		CaptchaProvider:           lr.CaptchaProvider,
		OAuthProviders:            providers,
		DisablePublicRegistration: lr.DisablePublicRegistration,
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

	if err = tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
	}
}
