package register

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication/oauth"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/pages"
)

type RegistrationPageHandler struct {
	AppName string
	// Only allow users with an invite code to register
	DisablePublicRegistration bool
	Styling                   pages.LocksmithPageStyling
	EmailAsUsername           bool
	HasOnboarding             bool
	InviteUsedRedirect        string
	MinimumLengthRequirement  int
	OAuthProviders            oauth.OAuthProviders
}

func (rr RegistrationPageHandler) servePublicHTML(w http.ResponseWriter, _ *http.Request, invite ...invitations.Invitation) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("register.html").Parse(string(pages.RegisterPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	providers := ""

	if rr.OAuthProviders != nil {
		js, _ := json.Marshal(rr.OAuthProviders.GetNames())
		providers = string(js)
	}
	type TemplateData struct {
		HasInvite             bool
		Invitation            invitations.Invitation
		Title                 string
		Styling               pages.LocksmithPageStyling
		EmailAsUsername       bool
		HasOnboarding         bool
		MinimumPasswordLength int
		OAuthProviders        string
	}
	inv := TemplateData{
		Title:                 rr.AppName,
		Styling:               rr.Styling,
		EmailAsUsername:       rr.EmailAsUsername,
		HasOnboarding:         rr.HasOnboarding,
		MinimumPasswordLength: rr.MinimumLengthRequirement,
		OAuthProviders:        providers,
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

	if len(invite) > 0 {
		inv.HasInvite = true
		inv.Invitation = invite[0]
	}

	err = tmpl.Execute(w, inv)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func (rr RegistrationPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("database").(database.DatabaseAccessor)

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	inviteCode := params.Get("invite")

	if rr.DisablePublicRegistration && len(inviteCode) == 0 {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		// w.Write([]byte("public registrations are not allowed."))
		return
	}

	if inviteCode != "" {
		invite, err := invitations.GetInviteFromCode(db, inviteCode)

		if err != nil {
			logger.LOGGER.Log(logger.INVITE_CODE_FAKE_VIEW, logger.GetIPFromRequest(*r), inviteCode)
			http.Redirect(w, r, rr.InviteUsedRedirect, http.StatusTemporaryRedirect)
			return
		}

		logger.LOGGER.Log(logger.INVITE_CODE_LOADED, logger.GetIPFromRequest(*r), inviteCode, invite.AttachUserID)

		invite.Code = inviteCode

		rr.servePublicHTML(w, r, invite)
	} else {
		rr.servePublicHTML(w, r)
	}
}
