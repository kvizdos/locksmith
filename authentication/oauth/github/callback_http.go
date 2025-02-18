package oauth_github

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

type githubCallbackHTTP struct {
	Database                     database.DatabaseAccessor
	ClientID                     string
	ClientSecret                 string
	BaseURL                      string
	Scope                        string // e.g., "user:email"
	RedirectToRegisterOnNotFound bool
	GetUserFunc                  func(email string, r *http.Request, customUserInterface ...users.LocksmithUserInterface) users.LocksmithUserInterface
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Username string `json:"login"`
	// There are additional fields but these are the key ones.
}

func (g githubCallbackHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for errors returned in the callback query.
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Error(w, "OAuth error: "+errParam, http.StatusBadRequest)
		return
	}

	// Retrieve the code from query parameters.
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code in callback", http.StatusBadRequest)
		return
	}

	// Prepare form data for token exchange.
	data := url.Values{}
	data.Set("client_id", g.ClientID)
	data.Set("client_secret", g.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", g.BaseURL+"/api/auth/oauth/github/callback")

	// Create the POST request to GitHub's token endpoint.
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, "Failed to create token request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Execute the token request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON token response.
	var tokenResp githubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
		return
	}

	// Now use the access token to fetch the user's email addresses.
	reqEmail, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		http.Error(w, "Failed to create email request", http.StatusInternalServerError)
		return
	}
	reqEmail.Header.Set("Authorization", "token "+tokenResp.AccessToken)
	reqEmail.Header.Set("Accept", "application/vnd.github+json")

	respEmail, err := client.Do(reqEmail)
	if err != nil {
		http.Error(w, "Failed to fetch emails", http.StatusInternalServerError)
		return
	}
	defer respEmail.Body.Close()

	// Decode the JSON response containing the info.
	var info githubEmail
	if err := json.NewDecoder(respEmail.Body).Decode(&info); err != nil {
		http.Error(w, "Failed to decode email response", http.StatusInternalServerError)
		return
	}

	user := g.GetUserFunc(info.Email, r)

	if user == nil {
		if g.RedirectToRegisterOnNotFound {
			http.Redirect(w, r, fmt.Sprintf("/register?err=oauth_email_not_found&username=%s&email=%s", info.Username, info.Email), http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, "/login?err=oauth_email_not_found", http.StatusTemporaryRedirect)
		return
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

	oauth.LoginUser(g.Database, user, "github", validRedirect, w, r)
}
