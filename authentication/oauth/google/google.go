package oauth_google

import (
	_ "embed"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/oauth"
)

//go:embed logo.webp
var logoBytes []byte

type GoogleOauth struct {
	oauth.BaseOAuthProvider
}

func (g GoogleOauth) GetName() string {
	return "google"
}

func (g GoogleOauth) RegisterRoutes(apiMux *http.ServeMux) {
	apiMux.Handle("/api/auth/oauth/google", googleRedirectHTTP{
		ClientID: g.ClientID,
		Scopes:   g.Scopes,
		BaseURL:  g.BaseURL,
	})
	apiMux.Handle("/api/auth/oauth/google/callback", googleCallbackHTTP{
		ClientID:                     g.ClientID,
		ClientSecret:                 g.ClientSecret,
		BaseURL:                      g.BaseURL,
		Scope:                        g.Scopes,
		RedirectToRegisterOnNotFound: g.RedirectToRegisterOnNotFound,
		Database:                     g.Database,
		GetUserFunc:                  g.GetUserByEmail,
	})
	apiMux.HandleFunc("/api/auth/oauth/google/logo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=2592000")
		w.Header().Set("Content-Type", "image/webp")
		w.Write(logoBytes)
	})
}
