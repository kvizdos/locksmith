package login

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

func TestLoginHandlerMissingBodyParams(t *testing.T) {
	handler := LoginHandler{}

	// Test Missing Username
	payload := `{"password": "password123"}`

	req, err := http.NewRequest("POST", "/api/login", strings.NewReader(payload))
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

	req, err = http.NewRequest("POST", "/api/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code (missing password): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestLoginHandlerInvalidUsername(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := LoginHandler{}

	payload := `{"username": "kenton", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestLoginHandlerInvalidPassword(t *testing.T) {
	testPassword, _ := authentication.CompileLocksmithPassword("securepassword123")
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"password": testPassword,
					"sessions": []interface{}{},
				},
			},
		},
	}

	handler := LoginHandler{}

	payload := `{"username": "kenton", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("unexpected status code (invalid password): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestLoginHandlerValidPassword(t *testing.T) {
	testPassword, _ := authentication.CompileLocksmithPassword("securepassword123")
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"password": testPassword,
					"sessions": []interface{}{},
				},
			},
		},
	}

	handler := LoginHandler{}

	payload := `{"username": "kenton", "password": "securepassword123"}`

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

	// Validate Session Token looks good on receivers end..
	cookies := rr.Result().Cookies()

	if len(cookies) == 0 {
		t.Errorf("no cookies sent.")
		return
	}

	tokenCookie := cookies[0]

	if !tokenCookie.HttpOnly {
		t.Errorf("Token is not HttpOnly (security issue)!")
		return
	}

	if !tokenCookie.Secure {
		t.Errorf("Token is not Secure only (security issue)!")
		return
	}

	if tokenCookie.Expires.Unix() < time.Now().Unix() {
		t.Errorf("Token expiration incorrect")
		return
	}

	if len(tokenCookie.Value) == 0 {
		t.Errorf("No token attached to cookie.")
		return
	}

	decodedCookie, err := base64.StdEncoding.DecodeString(tokenCookie.Value)

	if err != nil {
		t.Errorf("token not base64 encoded: %s", err.Error())
		return
	}

	splitValue := strings.Split(string(decodedCookie), ":")

	token := splitValue[0]
	username := splitValue[1]

	if len(token) != 64 {
		t.Errorf("invalid token length, expected %d got %d", 64, len(token))
		return
	}

	if username != "kenton" {
		t.Errorf("token username invalid, expected '%s' got '%s'", "kenton", username)
	}

	// Validate token exists in Database..
	dbUser, _ := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	var tmpUser users.LocksmithUserStruct
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

	if len(user.PasswordSessions) != 1 {
		t.Errorf("got %d sessions, expected %d", len(user.PasswordSessions), 1)
		return
	}

	if len(user.PasswordSessions[0].Token) != 64 {
		t.Errorf("got %d token length, expected %d", len(user.PasswordSessions[0].Token), 64)
		return
	}
}
