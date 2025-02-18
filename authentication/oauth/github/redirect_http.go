package oauth_github

import (
	"net/http"
	"net/url"
	"strings"
)

type githubRedirectHTTP struct {
	ClientID string
	Scopes   string
	BaseURL  string
}

func (g githubRedirectHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encodedScopes, err := url.QueryUnescape(g.Scopes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed"))
		return
	}

	// Retrieve the current page from a query parameter.
	rawPage := r.URL.Query().Get("page")

	// Validate that the page is a relative path.
	// If the page is not provided or is invalid, default to "/app".
	validPage := "/app"
	if rawPage != "" {
		// Parse the rawPage.
		parsed, err := url.Parse(rawPage)
		if err == nil && parsed.Scheme == "" && parsed.Host == "" && strings.HasPrefix(parsed.Path, "/") {
			validPage = parsed.Path
			if parsed.RawQuery != "" {
				validPage += "?" + parsed.RawQuery
			}
			if parsed.Fragment != "" {
				validPage += "#" + parsed.Fragment
			}
		}
	}

	// URL-encode the valid page value.
	encodedState := url.QueryEscape(validPage)

	// Build your GitHub OAuth URL. Replace the placeholders with your actual values.
	githubAuthURL := "https://github.com/login/oauth/authorize?" +
		"client_id=" + g.ClientID +
		"&redirect_uri=" + g.BaseURL + "/api/auth/oauth/github/callback" +
		"&scope=" + encodedScopes +
		"&state=" + encodedState

	http.Redirect(w, r, githubAuthURL, http.StatusFound)
}
