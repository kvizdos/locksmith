package oauth_oidc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"

	"github.com/kvizdos/locksmith/authentication/oauth"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

// OIDCConnection implements a generic OIDC provider.
type OIDCConnection struct {
	oauth.BaseOAuthProvider
	ProviderName string

	Provider *oidc.Provider
	Config   oauth2.Config
	Verifier *oidc.IDTokenVerifier

	LogoBytes      []byte
	DynamicBaseURL func(r *http.Request) string
}

type OIDCConnectionParams struct {
	Issuer                 string
	ClientID               string
	ClientSecret           string
	BaseURL                string
	ProviderName           string
	DB                     database.DatabaseAccessor
	LogoBytes              []byte
	CustomizedGetUserQuery func(email string, r *http.Request) map[string]interface{}
	DynamicBaseURL         func(r *http.Request) string
	LoginInfoCallback      func(method string, user map[string]any)
}

// NewOIDCConnection creates a new OIDCConnection instance.
// providerName should be unique (e.g., "google", "azure", etc.) for each connection.
func NewOIDCConnection(ctx context.Context, params OIDCConnectionParams) (*OIDCConnection, error) {
	provider, err := oidc.NewProvider(ctx, params.Issuer)
	if err != nil {
		return nil, err
	}

	callbackURL := fmt.Sprintf("%s/api/auth/oauth/%s/callback", params.BaseURL, params.ProviderName)

	config := oauth2.Config{
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  callbackURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: params.ClientID})

	return &OIDCConnection{
		BaseOAuthProvider: oauth.BaseOAuthProvider{
			ClientID:               params.ClientID,
			ClientSecret:           params.ClientSecret,
			BaseURL:                callbackURL,
			Database:               params.DB,
			CustomizedGetUserQuery: params.CustomizedGetUserQuery,
			LoginInfoCallback:      params.LoginInfoCallback,
		},
		ProviderName:   params.ProviderName,
		Provider:       provider,
		Config:         config,
		Verifier:       verifier,
		LogoBytes:      params.LogoBytes,
		DynamicBaseURL: params.DynamicBaseURL,
	}, nil
}

// GetName returns the configured provider name.
func (o *OIDCConnection) GetName() string {
	if o.ProviderName != "" {
		return o.ProviderName
	}
	return "oidc"
}

// RegisterRoutes registers the OIDC endpoints.
func (o *OIDCConnection) RegisterRoutes(apiMux *http.ServeMux) {
	apiMux.Handle("/api/auth/oauth/"+o.GetName(), http.HandlerFunc(o.handleRedirect))
	apiMux.Handle("/api/auth/oauth/"+o.GetName()+"/callback", http.HandlerFunc(o.handleCallback))
	apiMux.Handle("/api/auth/oauth/"+o.GetName()+"/refresh", http.HandlerFunc(o.handleRefresh))
	apiMux.HandleFunc("/api/auth/oauth/"+o.GetName()+"/logo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=2592000")
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(o.LogoBytes)
	})
}

// handleRedirect initiates the OIDC flow.
func (o *OIDCConnection) handleRedirect(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("page")
	if state == "" {
		state = "/app"
	} else {
		parsed, err := url.Parse(state)
		if err != nil || parsed.IsAbs() || !strings.HasPrefix(parsed.Path, "/") {
			state = "/app"
		}
	}
	cfg := o.Config
	if o.DynamicBaseURL != nil {
		cfg.RedirectURL = fmt.Sprintf("%s/api/auth/oauth/%s/callback", o.DynamicBaseURL(r), o.ProviderName)
	}
	opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}

	// check cookie
	if c, err := r.Cookie("ls_oauth_hint"); err == nil && c.Value != "" {
		opts = append(opts, oauth2.SetAuthURLParam("login_hint", c.Value))
	}

	authURL := cfg.AuthCodeURL(state, opts...)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleCallback processes the OIDC callback.
func (o *OIDCConnection) handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	cfg := o.Config
	if o.DynamicBaseURL != nil {
		cfg.RedirectURL = fmt.Sprintf("%s/api/auth/oauth/%s/callback", o.DynamicBaseURL(r), o.ProviderName)
	}

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}
	idToken, err := o.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	if !claims.EmailVerified {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	user := o.GetUserByEmail(claims.Email, r)
	if user == nil {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	state := r.URL.Query().Get("state")
	redirectPath := "/app"
	if state != "" {
		parsed, err := url.Parse(state)
		if err == nil && !parsed.IsAbs() && strings.HasPrefix(parsed.Path, "/") {
			redirectPath = parsed.String()
		}
	}

	oauth.LoginUser(o.Database, user, o.GetName(), redirectPath, o.LoginInfoCallback, w, r)
}

func (o *OIDCConnection) handleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	refreshCookie, err := r.Cookie("oidc_refresh")
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusBadRequest)
		return
	}
	if refreshCookie.Value == "" {
		http.Error(w, "missing refresh token", http.StatusBadRequest)
		return
	}

	// Create a token with the refresh token set.
	token := &oauth2.Token{
		RefreshToken: refreshCookie.Value,
	}

	// Use the TokenSource to refresh the token.
	ts := o.Config.TokenSource(ctx, token)
	newToken, err := ts.Token()
	if err != nil {
		http.Error(w, "failed to refresh token: "+err.Error(), http.StatusForbidden)
		return
	}

	// Return the new access token (and other fields as needed)
	fmt.Fprintf(w, "New access token: %s", newToken.AccessToken)
}

// GetUserByEmail retrieves the user given an email.
func (o *OIDCConnection) GetUserByEmail(email string, r *http.Request, customUserInterface ...users.LocksmithUserInterface) users.LocksmithUserInterface {
	return o.BaseOAuthProvider.GetUserByEmail(email, r, customUserInterface...)
}
