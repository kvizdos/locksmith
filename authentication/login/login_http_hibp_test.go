package login

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/database"
)

func TestLoginNotFoundInHIBPSucceeds(t *testing.T) {
	pass, _ := authentication.GenerateRandomString(64)
	testPassword, _ := authentication.CompileLocksmithPassword(pass)
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"password": testPassword,
					"sessions": []interface{}{},
					"role":     "user",
				},
			},
		},
	}

	handler := LoginHandler{
		HIBP: hibp.HIBPSettings{
			Enabled:                  true,
			AppName:                  "Locksmith Integration Tests",
			Enforcement:              hibp.STRICT,
			HTTPClient:               &http.Client{},
			PasswordSecurityInfoLink: "#",
		},
	}
	payload := fmt.Sprintf(`{"username": "kenton", "password": "%s"}`, pass)

	req, err := http.NewRequest("POST", "/api/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code (correct info): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestLoginFoundInHIBPDoesNotAuthenticateAndRedirects(t *testing.T) {
	pass := "password" // test with an insecure password.
	testPassword, _ := authentication.CompileLocksmithPassword(pass)
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"password": testPassword,
					"sessions": []interface{}{},
					"role":     "user",
				},
			},
		},
	}

	handler := LoginHandler{
		HIBP: hibp.HIBPSettings{
			Enabled:                  true,
			AppName:                  "Locksmith Integration Tests",
			Enforcement:              hibp.STRICT,
			HTTPClient:               &http.Client{},
			PasswordSecurityInfoLink: "#",
		},
	}
	payload := fmt.Sprintf(`{"username": "kenton", "password": "%s"}`, pass)

	req, err := http.NewRequest("POST", "/api/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("unexpected status code (correct info): got %v, want %v", status, http.StatusTemporaryRedirect)
	}

	loc, _ := rr.Result().Location()

	if loc.String() != "/reset-password?hibp=true" {
		t.Errorf("expected to be redirected to reset-password page! got: %s", loc.String())
	}
}
