package method_oidc

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"github.com/kvizdos/locksmith/database"
)

func TestDetectFlowCode(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/login?state=abc&code=xyz", nil)
	if got := detectFlow(req); got != flowCode {
		t.Fatalf("expected flowCode, got %v", got)
	}
}

func TestDetectFlowNoneOnGetWithoutParams(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	if got := detectFlow(req); got != flowNone {
		t.Fatalf("expected flowNone, got %v", got)
	}
}

func TestDetectFlowNoneOnGetWithOnlyState(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/login?state=abc", nil)
	if got := detectFlow(req); got != flowNone {
		t.Fatalf("expected flowNone, got %v", got)
	}
}

func TestDetectFlowCredentialOnPost(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	if got := detectFlow(req); got != flowCredential {
		t.Fatalf("expected flowCredential, got %v", got)
	}
}

func TestOidcHandlerCanHandle(t *testing.T) {
	t.Parallel()

	h := oidcHandler{options: authenticator_methods.OIDCValidatorOptions{ProviderName: "google"}}

	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{
			name: "authorization code callback",
			req:  httptest.NewRequest(http.MethodGet, "/api/login?state=abc&code=xyz", nil),
			want: true,
		},
		{
			name: "fedcm credential post",
			req:  httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"credential":"token"}`)),
			want: true,
		},
		{
			name: "get without callback params",
			req:  httptest.NewRequest(http.MethodGet, "/api/login", nil),
			want: false,
		},
		{
			name: "get with only code",
			req:  httptest.NewRequest(http.MethodGet, "/api/login?code=xyz", nil),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.CanHandle(tt.req); got != tt.want {
				t.Fatalf("CanHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOidcHandlerNameAndPasswordless(t *testing.T) {
	t.Parallel()

	h := oidcHandler{options: authenticator_methods.OIDCValidatorOptions{ProviderName: "google"}}

	if h.Name() != "oidc-google" {
		t.Fatalf("expected name 'oidc-google', got %s", h.Name())
	}
	if !h.Passwordless() {
		t.Fatal("expected Passwordless to be true")
	}
}

func TestOidcHandlerSession(t *testing.T) {
	t.Parallel()

	h := oidcHandler{options: authenticator_methods.OIDCValidatorOptions{ProviderName: "google"}}
	db := database.TestDatabase{Tables: map[string]map[string]interface{}{}}

	session := h.Session(db)
	if session == nil {
		t.Fatal("expected non-nil session")
	}
}

func TestNewOIDCValidatorAppliesOptions(t *testing.T) {
	t.Parallel()

	h := NewOIDCValidator(func(opts *authenticator_methods.OIDCValidatorOptions) {
		opts.ProviderName = "google"
		opts.Rosterable = true
	})

	if h.Name() != "oidc-google" {
		t.Fatalf("expected option-provided name, got %s", h.Name())
	}
}

func TestNewOIDCValidatorPanicsWithoutProviderName(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when ProviderName is not set")
		}
	}()

	NewOIDCValidator()
}

func TestLoadRequestFlowCodeMissingPKCECookie(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})

	req := httptest.NewRequest(http.MethodGet, "/api/login?state=abc&code=xyz", nil)
	err := session.LoadRequest(req)
	if err == nil {
		t.Fatal("expected error due to missing pkce cookie")
	}
}

func TestLoadRequestFlowCodeSuccess(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})

	req := httptest.NewRequest(http.MethodGet, "/api/login?state=abc&code=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "ls_oidc_pkce", Value: "verifier-value"})

	if err := session.LoadRequest(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.pkceVerifier != "verifier-value" {
		t.Fatalf("expected pkceVerifier 'verifier-value', got %s", session.pkceVerifier)
	}
	if session.untrustedParsedCode != "xyz" {
		t.Fatalf("expected untrustedParsedCode 'xyz', got %s", session.untrustedParsedCode)
	}
	if session.flow != flowCode {
		t.Fatalf("expected flow flowCode, got %v", session.flow)
	}
}

func TestLoadRequestFlowCredentialSuccess(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})

	body := `{"credential":"some-jwt-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))

	if err := session.LoadRequest(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.untrustedCredentialToken != "some-jwt-token" {
		t.Fatalf("expected untrustedCredentialToken 'some-jwt-token', got %s", session.untrustedCredentialToken)
	}
	if session.flow != flowCredential {
		t.Fatalf("expected flow flowCredential, got %v", session.flow)
	}
}

func TestLoadRequestFlowCredentialMissingField(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))

	err := session.LoadRequest(req)
	if err == nil {
		t.Fatal("expected error for missing credential field")
	}
}

func TestLoadRequestFlowCredentialMalformedJSON(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})

	body := `not-json`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))

	err := session.LoadRequest(req)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestLoadRequestUnsupportedFlow(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})
	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)

	err := session.LoadRequest(req)
	if err == nil {
		t.Fatal("expected error for unsupported flow")
	}
	if session.flow != flowNone {
		t.Fatalf("expected flowNone to be recorded, got %v", session.flow)
	}
}

func TestIsAuthorizedAlwaysNil(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})

	if err := session.IsAuthorized(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRegistrationHintNotRosterable(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		Rosterable:   false,
	})

	if hint := session.RegistrationHint(); hint != nil {
		t.Fatalf("expected nil hint when not rosterable, got %+v", hint)
	}
}

func TestRegistrationHintRosterable(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		Rosterable:   true,
	})
	session.SetEmail("kenton@example.com")
	session.displayName = "Kenton V"

	hint := session.RegistrationHint()
	if hint == nil {
		t.Fatal("expected non-nil hint when rosterable")
	}
	if hint.ProviderName != "google" {
		t.Fatalf("unexpected provider name: %s", hint.ProviderName)
	}
	if hint.Email != "kenton@example.com" {
		t.Fatalf("unexpected email: %s", hint.Email)
	}
	if hint.DisplayName != "Kenton V" {
		t.Fatalf("unexpected display name: %s", hint.DisplayName)
	}
	if !hint.Rosterable {
		t.Fatal("expected Rosterable to be true")
	}
}
