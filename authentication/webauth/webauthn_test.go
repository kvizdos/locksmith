package webauth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func SetupTesting(t *testing.T) {
	t.Helper()
	t.Setenv("AppName", "Locksmith Test App")
	t.Setenv("DomainName", "locksmith.local")
	t.Setenv("AuthOrigin", "locksmith.local")

	InitializeWebAuthn()
}

func TestBeginRegister(t *testing.T) {
	SetupTesting(t)

	passwordInfo, _ := authentication.CompileLocksmithPassword("helloworld", "salt")
	testUser := users.LocksmithUser{
		ID:           uuid.New().String(),
		Username:     "kvizdos",
		PasswordInfo: passwordInfo,
	}
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":          "c8531661-22a7-493f-b228-028842e09a05",
					"username":    "kvizdos",
					"websessions": []interface{}{},
				},
			},
		},
	}
	_, err := BeginRegisterWebAuthn(testUser, testDb)

	if err != nil {
		t.Errorf("failed to create credentail: %s", err.Error())
		return
	}
}
