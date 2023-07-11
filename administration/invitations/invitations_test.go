package invitations

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/database"
)

func TestInviteUserInvalidRole(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	_, err := InviteUser(testDb, "email@email.com", "fakerole", "random-uid")

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

	_, err := InviteUser(testDb, "malformed.email@", "user", "random-uid")

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

	_, err := InviteUser(testDb, "email@email.com", "user", "random-uid")

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

	_, err := InviteUser(testDb, "new@email.com", "user", "random-uid")

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

	_, err := InviteUser(testDb, "new@email.com", "user", "random-uid")

	if err != nil {
		t.Errorf("received unexpected error message: %s", err.Error())
		return
	}
}

func TestGetInviteCodeMalformedToken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {},
			"users":   {},
		},
	}

	_, err := GetInviteFromCode(testDb, "too-short")

	if err.Error() != "invalid token length" {
		t.Errorf("received unexpected error message: %s", err.Error())
		return
	}
}

func TestGetInviteCodeInvalidToken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {
				"id": map[string]interface{}{
					"code": "jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M",
				},
			},
			"users": {},
		},
	}

	_, err := GetInviteFromCode(testDb, "BBBBBBBBB-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M")

	if err.Error() != "could not find token" {
		t.Errorf("received unexpected error message: %s", err.Error())
		return
	}
}

func TestGetInviteCodeValidToken(t *testing.T) {
	hasher := sha256.New()
	hasher.Write([]byte("jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M"))
	hashedCode := hasher.Sum(nil)

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {
				"id": map[string]interface{}{
					"email":   "kvizdos@email.com",
					"role":    "user",
					"inviter": "a-uuid",
					"sentAt":  time.Now().Unix(),
					"code":    fmt.Sprintf("%x", hashedCode),
				},
			},
			"users": {},
		},
	}

	invite, err := GetInviteFromCode(testDb, "jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M")

	if err != nil {
		t.Errorf("received unexpected error message: %s", err.Error())
		return
	}

	if invite.Email != "kvizdos@email.com" {
		t.Errorf("could not find correct email")
	}
}
