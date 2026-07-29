package method_password

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func TestPasswordValidatorCanHandle(t *testing.T) {
	t.Parallel()

	pv := NewPasswordValidator().(passwordValidator)

	postReq := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	if !pv.CanHandle(postReq) {
		t.Fatal("expected CanHandle to be true for POST")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	if pv.CanHandle(getReq) {
		t.Fatal("expected CanHandle to be false for GET")
	}
}

func TestPasswordValidatorName(t *testing.T) {
	t.Parallel()

	pv := NewPasswordValidator().(passwordValidator)
	if pv.Name() != "password" {
		t.Fatalf("expected name 'password', got %s", pv.Name())
	}
}

func TestPasswordValidatorPasswordless(t *testing.T) {
	t.Parallel()

	pv := NewPasswordValidator().(passwordValidator)
	if pv.Passwordless() {
		t.Fatal("expected Passwordless to be false")
	}
}

func TestPasswordValidatorSessionReturnsSession(t *testing.T) {
	t.Parallel()

	pv := NewPasswordValidator().(passwordValidator)
	db := database.TestDatabase{Tables: map[string]map[string]interface{}{}}

	session := pv.Session(db)
	if session == nil {
		t.Fatal("expected non-nil session")
	}
}

func TestLoadRequestRejectsNonPost(t *testing.T) {
	t.Parallel()

	session := newPasswordValidatorSession(database.TestDatabase{}, authenticator_methods.DefaultPasswordValidatorOptions())

	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	err := session.LoadRequest(req)

	if !errors.Is(err, authenticator_domain.ErrMethodNotSupported) {
		t.Fatalf("expected ErrMethodNotSupported, got %v", err)
	}
}

func TestLoadRequestValidJSON(t *testing.T) {
	t.Parallel()

	session := newPasswordValidatorSession(database.TestDatabase{}, authenticator_methods.DefaultPasswordValidatorOptions())

	body := `{"username":"kenton","password":"supersecret"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	if err := session.LoadRequest(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.GetPresentedUser() != "kenton" {
		t.Fatalf("expected presented user 'kenton', got %s", session.GetPresentedUser())
	}
	if session.presentedPassword != "supersecret" {
		t.Fatalf("expected presented password 'supersecret', got %s", session.presentedPassword)
	}
}

func TestLoadRequestMalformedJSON(t *testing.T) {
	t.Parallel()

	session := newPasswordValidatorSession(database.TestDatabase{}, authenticator_methods.DefaultPasswordValidatorOptions())

	body := `{"username":"kenton",` // malformed
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	err := session.LoadRequest(req)
	if !errors.Is(err, authenticator_domain.ErrFailedToParse) {
		t.Fatalf("expected ErrFailedToParse, got %v", err)
	}
}

func compiledPassword(plain string) authentication.PasswordInfo {
	p, _ := authentication.CompileLocksmithPassword(plain)
	return p
}

func TestIsAuthorizedCorrectPassword(t *testing.T) {
	t.Parallel()

	session := newPasswordValidatorSession(database.TestDatabase{}, authenticator_methods.DefaultPasswordValidatorOptions())
	session.presentedPassword = "correct-password"

	user := users.LocksmithUser{
		ID:           "user-1",
		Username:     "kenton",
		PasswordInfo: compiledPassword("correct-password"),
	}

	if err := session.IsAuthorized(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsAuthorizedWrongPassword(t *testing.T) {
	t.Parallel()

	session := newPasswordValidatorSession(database.TestDatabase{}, authenticator_methods.DefaultPasswordValidatorOptions())
	session.presentedPassword = "wrong-password"

	user := users.LocksmithUser{
		ID:           "user-1",
		Username:     "kenton",
		PasswordInfo: compiledPassword("correct-password"),
	}

	err := session.IsAuthorized(user)
	if !errors.Is(err, authenticator_domain.ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
