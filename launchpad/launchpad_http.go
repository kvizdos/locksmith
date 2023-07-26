//go:build enable_launchpad
// +build enable_launchpad

package launchpad

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/pages"
)

type LaunchpadHTTPHandler struct {
	AppName           string
	Subtitle          string
	Styling           pages.LocksmithPageStyling
	AccessToken       string
	AvailableUsers    map[string]LocksmithLaunchpadUserOptions
	RefreshButtonText string
	// Setting this to TRUE will make the
	// launchpad buttons RED to notify
	// users that it is an EARLY preview
	// before it hits an official staging
	// environment.
	IsEarlyDevelopmentEnvironment bool
}

func (lr LaunchpadHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("launchpad.html").Parse(string(pages.LaunchpadPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type PageData struct {
		Title             string
		Styling           pages.LocksmithPageStyling
		AvailableUsers    map[string]LocksmithLaunchpadUserOptions
		AccessToken       string
		Subtitle          string
		RefreshButtonText string
	}

	data := PageData{
		Title:             lr.AppName,
		Styling:           lr.Styling,
		AvailableUsers:    lr.AvailableUsers,
		AccessToken:       lr.AccessToken,
		Subtitle:          lr.Subtitle,
		RefreshButtonText: lr.RefreshButtonText,
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

	if lr.IsEarlyDevelopmentEnvironment {
		data.Styling = pages.LocksmithPageStyling{
			StartGradient: "#d9452b",
			EndGradient:   "#d9452b",
			SubmitColor:   "#d9452b",
		}
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

type LaunchpadRefreshHTTPHandler struct {
	BootstrapDatabase func(db database.DatabaseAccessor)
}

func (lr LaunchpadRefreshHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	database := r.Context().Value("database").(database.DatabaseAccessor)

	lr.BootstrapDatabase(database)

	w.WriteHeader(http.StatusOK)
	return
}
