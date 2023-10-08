package endpoints

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
)

func TestSecureEndpointHTTPMiddlewareValidationFailsWithBadMagicToken(t *testing.T) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}

	sp, _ := signing.CreateSigningPackage()
	magic.MagicSigningPackage = &sp

	mac, macIdentifier, _ := magic.CreateMagicAuthentication(magic.MagicAuthenticationVariables{
		ForUserID: "c8531661-22a7-493f-b228-028842e09a05",
		AllowedPermissions: []string{
			"magic.test",
		},
		TTL: time.Hour,
	})
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []interface{}{},
					"role":     "admin",
					"magic": magic.MagicAuthentications{
						mac,
					}.ToMap(),
				},
			},
		},
	}

	macIdentifier = macIdentifier[:len(macIdentifier)-1] + "s"

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/example?magic=%s", macIdentifier), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	opts := EndpointSecurityOptions{
		MinimalPermissions: []string{
			"magic.test",
		},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusSeeOther)
	}
}

func TestSecureEndpointHTTPMiddlewareValidationFailsWithInvalidatedMagicToken(t *testing.T) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}

	sp, _ := signing.CreateSigningPackage()
	magic.MagicSigningPackage = &sp

	_, macIdentifier, _ := magic.CreateMagicAuthentication(magic.MagicAuthenticationVariables{
		ForUserID: "c8531661-22a7-493f-b228-028842e09a05",
		AllowedPermissions: []string{
			"magic.test",
		},
		TTL: time.Hour,
	})
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []interface{}{},
					"role":     "admin",
					"magic":    magic.MagicAuthentications{}.ToMap(),
				},
			},
		},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/example?magic=%s", macIdentifier), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	opts := EndpointSecurityOptions{
		MinimalPermissions: []string{
			"magic.test",
		},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusSeeOther)
	}
}

func TestSecureEndpointHTTPMiddlewareValidationFailsWithValidMagicTokenButNoMatchingPermissions(t *testing.T) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}

	sp, _ := signing.CreateSigningPackage()
	magic.MagicSigningPackage = &sp

	_, macIdentifier, _ := magic.CreateMagicAuthentication(magic.MagicAuthenticationVariables{
		ForUserID: "c8531661-22a7-493f-b228-028842e09a05",
		AllowedPermissions: []string{
			"magic.fails",
		},
		TTL: time.Hour,
	})
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []interface{}{},
					"role":     "admin",
					"magic":    magic.MagicAuthentications{}.ToMap(),
				},
			},
		},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/example?magic=%s", macIdentifier), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	opts := EndpointSecurityOptions{
		MinimalPermissions: []string{
			"magic.test",
		},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusSeeOther)
	}
}

func TestSecureEndpointHTTPMiddlewareValidationSucceedsWithMagicToken(t *testing.T) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}

	sp, _ := signing.CreateSigningPackage()
	magic.MagicSigningPackage = &sp

	mac, macIdentifier, _ := magic.CreateMagicAuthentication(magic.MagicAuthenticationVariables{
		ForUserID: "c8531661-22a7-493f-b228-028842e09a05",
		AllowedPermissions: []string{
			"magic.test",
		},
		TTL: time.Hour,
	})
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []interface{}{},
					"role":     "admin",
					"magic": magic.MagicAuthentications{
						mac,
					}.ToMap(),
				},
			},
		},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/example?magic=%s", macIdentifier), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	opts := EndpointSecurityOptions{
		MinimalPermissions: []string{
			"magic.test",
		},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	if rr.Body.String() != "admin 1" {
		t.Errorf("got a bad response: %s", rr.Body.String())
	}
}
