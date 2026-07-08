package authenticator

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
)

// helper: run attemptLogin with a mock federated session directly
func federatedAttemptLogin(t *testing.T, a *authorizers, sess *mockFederatedSession) (*authenticator_domain.Token, error) {
	t.Helper()
	handler := mockFederatedHandlerFull{session: sess}
	req := httptest.NewRequest("GET", "/", nil)
	token, _, err := a.attemptLogin(context.Background(), handler, req)
	return token, err
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestFederated_ExistingAuthLink_ResolvesUser(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", authentication.PasswordInfo{Passwordless: true}, "user")
	link := map[string]interface{}{
		"provider":  "mock-federated",
		"subject":   "sub-123",
		"issuer":    "https://issuer.example.com",
		"user_id":   "u1",
		"linked_at": int64(0),
	}
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {"u1": user},
		"auth_links": {"l1": link},
	})
	a := newTestAuthorizer(db)

	token, err := federatedAttemptLogin(t, a, &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-123",
		issuer:        "https://issuer.example.com",
		email:         "k@example.com",
		emailVerified: true,
	})

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if token == nil {
		t.Fatal("expected a token, got nil")
	}
}

func TestFederated_IssuerMismatch_Denied(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", authentication.PasswordInfo{Passwordless: true}, "user")
	link := map[string]interface{}{
		"provider":  "mock-federated",
		"subject":   "sub-123",
		"issuer":    "https://real-issuer.example.com",
		"user_id":   "u1",
		"linked_at": int64(0),
	}
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {"u1": user},
		"auth_links": {"l1": link},
	})
	a := newTestAuthorizer(db)

	_, err := federatedAttemptLogin(t, a, &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-123",
		issuer:        "https://attacker.example.com", // mismatch
		email:         "k@example.com",
		emailVerified: true,
	})

	if err == nil {
		t.Fatal("expected error for issuer mismatch, got nil")
	}
	var notFound *authenticator_domain.UserNotFoundError
	if !errors.As(err, &notFound) {
		t.Errorf("expected UserNotFoundError, got %T: %v", err, err)
	}
}

func TestFederated_NoAuthLink_EmailVerified_AutoLinks(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", authentication.PasswordInfo{Passwordless: true}, "user")
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {"u1": user},
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	token, err := federatedAttemptLogin(t, a, &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-new",
		issuer:        "https://issuer.example.com",
		email:         "k@example.com",
		emailVerified: true,
	})

	if err != nil {
		t.Fatalf("expected auto-link success, got: %v", err)
	}
	if token == nil {
		t.Fatal("expected token after auto-link")
	}

	// Verify the auth_link was created
	_, found := db.Tables["auth_links"]
	if !found || len(db.Tables["auth_links"]) == 0 {
		t.Error("expected auth_link to be created in database")
	}
}

func TestFederated_NoAuthLink_EmailNotVerified_Denied(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", authentication.PasswordInfo{Passwordless: true}, "user")
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {"u1": user},
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	_, err := federatedAttemptLogin(t, a, &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-new",
		issuer:        "https://issuer.example.com",
		email:         "k@example.com",
		emailVerified: false, // not verified
	})

	if err == nil {
		t.Fatal("expected error for unverified email, got nil")
	}
	var notFound *authenticator_domain.UserNotFoundError
	if !errors.As(err, &notFound) {
		t.Errorf("expected UserNotFoundError, got %T: %v", err, err)
	}
}

func TestFederated_NoAuthLink_EmailVerified_NoMatchingUser(t *testing.T) {
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {},
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	_, err := federatedAttemptLogin(t, a, &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-new",
		issuer:        "https://issuer.example.com",
		email:         "nobody@example.com",
		emailVerified: true,
	})

	if err == nil {
		t.Fatal("expected error when no matching user, got nil")
	}
	var notFound *authenticator_domain.UserNotFoundError
	if !errors.As(err, &notFound) {
		t.Errorf("expected UserNotFoundError, got %T: %v", err, err)
	}
}

func TestFederated_Rosterable_NoUser_ReturnsHint(t *testing.T) {
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {},
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	_, err := federatedAttemptLogin(t, a, &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-new",
		issuer:        "https://issuer.example.com",
		email:         "new@example.com",
		emailVerified: false, // no auto-link possible
		rosterable:    true,
		displayName:   "New User",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var notFound *authenticator_domain.UserNotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected UserNotFoundError, got %T", err)
	}
	if notFound.RegistrationHint == nil {
		t.Fatal("expected RegistrationHint to be populated")
	}
	if !notFound.RegistrationHint.Rosterable {
		t.Error("expected RegistrationHint.Rosterable to be true")
	}
	if notFound.RegistrationHint.ProviderName != "mock-federated" {
		t.Errorf("expected provider name 'mock-federated', got %q", notFound.RegistrationHint.ProviderName)
	}
}

func TestFederated_NotRosterable_NoUser_NoHint(t *testing.T) {
	db := newTestDB(map[string]map[string]interface{}{
		"users":      {},
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	_, err := federatedAttemptLogin(t, a, &mockFederatedSession{
		provider:      "mock-federated",
		subject:       "sub-new",
		issuer:        "https://issuer.example.com",
		email:         "new@example.com",
		emailVerified: false,
		rosterable:    false,
	})

	var notFound *authenticator_domain.UserNotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected UserNotFoundError, got %T", err)
	}
	if notFound.RegistrationHint != nil {
		t.Error("expected no RegistrationHint for non-rosterable session")
	}
}
