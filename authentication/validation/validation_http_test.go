package validation

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
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

func InjectTokenToDatabase(db database.DatabaseAccessor) string {
	token, _ := authentication.GenerateRandomString(64)

	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedCode := hasher.Sum(nil)
	hashedToken := fmt.Sprintf("%x", hashedCode)

	db.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": authentication.PasswordSession{
				Token:     hashedToken,
				ExpiresAt: time.Now().Unix() + 60000,
			},
		},
	})

	return token
}

func TestValidateHTTPNoCookiePresent(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	req, err := http.NewRequest("GET", "/app", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateHTTPUserToken(req, testDb)

	if err == nil {
		fmt.Printf("expected to receive an error")
		return
	}
}

func TestValidationHTTPMalformedTokenNotBase64Encoded(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

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

	_, err = ValidateHTTPUserToken(req, testDb)

	if err == nil {
		t.Error("expected error")
		return
	}

	if err.Error() != "token could not be parsed" {
		t.Errorf("Received incorrect error: %s", err.Error())
	}
}

func TestValidationHTTPalformedTokenBase64EncodedInvalidTokenLength(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

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

	_, err = ValidateHTTPUserToken(req, testDb)

	if err.Error() != "token could not be parsed" {
		t.Errorf("Received incorrect error: %s", err.Error())
	}
}

func TestValidationMiddlewareInvalidTokenUserDoesNotExist(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

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

	_, err = ValidateHTTPUserToken(req, testDb)

	if err.Error() != "token could not be validated" {
		t.Errorf("Received incorrect error: %s", err.Error())
	}
}

func TestValidationMiddlewareInvalidTokenBadUsername(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

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

	_, err = ValidateHTTPUserToken(req, testDb)

	if err.Error() != "token could not be validated" {
		t.Errorf("Received incorrect error: %s", err.Error())
	}
}

func TestValidationMiddlewareInvalidTokenBadToken(t *testing.T) {
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

	InjectTokenToDatabase(testDb)

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

	_, err = ValidateHTTPUserToken(req, testDb)

	if err.Error() != "token did not validate" {
		t.Errorf("Received incorrect error: %s", err.Error())
	}
}

func TestValidationMiddlewareValidToken(t *testing.T) {
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

	validToken := InjectTokenToDatabase(testDb)

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

	user, err := ValidateHTTPUserToken(req, testDb)

	if err != nil {
		t.Errorf("received unexpected error: %s", err.Error())
		return
	}

	if user.Username != "kenton" {
		t.Errorf("received invalid username")
	}
}

func TestValidationMiddlewareExpiredToken(t *testing.T) {
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

	_, err = ValidateHTTPUserToken(req, testDb)

	if err.Error() != "token did not validate" {
		t.Errorf("Received incorrect error: %s", err.Error())
	}

	// Validate that expired token was removed from database
	dbUser, _ := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	var tmpUser users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

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
					"email":    "email@email.com",
					"sessions": []interface{}{},
					"role":     "user",
				},
			},
		},
	}

	tokenString := InjectTokenToDatabase(testDb)

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

	_, err = ValidateHTTPUserToken(req, testDb)

	if err != nil {
		t.Errorf("Received error: %s", err.Error())
	}

	// Validate that expired token was removed from database
	dbUser, _ := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	var tmpUser users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

	if len(user.PasswordSessions) != 1 {
		t.Errorf("expired token was not removed from database")
	}
}
