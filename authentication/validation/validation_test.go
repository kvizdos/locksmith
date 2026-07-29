package validation

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
)

func TestMain(m *testing.M) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
		"user": {
			"view.admin",
			"user.delete.self",
		},
	}

	m.Run()

	roles.AVAILABLE_ROLES = map[string][]string{}
}

func pushSession(db database.DatabaseAccessor) {
	hasher := sha256.New()
	hasher.Write([]byte("correct-token"))
	hashedCode := hasher.Sum(nil)
	hashedToken := fmt.Sprintf("%x", hashedCode)

	db.UpdateOne("users", map[string]any{
		"username": "kenton",
	}, map[database.DatabaseUpdateActions]map[string]any{
		database.PUSH: {
			"sessions": authentication.PasswordSession{
				Token:     hashedToken,
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
		Tables: map[string]map[string]any{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]any{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []any{},
					"role":     "user",
				},
			},
		},
	}

	_, _, err := ValidateToken(testToken, testDb, "")

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
		Username: "c8531661-22a7-493f-b228-028842e09a05",
	}

	testDb := database.TestDatabase{
		Tables: map[string]map[string]any{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]any{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []any{},
					"role":     "user",
				},
			},
		},
	}

	// Push a session..
	pushSession(testDb)

	_, validated, err := ValidateToken(testToken, testDb, "")

	if err != nil {
		t.Error(err)
		return
	}

	if validated {
		t.Errorf("token should not be valid.")
	}
}

func TestValidateValidToken(t *testing.T) {
	testToken := authentication.Token{
		Token:    "correct-token",
		Username: "c8531661-22a7-493f-b228-028842e09a05",
	}

	testDb := database.TestDatabase{
		Tables: map[string]map[string]any{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]any{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []any{},
					"role":     "user",
				},
			},
		},
	}

	// Push a session..
	pushSession(testDb)

	_, validated, err := ValidateToken(testToken, testDb, "")

	if err != nil {
		t.Error(err)
		return
	}

	if !validated {
		t.Errorf("token should be valid.")
	}
}
