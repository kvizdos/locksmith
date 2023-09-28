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

	_, _, err := InviteUser(testDb, "email@email.com", "fakerole", "random-uid")

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

	_, _, err := InviteUser(testDb, "malformed.email@", "user", "random-uid")

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

	_, _, err := InviteUser(testDb, "email@email.com", "user", "random-uid")

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

	_, _, err := InviteUser(testDb, "new@email.com", "user", "random-uid")

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

	_, _, err := InviteUser(testDb, "new@email.com", "user", "random-uid")

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
					"userid":  "abc123",
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

	if invite.AttachUserID != "abc123" {
		t.Errorf("recevied incorrect attach user id: %s", invite.AttachUserID)
	}
}

func TestReinviteUserNoTokenFound(t *testing.T) {
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
					"userid":  "abc123",
				},
			},
			"users": {},
		},
	}

	_, err := ReinviteUser(testDb, "invalid-user-id", "auth-user-id")

	if err == nil {
		t.Errorf("expected to receive an error message!")
		return
	}

	if err.Error() != "could not find invite" {
		t.Errorf("got wrong error mesage: %s", err)
		return
	}
}

func TestReinviteUserNoEmailChangeSuccess(t *testing.T) {
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
					"sentAt":  int64(0),
					"code":    fmt.Sprintf("%x", hashedCode),
					"userid":  "abc123",
				},
			},
			"users": {},
		},
	}

	code, err := ReinviteUser(testDb, "abc123", "b-uuid")

	if err != nil {
		t.Errorf("should not have received an error message: %s", err)
		return
	}

	rawInvite, found := testDb.FindOne("invites", map[string]interface{}{
		"userid": "abc123",
	})

	if !found {
		fmt.Errorf("could not find the invite in the database")
		return
	}

	invite := rawInvite.(map[string]interface{})

	hasher = sha256.New()
	hasher.Write([]byte(code))
	hahedReceivedCode := hasher.Sum(nil)

	if invite["code"].(string) != fmt.Sprintf("%x", hahedReceivedCode) {
		t.Errorf("received incorrect code: %s", invite["code"].(string))
	}

	if invite["inviter"].(string) != "b-uuid" {
		t.Errorf("inviter did not change: %s", invite["inviter"].(string))
	}

	if invite["sentAt"].(int64) == 0 {
		t.Error("sentAt time did not change")
	}

	if invite["email"].(string) != "kvizdos@email.com" {
		t.Error("email is incorrect: %", invite["email"].(string))
	}

	if invite["role"].(string) != "user" {
		t.Error("role is incorrect: %", invite["role"].(string))
	}
}

func TestReinviteWithEmailChangeEmailAlreadyRegistered(t *testing.T) {
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
					"userid":  "abc123",
				},
			},
			"users": {
				"id": map[string]interface{}{
					"email": "an-email@example.com",
				},
			},
		},
	}

	_, err := ReinviteUser(testDb, "abc123", "auth-user-id", "an-email@example.com")

	if err == nil {
		t.Errorf("expected to receive an error message!")
		return
	}

	if err.Error() != "email already registered" {
		t.Errorf("got wrong error mesage: %s", err)
		return
	}
}

func TestReinviteWithEmailChangeEmailAlreadyInvited(t *testing.T) {
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
					"userid":  "abc123",
				},
				"id2": map[string]interface{}{
					"email":   "kvizdos2@gmail.com",
					"role":    "user",
					"inviter": "a-uuid",
					"sentAt":  time.Now().Unix(),
					"code":    fmt.Sprintf("%x", hashedCode),
					"userid":  "abc123",
				},
			},
			"users": {
				"id": map[string]interface{}{
					"email": "an-email@example.com",
				},
			},
		},
	}

	_, err := ReinviteUser(testDb, "abc123", "auth-user-id", "kvizdos2@gmail.com")

	if err == nil {
		t.Errorf("expected to receive an error message!")
		return
	}

	if err.Error() != "email already invited" {
		t.Errorf("got wrong error mesage: %s", err)
		return
	}
}

func TestReinviteWithEmailChangeEmailSuccess(t *testing.T) {
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
					"userid":  "abc123",
				},
				"id2": map[string]interface{}{
					"email":   "kvizdos2@gmail.com",
					"role":    "user",
					"inviter": "a-uuid",
					"sentAt":  time.Now().Unix(),
					"code":    fmt.Sprintf("%x", hashedCode),
					"userid":  "xyz123",
				},
			},
			"users": {
				"id": map[string]interface{}{
					"email": "an-email@example.com",
				},
			},
		},
	}

	code, err := ReinviteUser(testDb, "abc123", "b-uuid", "new-email@example.com")

	if err != nil {
		t.Errorf("should not have received an error message: %s", err)
		return
	}

	rawInvite, found := testDb.FindOne("invites", map[string]interface{}{
		"userid": "abc123",
	})

	if !found {
		fmt.Errorf("could not find the invite in the database")
		return
	}

	invite := rawInvite.(map[string]interface{})
	hasher = sha256.New()
	hasher.Write([]byte(code))
	hahedReceivedCode := hasher.Sum(nil)

	if invite["code"].(string) != fmt.Sprintf("%x", hahedReceivedCode) {
		t.Errorf("received incorrect code: %s (expected %s)", invite["code"].(string), code)
	}

	if invite["inviter"].(string) != "b-uuid" {
		t.Errorf("inviter did not change: %s", invite["inviter"].(string))
	}

	if invite["sentAt"].(int64) == 0 {
		t.Error("sentAt time did not change")
	}

	if invite["email"].(string) != "new-email@example.com" {
		t.Errorf("email is incorrect: %s", invite["email"].(string))
	}

	if invite["role"].(string) != "user" {
		t.Error("role is incorrect: %", invite["role"].(string))
	}
}
