package users

import (
	"testing"

	"kv.codes/locksmith/authentication"
)

func TestLoadLocksmithUserFromMap(t *testing.T) {
	var sessions []interface{}
	tempSessions := []authentication.PasswordSession{
		{
			Token:     "abc",
			ExpiresAt: 123,
		},
		{
			Token:     "bca",
			ExpiresAt: 321,
		},
	}

	for _, sess := range tempSessions {
		sessions = append(sessions, sess)
	}

	lsu := map[string]interface{}{
		"id":       "c8531661-22a7-493f-b228-028842e09a05",
		"username": "kenton",
		"password": map[string]interface{}{
			"password": "passwordhere",
			"salt":     "salthere",
		},
		"sessions": sessions,
	}

	var user LocksmithUserStruct
	LocksmithUser{}.ReadFromMap(&user, lsu)

	converted := user.(LocksmithUser)

	if converted.ID != "c8531661-22a7-493f-b228-028842e09a05" {
		t.Errorf("invalid id: %s\n", converted.ID)
	}

	if converted.Username != "kenton" {
		t.Errorf("invalid username: %s\n", converted.Username)
	}

	if converted.PasswordInfo.Password != "passwordhere" {
		t.Errorf("invalid password: %s\n", converted.PasswordInfo.Password)
	}

	if converted.PasswordInfo.Salt != "salthere" {
		t.Errorf("invalid salt: %s\n", converted.PasswordInfo.Salt)
	}

	if len(converted.PasswordSessions) != 2 {
		t.Errorf("invalid password session length, expected 2, got %d", len(converted.PasswordSessions))
	}
}

type customUserInterface interface {
	LocksmithUserStruct
}

type customUser struct {
	LocksmithUser

	customObject string
}

func (c customUser) ReadFromMap(writeTo *LocksmithUserStruct, u map[string]interface{}) {
	// Load initial locksmith data
	var user LocksmithUserStruct
	c.LocksmithUser.ReadFromMap(&user, u)
	lsu := user.(LocksmithUser)

	converted := customUser{
		LocksmithUser: lsu,
	}

	converted.customObject = u["customObject"].(string)

	*writeTo = converted
}

func TestLoadCustomUserFromMap(t *testing.T) {

	lsu := map[string]interface{}{
		"id":       "c8531661-22a7-493f-b228-028842e09a05",
		"username": "kenton",
		"password": map[string]interface{}{
			"password": "passwordhere",
			"salt":     "salthere",
		},
		"sessions":     []interface{}{},
		"customObject": "helloworld",
	}

	var user LocksmithUserStruct
	customUser{}.ReadFromMap(&user, lsu)
	converted := user.(customUser)

	if converted.ID != "c8531661-22a7-493f-b228-028842e09a05" {
		t.Errorf("invalid id: %s\n", converted.ID)
	}

	if converted.Username != "kenton" {
		t.Errorf("invalid username: %s\n", converted.Username)
	}

	if converted.PasswordInfo.Password != "passwordhere" {
		t.Errorf("invalid password: %s\n", converted.PasswordInfo.Password)
	}

	if converted.PasswordInfo.Salt != "salthere" {
		t.Errorf("invalid salt: %s\n", converted.PasswordInfo.Salt)
	}

	if len(converted.PasswordSessions) != 0 {
		t.Errorf("invalid password session length, expected 2, got %d", len(converted.PasswordSessions))
	}

	// confirm custom field is set
	if converted.customObject != "helloworld" {
		t.Errorf("Custom object was not set!")
	}
}
