package tokens

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func TestSetBaseCookiesNoOAuthProvider(t *testing.T) {
	t.Parallel()

	token := &authenticator_domain.Token{
		AuthToken: "tok",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	rr := httptest.NewRecorder()
	SetBaseCookies(rr, token)

	cookies := rr.Result().Cookies()
	var expires, oauthProvider *string
	var oauthHintFound bool
	for _, c := range cookies {
		switch c.Name {
		case "ls_expires_at":
			v := c.Value
			expires = &v
		case "ls_oauth_provider":
			v := c.Value
			oauthProvider = &v
			if !c.Expires.Equal(time.Unix(0, 0)) {
				t.Fatalf("ls_oauth_provider expires = %v, want epoch", c.Expires)
			}
		case "ls_oauth_hint":
			oauthHintFound = true
		}
	}
	if expires == nil {
		t.Fatal("expected ls_expires_at cookie to be set")
	}
	if oauthProvider == nil || *oauthProvider != "" {
		t.Fatal("expected ls_oauth_provider cookie to be cleared")
	}
	if oauthHintFound {
		t.Fatal("ls_oauth_hint should not be set without an OAuth provider")
	}
}

func TestSetBaseCookiesWithOAuthProvider(t *testing.T) {
	t.Parallel()

	token := &authenticator_domain.Token{
		AuthToken:     "tok",
		ExpiresAt:     time.Now().Add(time.Hour),
		OAuthProvider: "google",
		OAuthHint:     "hint-value",
	}

	rr := httptest.NewRecorder()
	SetBaseCookies(rr, token)

	cookies := rr.Result().Cookies()
	var sawProvider, sawHint bool
	for _, c := range cookies {
		switch c.Name {
		case "ls_oauth_provider":
			sawProvider = true
			if c.Value != "google" {
				t.Fatalf("ls_oauth_provider value = %q, want %q", c.Value, "google")
			}
		case "ls_oauth_hint":
			sawHint = true
			if c.Value != "hint-value" {
				t.Fatalf("ls_oauth_hint value = %q, want %q", c.Value, "hint-value")
			}
			if !c.HttpOnly {
				t.Fatal("ls_oauth_hint must be HttpOnly")
			}
		}
	}
	if !sawProvider {
		t.Fatal("expected ls_oauth_provider cookie to be set")
	}
	if !sawHint {
		t.Fatal("expected ls_oauth_hint cookie to be set")
	}
}

func TestPassToClientRedirectsToConfiguredDefault(t *testing.T) {
	t.Parallel()

	db := database.TestDatabase{Tables: map[string]map[string]interface{}{}}
	cm := NewCookieManager(db, "/app")

	token := &authenticator_domain.Token{
		AuthToken: "tok",
		ExpiresAt: time.Now().Add(time.Hour),
		User:      &users.LocksmithUser{ID: "u1"},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	rr := httptest.NewRecorder()

	if err := cm.PassToClient(rr, req, token); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusSeeOther)
	}
	if got := rr.Header().Get("Location"); got != "/app" {
		t.Fatalf("Location = %q, want %q", got, "/app")
	}
}

func TestPassToClientRedirectsToTokenOverride(t *testing.T) {
	t.Parallel()

	db := database.TestDatabase{Tables: map[string]map[string]interface{}{}}
	cm := NewCookieManager(db, "/app")

	token := &authenticator_domain.Token{
		AuthToken:    "tok",
		ExpiresAt:    time.Now().Add(time.Hour),
		User:         &users.LocksmithUser{ID: "u1"},
		RedirectPath: "/dashboard",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	rr := httptest.NewRecorder()

	if err := cm.PassToClient(rr, req, token); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := rr.Header().Get("Location"); got != "/dashboard" {
		t.Fatalf("Location = %q, want %q", got, "/dashboard")
	}
}
