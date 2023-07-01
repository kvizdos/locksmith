package invitations

import (
	"testing"

	"kv.codes/locksmith/database"
)

func TestInviteUserInvalidRole(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	_, err := InviteUser(testDb, "email@email.com", "fakerole")

	if err.Error() != "invalid role" {
		t.Errorf("received unexpected error message: %s", err.Error())
	}
}

func TestInviteUserInvalidEmail(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	_, err := InviteUser(testDb, "malformed.email@", "user")

	if err.Error() != "invalid email address" {
		t.Errorf("received unexpected error message: %s", err.Error())
	}
}

func TestInviteUserEmailExistsAsRegisteredUser(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"password": "password",
					"sessions": []interface{}{},
					"role":     "user",
				},
			},
		},
	}

	_, err := InviteUser(testDb, "email@email.com", "user")

	if err.Error() != "email already registered" {
		t.Errorf("received unexpected error message: %s", err.Error())
	}
}

func TestInviteUserEmailExistsAsInvite(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {
				"id": map[string]interface{}{
					"code":  "invite-token",
					"email": "new@email.com",
				},
			},
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"password": "password",
					"sessions": []interface{}{},
					"role":     "user",
				},
			},
		},
	}

	_, err := InviteUser(testDb, "new@email.com", "user")

	if err.Error() != "email already invited" {
		t.Errorf("received unexpected error message: %s", err.Error())
	}
}

func TestInviteUserSuccess(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {
				"id": map[string]interface{}{
					"code":  "old-invite-token",
					"email": "old@email.com",
				},
			},
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"password": "password",
					"sessions": []interface{}{},
					"role":     "user",
				},
			},
		},
	}

	_, err := InviteUser(testDb, "new@email.com", "user")

	if err != nil {
		t.Errorf("received unexpected error message: %s", err.Error())
		return
	}
}
