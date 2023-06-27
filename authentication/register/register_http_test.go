package register

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
)

func TestMain(m *testing.M) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}

	m.Run()

	roles.AVAILABLE_ROLES = map[string][]string{}
}

func TestRegistrationHandlerMissingRole(t *testing.T) {
	handler := RegistrationHandler{}

	// Test Missing Username
	payload := `{"username": "kvizdos", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestRegistrationHandlerInvalidRole(t *testing.T) {
	handler := RegistrationHandler{
		DefaultRoleName: "not-set",
	}

	// Test Missing Username
	payload := `{"username": "kvizdos", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestRegistrationHandlerMissingBodyParams(t *testing.T) {
	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	// Test Missing Username
	payload := `{"password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}

	// Test Missing Password
	payload = `{"username": "kenton"}`

	req, err = http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code (missing password): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestRegistrationHandlerUsernameTaken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "kenton", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusConflict)
	}
}

func TestRegistrationHandlerSuccess(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "kenton", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
		return
	}

	newUser, found := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	if !found {
		t.Errorf("new user did not get added to database.")
		return
	}

	user := newUser.(map[string]interface{})
	passwordInfo := user["password"].(authentication.PasswordInfo)
	if passwordInfo.Password == "password123" {
		t.Errorf("PASSWORD NOT ENCRYPTED!")
		return
	}

	if len(passwordInfo.Salt) != 32 {
		t.Errorf("Salt not a correct length, expected %d, got %d", 32, len(passwordInfo.Salt))
	}
}

// Tests to see if the username does not match [A-z0-9]
func TestRegistrationHandlerInvalidUsername(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "<i want xss>!", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
		return
	}
}
