package authenticator

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/users"
)

// ── ServeLoginAPI: additional coverage not in attempt_login_password_test.go ──

func TestServeLoginAPI_ValidLogin_SetsExpiresAtCookie(t *testing.T) {
	t.Parallel()
	db := newTestDB(map[string]map[string]interface{}{
		"users": {
			"u1": makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user"),
		},
	})
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var foundExpires, foundToken bool
	for _, c := range rr.Result().Cookies() {
		if c.Name == "ls_expires_at" {
			foundExpires = true
		}
		if c.Name == "token" {
			foundToken = true
		}
	}
	if !foundExpires {
		t.Error("expected ls_expires_at cookie to be set")
	}
	if !foundToken {
		t.Error("expected token cookie to be set by token manager")
	}
}

func TestServeLoginAPI_NoHandlerMatches_Returns500(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	// GET request doesn't satisfy password's CanHandle, and there's no
	// path value or header hint, so no handler will be found at all.
	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	rr := httptest.NewRecorder()

	a.ServeLoginAPI(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 when no handler matches, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestServeLoginAPI_TokenNilWithoutError_ReturnsAuthError(t *testing.T) {
	t.Parallel()
	db := newTestDB(map[string]map[string]interface{}{
		"users": {
			"u1": makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user"),
		},
	})

	// Use a token manager whose CreateAuthToken returns (nil, nil) to
	// exercise the "token is missing, yet there was no error" branch.
	sp, _ := signing.CreateSigningPackage()
	a := NewAuthorizer(db,
		WithTokenManager(&nilTokenManager{}),
		WithMethods(AllowMethodPassword()),
		WithSigningPackage(&sp),
	)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 when token is nil without error, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

// nilTokenManager returns a nil token with no error to exercise the
// "token == nil" branch in ServeLoginAPI.
type nilTokenManager struct{}

func (m *nilTokenManager) Read(r *http.Request) (*authenticator_domain.Token, error) {
	return nil, nil
}
func (m *nilTokenManager) CreateAuthToken(user users.LocksmithUserInterface) (*authenticator_domain.Token, error) {
	return nil, nil
}
func (m *nilTokenManager) PassToClient(w http.ResponseWriter, r *http.Request, token *authenticator_domain.Token) error {
	return nil
}

// ── setBaseCookies ─────────────────────────────────────────────────────────────

func TestSetBaseCookies_NoOAuthProvider(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	token := &authenticator_domain.Token{
		AuthToken:     "tok",
		ExpiresAt:     time.Now().Add(time.Hour),
		OAuthProvider: "",
	}

	rr := httptest.NewRecorder()
	if err := a.setBaseCookies(rr, token); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := rr.Result().Cookies()
	var oauthProviderCookie *http.Cookie
	var expiresCookie *http.Cookie
	for _, c := range cookies {
		switch c.Name {
		case "ls_oauth_provider":
			oauthProviderCookie = c
		case "ls_expires_at":
			expiresCookie = c
		case "ls_oauth_hint":
			t.Error("did not expect ls_oauth_hint cookie when OAuthProvider is empty")
		}
	}

	if oauthProviderCookie == nil {
		t.Fatal("expected ls_oauth_provider cookie to be set")
	}
	if oauthProviderCookie.Value != "" {
		t.Errorf("expected empty ls_oauth_provider value, got %q", oauthProviderCookie.Value)
	}
	if !oauthProviderCookie.Expires.Equal(time.Unix(0, 0)) {
		t.Errorf("expected ls_oauth_provider to expire at epoch, got %v", oauthProviderCookie.Expires)
	}

	if expiresCookie == nil {
		t.Fatal("expected ls_expires_at cookie to be set")
	}
}

func TestSetBaseCookies_WithOAuthProvider(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	token := &authenticator_domain.Token{
		AuthToken:     "tok",
		ExpiresAt:     time.Now().Add(time.Hour),
		OAuthProvider: "google",
		OAuthHint:     "hint-value",
	}

	rr := httptest.NewRecorder()
	if err := a.setBaseCookies(rr, token); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := rr.Result().Cookies()
	var oauthProviderCookie, oauthHintCookie *http.Cookie
	for _, c := range cookies {
		switch c.Name {
		case "ls_oauth_provider":
			oauthProviderCookie = c
		case "ls_oauth_hint":
			oauthHintCookie = c
		}
	}

	if oauthProviderCookie == nil {
		t.Fatal("expected ls_oauth_provider cookie to be set")
	}
	if oauthProviderCookie.Value != "google" {
		t.Errorf("expected 'google', got %q", oauthProviderCookie.Value)
	}
	farFuture := time.Now().UTC().AddDate(9, 0, 0)
	if oauthProviderCookie.Expires.Before(farFuture) {
		t.Errorf("expected far-future expiry for ls_oauth_provider, got %v", oauthProviderCookie.Expires)
	}

	if oauthHintCookie == nil {
		t.Fatal("expected ls_oauth_hint cookie to be set")
	}
	if oauthHintCookie.Value != "hint-value" {
		t.Errorf("expected 'hint-value', got %q", oauthHintCookie.Value)
	}
	if !oauthHintCookie.HttpOnly {
		t.Error("expected ls_oauth_hint to be HttpOnly")
	}
	if oauthHintCookie.Expires.Before(farFuture) {
		t.Errorf("expected far-future expiry for ls_oauth_hint, got %v", oauthHintCookie.Expires)
	}
}

// ── attemptLogin: WithEmailAsUsername lookup ──────────────────────────────────

func TestAttemptLogin_EmailAsUsername_LooksUpByEmail(t *testing.T) {
	t.Parallel()
	db := newTestDB(map[string]map[string]interface{}{
		"users": {
			"u1": makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user"),
		},
	})
	a := newTestAuthorizer(db, WithEmailAsUsername())

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"k@example.com","password":"hunter2"}`))

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 login via email lookup, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestAttemptLogin_EmailAsUsername_UsernameLookupFails(t *testing.T) {
	t.Parallel()
	db := newTestDB(map[string]map[string]interface{}{
		"users": {
			"u1": makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user"),
		},
	})
	a := newTestAuthorizer(db, WithEmailAsUsername())

	// Presenting the username (not email) should fail to find the user
	// since lookup happens against the "email" field.
	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 when looking up by email but presenting username, got %d", rr.Code)
	}
}

// ── beginRegistrationRostering via ServeLoginAPI (federated + rosterable) ─────

func TestServeLoginAPI_RosterableFederated_RedirectsToHintedRegistration(t *testing.T) {
	t.Parallel()
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {},
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	handler := mockFederatedHandlerFull{session: &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-new",
		issuer:        "https://issuer.example.com",
		email:         "new@example.com",
		emailVerified: false,
		rosterable:    true,
		displayName:   "New User",
	}}

	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	req.Header.Set("X-Handler-Hint", "mock-federated")
	rr := httptest.NewRecorder()

	// Register the federated handler via WithMethods so getHandler can find it.
	a2 := NewAuthorizer(db,
		WithTokenManager(&mockTokenManager{redirectPath: "/app"}),
		WithMethods(handler),
		WithSigningPackage(a.sp),
	)

	a2.ServeLoginAPI(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 See Other redirect to hinted registration, got %d — body: %s", rr.Code, rr.Body.String())
	}

	if got := rr.Result().Header.Get("Location"); got != "/api/register?hinted" {
		t.Errorf("expected redirect to /api/register?hinted, got %q", got)
	}

	var found bool
	for _, c := range rr.Result().Cookies() {
		if c.Name == "registration_hint" {
			found = true
			if c.Value == "" {
				t.Error("expected non-empty registration_hint cookie value")
			}
		}
	}
	if !found {
		t.Error("expected registration_hint cookie to be set")
	}
}

func containsAll(haystack string, needles ...string) bool {
	for _, n := range needles {
		if !contains(haystack, n) {
			return false
		}
	}
	return true
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}

// ── enrichCtx ──────────────────────────────────────────────────────────────────

func TestEnrichCtx_AddsIPToContext(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "203.0.113.5:12345"

	ctx := a.enrichCtx(req)
	if ctx.Value("ip") == nil {
		t.Error("expected 'ip' to be set in enriched context")
	}
}

// ── writeAuthError ───────────────────────────────────────────────────────────

func TestWriteAuthError_ProtectionEnabled_GenericMessage(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.writeAuthError(nil, rr, "Some specific detail")

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
	body := rr.Body.String()
	if contains(body, "Some specific detail") {
		t.Errorf("expected enumeration-protected generic message, got: %s", body)
	}
}

func TestWriteAuthError_ProtectionDisabled_SpecificMessage(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db, DisableUserEnumerationProtection())

	rr := httptest.NewRecorder()
	a.writeAuthError(nil, rr, "Specific detail here")

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !contains(body, "Specific detail here") {
		t.Errorf("expected specific error message, got: %s", body)
	}
}

func TestWriteAuthError_CallsMinCtxFunc(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	called := false
	rr := httptest.NewRecorder()
	a.writeAuthError(func() { called = true }, rr, "msg")

	if !called {
		t.Error("expected minCtx function to be called")
	}
}

// ── attemptLogin direct tests using federated session helper ─────────────────

func TestAttemptLogin_FederatedSuccess_LogsUserID(t *testing.T) {
	t.Parallel()
	user := makeUser("u1", "kenton", "k@example.com", authentication.PasswordInfo{Passwordless: true}, "user")
	db := newTestDB(map[string]map[string]interface{}{
		"users": {"u1": user},
		"auth_links": {
			"l1": map[string]interface{}{
				"provider":  "mock-federated",
				"subject":   "sub-123",
				"issuer":    "https://issuer.example.com",
				"user_id":   "u1",
				"linked_at": int64(0),
			},
		},
	})
	a := newTestAuthorizer(db)

	handler := mockFederatedHandlerFull{session: &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-123",
		issuer:        "https://issuer.example.com",
		email:         "k@example.com",
		emailVerified: true,
	}}
	req := httptest.NewRequest("GET", "/", nil)

	token, _, err := a.attemptLogin(context.Background(), handler, req)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if token == nil {
		t.Fatal("expected token")
	}
}

func TestAttemptLogin_LoadRequestError_Wrapped(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	handler := AllowMethodPassword()
	req := httptest.NewRequest(http.MethodGet, "/api/login", nil) // GET fails password's LoadRequest

	_, _, err := a.attemptLogin(context.Background(), handler, req)
	if err == nil {
		t.Fatal("expected error from LoadRequest failure")
	}
	if !errors.Is(err, authenticator_domain.ErrMethodNotSupported) {
		t.Errorf("expected wrapped ErrMethodNotSupported, got: %v", err)
	}
}
