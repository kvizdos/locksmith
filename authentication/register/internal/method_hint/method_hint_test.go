package method_hint

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/register/register_methods"
	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/authentication/signing"
)

func TestHintRegistrationCanHandle(t *testing.T) {
	t.Parallel()

	handler := NewHintRegistration()

	cookieReq := httptest.NewRequest(http.MethodGet, "/register?hinted", nil)
	cookieReq.AddCookie(&http.Cookie{Name: registrationhints.CookieName, Value: "token"})
	if !handler.CanHandle(cookieReq) {
		t.Fatal("expected cookie hint request to be handled")
	}

	authReq := httptest.NewRequest(http.MethodGet, "/register?hinted", nil)
	authReq.Header.Set("Authorization", "RegistrationHint token")
	if !handler.CanHandle(authReq) {
		t.Fatal("expected authorization hint request to be handled")
	}

	plainReq := httptest.NewRequest(http.MethodGet, "/register?hinted", nil)
	if handler.CanHandle(plainReq) {
		t.Fatal("expected request without hint to be ignored")
	}

	missingQueryReq := httptest.NewRequest(http.MethodGet, "/register", nil)
	missingQueryReq.AddCookie(&http.Cookie{Name: registrationhints.CookieName, Value: "token"})
	if handler.CanHandle(missingQueryReq) {
		t.Fatal("expected request without ?hinted to be ignored")
	}

	wrongMethodReq := httptest.NewRequest(http.MethodPost, "/register?hinted", nil)
	wrongMethodReq.AddCookie(&http.Cookie{Name: registrationhints.CookieName, Value: "token"})
	if handler.CanHandle(wrongMethodReq) {
		t.Fatal("expected non-GET request to be ignored")
	}
}

func TestHintRegistrationSessionLoadRequestFromAuthorization(t *testing.T) {
	t.Parallel()

	signer, err := signing.CreateSigningPackage()
	if err != nil {
		t.Fatalf("create signer: %v", err)
	}
	svc := registrationhints.Service{Signer: signer, TTL: time.Minute}
	token, err := svc.Create(registrationhints.Hint{
		ProviderName: "google",
		Email:        "alice@example.com",
		DisplayName:  "Alice",
		Issuer:       "https://issuer.example.com",
		Subject:      "subject",
		Rosterable:   true,
	})
	if err != nil {
		t.Fatalf("create hint: %v", err)
	}

	handler := NewHintRegistration(register_methods.WithHintService(svc))
	session := handler.Session(nil)
	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	req.Header.Set("Authorization", "RegistrationHint "+token)

	if err := session.LoadRequest(req); err != nil {
		t.Fatalf("load request: %v", err)
	}

	got := session.RegistrationRequest()
	if got.Username != "alice@example.com" || got.Email != "alice@example.com" || got.Password != "" {
		t.Fatalf("unexpected registration request: %+v", got)
	}
	if got.Hint == nil || got.Hint.ProviderName != "google" || got.Hint.Subject != "subject" {
		t.Fatalf("unexpected hint: %+v", got.Hint)
	}
}
