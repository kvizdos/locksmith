package validation

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

type testHandler struct{}

func (lh testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK!"))
}

func injectTokenToDatabase(db database.DatabaseAccessor) string {
	token, _ := authentication.GenerateRandomString(64)
	db.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": authentication.PasswordSession{
				Token:     token,
				ExpiresAt: time.Now().Unix() + 60000,
			},
		},
	})

	return token
}

func TestValidationMiddlewareNoCookiePresent(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code (no cookie): got %v, want %v", status, http.StatusSeeOther)
	}

	url, _ := rr.Result().Location()

	if url.Path != "/login" {
		t.Errorf("did not redirect to /login: %s", url.Path)
	}
}

func TestValidationMiddlewareMalformedTokenNotBase64Encoded(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject Token..
	token := http.Cookie{
		Name:     "token",
		Value:    "randomtoken",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code (invalid token): got %v, want %v", status, http.StatusSeeOther)
	}

	url, _ := rr.Result().Location()

	if url == nil || url.Path != "/login" {
		if url != nil {
			t.Errorf("did not redirect to /login: %s", url.Path)
		} else {
			t.Errorf("response URL is nil")
		}
	}
}

func TestValidationMiddlewareMalformedTokenBase64EncodedInvalidTokenLength(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	tok := base64.StdEncoding.EncodeToString([]byte("abc:123"))

	// Inject Token..
	token := http.Cookie{
		Name:     "token",
		Value:    tok,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code (invalid token): got %v, want %v", status, http.StatusSeeOther)
	}

	url, _ := rr.Result().Location()

	if url == nil || url.Path != "/login" {
		if url != nil {
			t.Errorf("did not redirect to /login: %s", url.Path)
		} else {
			t.Errorf("response URL is nil")
		}
	}
}

func TestValidationMiddlewareMalformedTokenBase64EncodedValidLength(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	randomness, _ := authentication.GenerateRandomString(64)
	tok := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:user", randomness)))

	// Inject Token..
	token := http.Cookie{
		Name:     "token",
		Value:    tok,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code (invalid token): got %v, want %v", status, http.StatusSeeOther)
	}

	url, _ := rr.Result().Location()

	if url == nil || url.Path != "/login" {
		if url != nil {
			t.Errorf("did not redirect to /login: %s", url.Path)
		} else {
			t.Errorf("response URL is nil")
		}
	}
}

func TestValidationMiddlewareInvalidTokenBadUsername(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject Token..
	u := users.LocksmithUser{
		Username: "kenton",
	}

	randToken, _ := authentication.GenerateRandomString(64)
	session := authentication.PasswordSession{
		Token:     randToken,
		ExpiresAt: time.Now().Unix(),
	}

	token := http.Cookie{
		Name:     "token",
		Value:    u.GenerateCookieValueFromSession(session),
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code (invalid username): got %v, want %v", status, http.StatusSeeOther)
	}

	url, _ := rr.Result().Location()

	if url == nil || url.Path != "/login" {
		if url != nil {
			t.Errorf("did not redirect to /login: %s", url.Path)
		} else {
			t.Errorf("response URL is nil")
		}
	}
}

func TestValidationMiddlewareInvalidTokenBadToken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
				},
			},
		},
	}

	injectTokenToDatabase(testDb)

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject Token..
	u := users.LocksmithUser{
		Username: "kenton",
	}

	randToken, _ := authentication.GenerateRandomString(64)

	session := authentication.PasswordSession{
		Token:     randToken,
		ExpiresAt: time.Now().Unix(),
	}

	v := u.GenerateCookieValueFromSession(session)

	token := http.Cookie{
		Name:     "token",
		Value:    v,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code (invalid username): got %v, want %v", status, http.StatusSeeOther)
	}

	url, _ := rr.Result().Location()

	if url == nil || url.Path != "/login" {
		if url != nil {
			t.Errorf("did not redirect to /login: %s", url.Path)
		} else {
			t.Errorf("response URL is nil")
		}
	}
}

func TestValidationMiddlewareValidToken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
				},
			},
		},
	}

	validToken := injectTokenToDatabase(testDb)

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject Token..
	u := users.LocksmithUser{
		Username: "kenton",
	}

	session := authentication.PasswordSession{
		Token:     validToken,
		ExpiresAt: time.Now().Unix(),
	}

	v := u.GenerateCookieValueFromSession(session)

	token := http.Cookie{
		Name:     "token",
		Value:    v,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code (valid token): got %v, want %v", status, http.StatusOK)
	}

	url, _ := rr.Result().Location()

	if url != nil {
		t.Errorf("unexpected redirect: %s", url.Path)
	}
}

func TestValidationMiddlewareExpiredToken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
				},
			},
		},
	}

	tokenString, _ := authentication.GenerateRandomString(64)
	testDb.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": authentication.PasswordSession{
				Token:     tokenString,
				ExpiresAt: time.Now().Unix() - 5000,
			},
		},
	})

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject Token..
	u := users.LocksmithUser{
		Username: "kenton",
	}

	session := authentication.PasswordSession{
		Token:     tokenString,
		ExpiresAt: time.Now().Unix(),
	}

	v := u.GenerateCookieValueFromSession(session)

	token := http.Cookie{
		Name:     "token",
		Value:    v,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("unexpected status code (expired token): got %v, want %v", status, http.StatusSeeOther)
	}

	url, _ := rr.Result().Location()

	if url == nil || url.Path != "/login" {
		if url != nil {
			t.Errorf("did not redirect to /login: %s", url.Path)
		} else {
			t.Errorf("response URL is nil")
		}
	}

	// Validate that expired token was removed from database
	dbUser, _ := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})
	user := users.LocksmithUserFromMap(dbUser.(map[string]interface{}))

	if len(user.PasswordSessions) != 0 {
		t.Errorf("expired token was not removed from database")
	}
}

func TestValidationMiddlewareRemovesExpiredTokenAndPreservesValid(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
				},
			},
		},
	}

	tokenString, _ := authentication.GenerateRandomString(64)
	testDb.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": authentication.PasswordSession{
				Token:     tokenString,
				ExpiresAt: time.Now().Unix() - 5000,
			},
		},
	})

	tokenString, _ = authentication.GenerateRandomString(64)
	testDb.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": authentication.PasswordSession{
				Token:     tokenString,
				ExpiresAt: time.Now().Unix() + 50000,
			},
		},
	})

	handler := testHandler{}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject Token..
	u := users.LocksmithUser{
		Username: "kenton",
	}

	session := authentication.PasswordSession{
		Token:     tokenString,
		ExpiresAt: time.Now().Unix(),
	}

	v := u.GenerateCookieValueFromSession(session)

	token := http.Cookie{
		Name:     "token",
		Value:    v,
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	req.AddCookie(&token)

	rr := httptest.NewRecorder()
	middleware := ValidateUserToken(handler, testDb)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code (valid token): got %v, want %v", status, http.StatusOK)
	}

	url, _ := rr.Result().Location()

	if url != nil {
		t.Errorf("unexpected redirect: %s", url.Path)
	}

	// Validate that expired token was removed from database
	dbUser, _ := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})
	user := users.LocksmithUserFromMap(dbUser.(map[string]interface{}))

	if len(user.PasswordSessions) != 1 {
		t.Errorf("expired token was not removed from database")
	}
}
