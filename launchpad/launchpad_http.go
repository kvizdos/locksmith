//go:build enable_launchpad
// +build enable_launchpad

package launchpad

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kvizdos/locksmith/pages"
)

type LaunchpadHTTPHandler struct {
	AppName        string
	Subtitle       string
	Styling        pages.LocksmithPageStyling
	AccessToken    string
	AvailableUsers map[string]LocksmithLaunchpadUserOptions
}

func (lr LaunchpadHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("launchpad.html").Parse(string(pages.LaunchpadPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type PageData struct {
		Title          string
		Styling        pages.LocksmithPageStyling
		AvailableUsers map[string]LocksmithLaunchpadUserOptions
		AccessToken    string
		Subtitle       string
	}

	data := PageData{
		Title:          lr.AppName,
		Styling:        lr.Styling,
		AvailableUsers: lr.AvailableUsers,
		AccessToken:    lr.AccessToken,
		Subtitle: lr.Subtitle,
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
