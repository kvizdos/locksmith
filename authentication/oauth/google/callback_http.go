package oauth_google

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kvizdos/locksmith/authentication/oauth"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type googleCallbackHTTP struct {
	ClientID                     string
	ClientSecret                 string
	BaseURL                      string
	Scope                        string
	RedirectToRegisterOnNotFound bool
	Database                     database.DatabaseAccessor
	GetUserFunc                  func(email string, r *http.Request, customUserInterface ...users.LocksmithUserInterface) users.LocksmithUserInterface
}

func (g googleCallbackHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// In your callback handler
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		if r.URL.Query().Get("error_subtype") == "access_denied" || errParam == "immediate_failed" {
			// Redirect to full login without 'prompt=none'
			fullLoginURL := "https://accounts.google.com/o/oauth2/auth?" +
				"client_id=" + g.ClientID +
				"&redirect_uri=" + g.BaseURL + "/api/auth/oauth/google/callback" +
				"&scope=" + g.Scope +
				"&response_type=code"
			http.Redirect(w, r, fullLoginURL, http.StatusFound)
			return
		}

		if errParam == "access_denied" {
			http.Redirect(w, r, g.BaseURL+"/login", http.StatusFound)
		}
	}

	// Read the state parameter from the query string.
	state := r.URL.Query().Get("state")

	validRedirect := "/app"
	if state != "" {
		parsed, err := url.Parse(state)
		if err == nil && parsed.Scheme == "" && parsed.Host == "" && strings.HasPrefix(parsed.Path, "/") {
			validRedirect = parsed.String() // includes query and fragment if present
		}
	}

	// Retrieve the code from query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code in callback", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for tokens
	tokenResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {g.ClientID},
		"client_secret": {g.ClientSecret},
		"redirect_uri":  {g.BaseURL + "/api/auth/oauth/google/callback"},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		http.Error(w, "Error exchanging code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()

	var tokenData map[string]interface{}

	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		http.Error(w, "Error decoding token response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// After decoding tokenData
	accessToken, ok := tokenData["access_token"].(string)
	if !ok || accessToken == "" {
		http.Error(w, "Invalid access token", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		http.Error(w, "Error creating user info request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	userInfoResp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error fetching user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	if userInfoResp.StatusCode != http.StatusOK {
		http.Error(w, "Non-OK response from user info endpoint", http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Error decoding user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		http.Error(w, "Email not found in user info", http.StatusInternalServerError)
		return
	}

	emailVerified, ok := userInfo["email_verified"].(bool)
	if !ok {
		http.Error(w, "Email verified not found in user info", http.StatusInternalServerError)
		return
	}

	if !emailVerified {
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	givenName, _ := userInfo["given_name"].(string)
	familyName, _ := userInfo["family_name"].(string)

	user := g.GetUserFunc(email, r)

	if user == nil {
		if g.RedirectToRegisterOnNotFound {
			http.Redirect(w, r, fmt.Sprintf("/register?err=oauth_email_not_found&username=%s%s&email=%s", givenName, familyName, email), http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
	}

	oauth.LoginUser(g.Database, user, "google", validRedirect, w, r)
}
