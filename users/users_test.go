package users

import (
	"testing"

	"github.com/kvizdos/locksmith/authentication"
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
		"email":    "email@email.com",
		"password": map[string]interface{}{
			"password": "passwordhere",
			"salt":     "salthere",
		},
		"sessions": sessions,
		"role":     "user",
	}

	var user LocksmithUserInterface
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
	LocksmithUserInterface
}

type customUser struct {
	LocksmithUser

	customObject string
}

func (c customUser) ReadFromMap(writeTo *LocksmithUserInterface, u map[string]interface{}) {
	// Load initial locksmith data
	var user LocksmithUserInterface
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
		"email":    "email@email.com",
		"password": map[string]interface{}{
			"password": "passwordhere",
			"salt":     "salthere",
		},
		"sessions":     []interface{}{},
		"role":         "user",
		"customObject": "helloworld",
	}

	var user LocksmithUserInterface
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

func TestConvertUserToPublicUser(t *testing.T) {
	privateUser := LocksmithUser{
		ID:       "userIDhere",
		Username: "kvizdos",
		Email:    "email@email.com",
		PasswordInfo: authentication.PasswordInfo{
			Password: "password",
			Salt:     "salt",
		},
		PasswordSessions: []authentication.PasswordSession{
			{
				Token:     "token here",
				ExpiresAt: 0,
			},
			{
				Token:     "another token here",
				ExpiresAt: 0,
			},
		},
		Role: "user",
	}

	convertedUser, err := PublicLocksmithUser{}.FromRegular(privateUser)

	if err != nil {
		t.Errorf("received error while converting from regular: %s", err.Error())
		return
	}

	convertedUser, success := convertedUser.(PublicLocksmithUser)

	if !success {
		t.Errorf("could not convert to PublicLocksmithUser")
		return
	}
}
