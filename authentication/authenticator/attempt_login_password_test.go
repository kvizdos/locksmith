package authenticator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kvizdos/locksmith/authentication"
)

func TestPassword_ValidCredentials(t *testing.T) {
	db := newTestDB(map[string]map[string]any{
		"users": {
			"u1": makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user"),
		},
	})
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d — body: %s", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("Location") != "/app" {
		t.Errorf("expected redirect to /app, got %q", rr.Header().Get("Location"))
	}
}

func TestPassword_ValidCredentials_EmailAsUsername(t *testing.T) {
	db := newTestDB(map[string]map[string]any{
		"users": {
			"u1": makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user"),
		},
	})
	a := newTestAuthorizer(db, WithEmailAsUsername())

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"k@example.com","password":"hunter2"}`))

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestPassword_InvalidPassword(t *testing.T) {
	db := newTestDB(map[string]map[string]any{
		"users": {
			"u1": makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user"),
		},
	})
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"wrongpass"}`))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestPassword_UserNotFound(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"ghost","password":"hunter2"}`))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestPassword_MissingBody(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPassword_WrongHTTPMethod(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, req)

	// GET with no handler match → ErrUnhandleableRequest → 500
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for unroutable GET, got %d", rr.Code)
	}
}

func TestPassword_PasswordlessUserDenied(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", authentication.PasswordInfo{Passwordless: true}, "user")

	db := newTestDB(map[string]map[string]any{
		"users": {"u1": user},
	})
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for passwordless user, got %d", rr.Code)
	}
}

func TestPassword_EmailNotVerified(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user")
	user["needsEmailVerification"] = true
	user["emailVerified"] = false

	db := newTestDB(map[string]map[string]any{
		"users": {"u1": user},
	})
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 for unverified email, got %d", rr.Code)
	}

	// verify the new location is /verify..
	if location := rr.Header().Get("Location"); location != "/app" {
		t.Errorf("expected /app location, got %s", location)
	}
}

func TestPassword_OAuthRestrictedSourceMismatch(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user")
	user["oauth"] = "google" // restricted to google, not password

	db := newTestDB(map[string]map[string]any{
		"users": {"u1": user},
	})
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for oauth-restricted user, got %d", rr.Code)
	}
}

func TestPassword_OAuthRestrictedSourceMatches(t *testing.T) {
	user := makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user")
	user["oauth"] = "password" // matches handler name

	db := newTestDB(map[string]map[string]any{
		"users": {"u1": user},
	})
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestPassword_EnumerationProtectionDefault(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"ghost","password":"hunter2"}`))

	body := rr.Body.String()
	if strings.Contains(body, "ghost") {
		t.Errorf("enumeration protection: response should not reveal username, got: %s", body)
	}
	if !strings.Contains(body, "incorrect") {
		t.Errorf("expected generic error message, got: %s", body)
	}
}

func TestPassword_EnumerationProtectionDisabled(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db, DisableUserEnumerationProtection())

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"ghost","password":"hunter2"}`))

	body := rr.Body.String()
	if !strings.Contains(body, "not found") {
		t.Errorf("expected specific error with enumeration disabled, got: %s", body)
	}
}
