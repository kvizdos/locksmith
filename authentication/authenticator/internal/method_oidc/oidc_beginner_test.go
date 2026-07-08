package method_oidc

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"golang.org/x/oauth2"
)

func TestBeginRejectsUnsupportedMethod(t *testing.T) {
	t.Parallel()

	h := oidcHandler{options: authenticator_methods.OIDCValidatorOptions{ProviderName: "google"}}

	req := httptest.NewRequest(http.MethodDelete, "/api/start/google", nil)
	rr := httptest.NewRecorder()

	err := h.Begin(context.Background(), rr, req)
	if !errors.Is(err, authenticator_domain.ErrMethodNotSupported) {
		t.Fatalf("expected ErrMethodNotSupported, got %v", err)
	}
}

func TestBeginRejectsNilOauthConfig(t *testing.T) {
	t.Parallel()

	h := oidcHandler{options: authenticator_methods.OIDCValidatorOptions{ProviderName: "google", OauthConfig: nil}}

	req := httptest.NewRequest(http.MethodGet, "/api/start/google", nil)
	rr := httptest.NewRecorder()

	err := h.Begin(context.Background(), rr, req)
	if err == nil {
		t.Fatal("expected error when oauth config is nil")
	}
}

func TestBeginSetsCookiesAndRedirects(t *testing.T) {
	t.Parallel()

	h := oidcHandler{options: authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		OauthConfig: &oauth2.Config{
			ClientID:    "client-id",
			RedirectURL: "https://example.com/api/login",
			Endpoint: oauth2.Endpoint{
				AuthURL: "https://provider.example.com/auth",
			},
		},
	}}

	req := httptest.NewRequest(http.MethodGet, "/api/start/google", nil)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(context.Background(), "log", slog.Default())

	if err := h.Begin(ctx, rr, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rr.Code != http.StatusFound {
		t.Fatalf("expected 302 redirect, got %d", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location == "" {
		t.Fatal("expected Location header to be set")
	}

	parsedLocation, err := url.Parse(location)
	if err != nil {
		t.Fatalf("expected valid redirect location, got %q: %v", location, err)
	}
	if parsedLocation.Scheme != "https" || parsedLocation.Host != "provider.example.com" || parsedLocation.Path != "/auth" {
		t.Fatalf("unexpected redirect destination: %s", location)
	}
	query := parsedLocation.Query()
	if query.Get("client_id") != "client-id" {
		t.Fatalf("expected client_id to be set, got %q", query.Get("client_id"))
	}
	if query.Get("redirect_uri") != "https://example.com/api/login" {
		t.Fatalf("expected redirect_uri to be set, got %q", query.Get("redirect_uri"))
	}
	if query.Get("response_type") != "code" {
		t.Fatalf("expected response_type code, got %q", query.Get("response_type"))
	}
	if query.Get("code_challenge") == "" {
		t.Fatal("expected code_challenge to be set")
	}
	if query.Get("code_challenge_method") != "S256" {
		t.Fatalf("expected S256 code_challenge_method, got %q", query.Get("code_challenge_method"))
	}

	cookies := map[string]*http.Cookie{}
	for _, c := range rr.Result().Cookies() {
		cookies[c.Name] = c
	}

	stateCookie := cookies["ls_oidc_state"]
	if stateCookie == nil {
		t.Fatal("expected ls_oidc_state cookie to be set")
	}
	if stateCookie.Value == "" {
		t.Error("expected non-empty state cookie value")
	}
	if stateCookie.Value != query.Get("state") {
		t.Fatalf("expected state cookie to match auth URL state, got cookie %q and query %q", stateCookie.Value, query.Get("state"))
	}
	assertOIDCCookie(t, stateCookie, false)

	pkceCookie := cookies["ls_oidc_pkce"]
	if pkceCookie == nil {
		t.Fatal("expected ls_oidc_pkce cookie to be set")
	}
	if pkceCookie.Value == "" {
		t.Error("expected non-empty pkce cookie value")
	}
	if pkceChallenge(pkceCookie.Value) != query.Get("code_challenge") {
		t.Fatal("expected code_challenge to be derived from pkce verifier cookie")
	}
	assertOIDCCookie(t, pkceCookie, false)
}

func TestBeginSetsSecureCookiesForTLSRequest(t *testing.T) {
	t.Parallel()

	h := oidcHandler{options: authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		OauthConfig: &oauth2.Config{
			ClientID:    "client-id",
			RedirectURL: "https://example.com/api/login",
			Endpoint: oauth2.Endpoint{
				AuthURL: "https://provider.example.com/auth",
			},
		},
	}}

	req := httptest.NewRequest(http.MethodGet, "https://example.com/api/start/google", nil)
	rr := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), "log", slog.Default())

	if err := h.Begin(ctx, rr, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	seenSecureCookies := map[string]bool{}
	for _, c := range rr.Result().Cookies() {
		if c.Name == "ls_oidc_state" || c.Name == "ls_oidc_pkce" {
			assertOIDCCookie(t, c, true)
			seenSecureCookies[c.Name] = true
		}
	}
	if !seenSecureCookies["ls_oidc_state"] {
		t.Fatal("expected ls_oidc_state cookie to be set")
	}
	if !seenSecureCookies["ls_oidc_pkce"] {
		t.Fatal("expected ls_oidc_pkce cookie to be set")
	}
}

func assertOIDCCookie(t *testing.T, c *http.Cookie, wantSecure bool) {
	t.Helper()

	if c.Path != "/" {
		t.Fatalf("expected %s cookie path /, got %q", c.Name, c.Path)
	}
	if !c.HttpOnly {
		t.Fatalf("expected %s cookie to be HttpOnly", c.Name)
	}
	if c.Secure != wantSecure {
		t.Fatalf("expected %s cookie Secure=%v, got %v", c.Name, wantSecure, c.Secure)
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected %s cookie SameSite=Lax, got %v", c.Name, c.SameSite)
	}
	if c.MaxAge != 600 {
		t.Fatalf("expected %s cookie MaxAge 600, got %d", c.Name, c.MaxAge)
	}
}

func TestRandomURLSafeProducesDistinctValues(t *testing.T) {
	t.Parallel()

	a, err := randomURLSafe(32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := randomURLSafe(32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a == "" || b == "" {
		t.Fatal("expected non-empty random strings")
	}
	if a == b {
		t.Fatal("expected distinct random strings across calls")
	}
}

func TestRandomURLSafeLength(t *testing.T) {
	t.Parallel()

	s, err := randomURLSafe(16)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// base64 RawURLEncoding of 16 bytes should be 22 chars (no padding).
	if len(s) != 22 {
		t.Fatalf("expected length 22 for 16-byte input, got %d", len(s))
	}
}

func TestPkceChallengeDeterministic(t *testing.T) {
	t.Parallel()

	verifier := "some-fixed-verifier-value"
	c1 := pkceChallenge(verifier)
	c2 := pkceChallenge(verifier)

	if c1 != c2 {
		t.Fatalf("expected deterministic challenge for the same verifier, got %s vs %s", c1, c2)
	}
	if c1 == verifier {
		t.Fatal("expected challenge to differ from the verifier")
	}
}

func TestPkceChallengeDiffersForDifferentVerifiers(t *testing.T) {
	t.Parallel()

	c1 := pkceChallenge("verifier-one")
	c2 := pkceChallenge("verifier-two")

	if c1 == c2 {
		t.Fatal("expected different challenges for different verifiers")
	}
}
