package validation

import (
	"testing"
	"time"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
)

func pushSession(db database.DatabaseAccessor) {
	db.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": authentication.PasswordSession{
				Token:     "correct-token",
				ExpiresAt: time.Now().Unix() + 60000,
			},
		},
	})
}

func TestValidateInvalidUsername(t *testing.T) {
	testToken := authentication.Token{
		Token:    "invalid-token",
		Username: "jimbob",
	}

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

	_, _, err := ValidateToken(testToken, testDb)

	if err == nil {
		t.Errorf("expected error")
		return
	}

	if err.Error() != "invalid username" {
		t.Errorf("error should be invalid username")
		return
	}
}

func TestValidateInvalidToken(t *testing.T) {
	testToken := authentication.Token{
		Token:    "invalid-token",
		Username: "kenton",
	}

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

	// Push a session..
	pushSession(testDb)

	_, validated, err := ValidateToken(testToken, testDb)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if validated {
		t.Errorf("token should not be valid.")
	}
}

func TestValidateValidToken(t *testing.T) {
	testToken := authentication.Token{
		Token:    "correct-token",
		Username: "kenton",
	}

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

	// Push a session..
	pushSession(testDb)

	_, validated, err := ValidateToken(testToken, testDb)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if !validated {
		t.Errorf("token should be valid.")
	}
}
