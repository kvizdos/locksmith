package invitations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
	"kv.codes/locksmith/users"
)

func TestMain(m *testing.M) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"user.invite",
		},
		"user": {},
	}

	m.Run()

	roles.AVAILABLE_ROLES = map[string][]string{}
}

func TestInviteUserHTTPInvalidMethod(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users":   {},
			"invites": {},
		},
	}

	handler := AdministrationInviteUserHandler{}

	req, err := http.NewRequest("GET", "/api/invite", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestInviteUserHTTPAuthUserNotPassed(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users":   {},
			"invites": {},
		},
	}

	handler := AdministrationInviteUserHandler{}

	req, err := http.NewRequest("POST", "/api/invite", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestInviteUserHTTPInvalidPayload(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users":   {},
			"invites": {},
		},
	}

	handler := AdministrationInviteUserHandler{}

	payload := `{"email": "random@random.com"}`
	req, err := http.NewRequest("POST", "/api/invite", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	withUser := users.LocksmithUser{
		Username: "kenton",
		Role:     "admin",
	}
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestInviteUserHTTPAlreadyInvited(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
			"invites": {
				"ID": map[string]interface{}{
					"email": "random@random.com",
				},
			},
		},
	}

	handler := AdministrationInviteUserHandler{}

	payload := `{"email": "random@random.com", "role": "user"}`
	req, err := http.NewRequest("POST", "/api/invite", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	withUser := users.LocksmithUser{
		ID:       "user-id",
		Username: "kenton",
		Role:     "admin",
	}
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusConflict)
	}
}

func TestInviteUserHTTPAlreadyRegistered(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"ID": map[string]interface{}{
					"email": "random@random.com",
				},
			},
			"invites": {},
		},
	}

	handler := AdministrationInviteUserHandler{}

	payload := `{"email": "random@random.com", "role": "user"}`
	req, err := http.NewRequest("POST", "/api/invite", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	withUser := users.LocksmithUser{
		ID:       "user-id",
		Username: "kenton",
		Role:     "admin",
	}
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusConflict)
	}
}

func TestInviteUserHTTPSuccess(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users":   {},
			"invites": {},
		},
	}

	handler := AdministrationInviteUserHandler{}

	payload := `{"email": "random@random.com", "role": "user"}`
	req, err := http.NewRequest("POST", "/api/invite", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	withUser := users.LocksmithUser{
		ID:       "user-id",
		Username: "kenton",
		Role:     "admin",
	}
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusOK)
	}
}
