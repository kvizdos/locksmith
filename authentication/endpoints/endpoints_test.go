package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
	"kv.codes/locksmith/users"
)

type testHandler struct{}

func (lh testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK!"))
}

func InjectTokenToDatabase(db database.DatabaseAccessor) string {
	u := users.LocksmithUser{
		Username: "kenton",
	}

	token, _ := authentication.GenerateRandomString(64)
	session := authentication.PasswordSession{
		Token:     token,
		ExpiresAt: time.Now().Unix() + 60000,
	}
	db.UpdateOne("users", map[string]interface{}{
		"username": u.Username,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": session,
		},
	})

	return u.GenerateCookieValueFromSession(session)
}

// Validator is tested in the validation package, so
// im only going to test one fail case here to make sure
func TestSecureEndpointHTTPMiddlewareInvalidToken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	req, err := http.NewRequest("POST", "/api/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestSecureEndpointHTTPMiddlewareInvalidPermissions(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
					"role":     "user",
				},
			},
		},
	}

	token := InjectTokenToDatabase(testDb)

	req, err := http.NewRequest("POST", "/api/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject cookie..
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()

	opts := EndpointSecurityOptions{
		MinimalPermissions: []string{"use.api"},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestSecureEndpointHTTPMiddlewareValidPermissions(t *testing.T) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
					"role":     "admin",
				},
			},
		},
	}

	token := InjectTokenToDatabase(testDb)

	req, err := http.NewRequest("POST", "/api/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject cookie..
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()

	opts := EndpointSecurityOptions{
		MinimalPermissions: []string{"view.admin"},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestSecureEndpointHTTPMiddlewareFailsMultipleRequiredPermissions(t *testing.T) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
					"role":     "admin",
				},
			},
		},
	}

	token := InjectTokenToDatabase(testDb)

	req, err := http.NewRequest("POST", "/api/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject cookie..
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()

	opts := EndpointSecurityOptions{
		MinimalPermissions: []string{"view.admin", "modify.admin", "user.delete.self"},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
	}
}
