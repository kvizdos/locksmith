package oauth_github

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/oauth"
)

//go:embed logo.svg
var logoBytes []byte

type GitHubOauth struct {
	oauth.BaseOAuthProvider
}

func (g GitHubOauth) GetName() string {
	return "github"
}

func (g GitHubOauth) RegisterRoutes(apiMux *http.ServeMux) {
	apiMux.Handle("/api/auth/oauth/github", githubRedirectHTTP{
		ClientID: g.ClientID,
		Scopes:   g.Scopes,
		BaseURL:  g.BaseURL,
	})
	apiMux.Handle("/api/auth/oauth/github/callback", githubCallbackHTTP{
		ClientID:                     g.ClientID,
		ClientSecret:                 g.ClientSecret,
		BaseURL:                      g.BaseURL,
		Scope:                        g.Scopes,
		Database:                     g.Database,
		RedirectToRegisterOnNotFound: g.RedirectToRegisterOnNotFound,
		GetUserFunc:                  g.GetUserByEmail,
	})
	apiMux.HandleFunc("/api/auth/oauth/github/logo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=2592000")
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(logoBytes)
	})
	fmt.Println("Registered github oauth!")
}
