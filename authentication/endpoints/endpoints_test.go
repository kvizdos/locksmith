package endpoints

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/users"
)

type testHandler struct{}

func (lh testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authUser := r.Context().Value("authUser").(users.LocksmithUserInterface)
	role, _ := authUser.GetRole()
	w.Write([]byte(fmt.Sprintf("%s %d", role.Name, len(role.Permissions))))
}

func InjectTokenToDatabase(db database.DatabaseAccessor) string {
	u := users.LocksmithUser{
		ID:       "c8531661-22a7-493f-b228-028842e09a05",
		Username: "kenton",
	}

	token, _ := authentication.GenerateRandomString(64)
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedCode := hasher.Sum(nil)
	hashedToken := fmt.Sprintf("%x", hashedCode)

	session := authentication.PasswordSession{
		Token:     hashedToken,
		ExpiresAt: time.Now().Unix() + 60000,
	}
	db.UpdateOne("users", map[string]interface{}{
		"username": u.Username,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": session,
		},
	})

	session.Token = token

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
	roles.AVAILABLE_ROLES = map[string][]string{
		"user": {},
	}
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
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
					"email":    "email@email.com",
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
					"email":    "email@email.com",
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

func TestSecureEndpointHTTPMiddlewareSecondaryValidationFails(t *testing.T) {
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
					"email":    "email@email.com",
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
		SecondaryValidation: func(lui users.LocksmithUserInterface, da database.DatabaseAccessor) int {
			return 401
		},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusUnauthorized)
	}
}

func TestSecureEndpointHTTPMiddlewareSecondaryValidationSucceeds(t *testing.T) {
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
					"email":    "email@email.com",
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
		SecondaryValidation: func(lui users.LocksmithUserInterface, da database.DatabaseAccessor) int {
			return 200
		},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}
}

type customUserInterface interface {
	users.LocksmithUserInterface
}

type customUser struct {
	users.LocksmithUser

	customObject string
}

func (c customUser) ReadFromMap(writeTo *users.LocksmithUserInterface, u map[string]interface{}) {
	// Load initial locksmith data
	var user users.LocksmithUserInterface
	c.LocksmithUser.ReadFromMap(&user, u)
	lsu := user.(users.LocksmithUser)

	converted := customUser{
		LocksmithUser: lsu,
	}

	converted.customObject = u["customObject"].(string)

	*writeTo = converted
}

func TestSecureEndpointHTTPMiddlewareSecondaryValidationSucceedsWithCustomUser(t *testing.T) {
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
					"id":           "c8531661-22a7-493f-b228-028842e09a05",
					"username":     "kenton",
					"email":        "custom@email.com",
					"sessions":     []interface{}{},
					"role":         "admin",
					"customObject": "hello",
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
		CustomUser: customUser{},
		SecondaryValidation: func(lui users.LocksmithUserInterface, da database.DatabaseAccessor) int {
			user := lui.(customUser)

			if user.customObject != "hello" {
				return 400
			}

			return 200
		},
	}
	middleware := SecureEndpointHTTPMiddleware(testHandler{}, testDb, opts)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}
}
