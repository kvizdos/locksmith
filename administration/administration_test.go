package administration

import (
	"testing"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

func TestListUsersNoUsers(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	users, err := ListUsers(testDb)

	if err != nil {
		t.Errorf("unexpected listing error: %s", err)
		return
	}

	expecting := 0
	if len(users) != expecting {
		t.Errorf("did not receive correct number of users, expected %d, got %d", expecting, len(users))
		return
	}
}

func TestListUsersOneUser(t *testing.T) {
	testPassword, _ := authentication.CompileLocksmithPassword("securepassword123")
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"password": testPassword,
					"sessions": []interface{}{},
				},
			},
		},
	}

	usersArr, err := ListUsers(testDb)

	if err != nil {
		t.Errorf("unexpected listing error: %s", err)
		return
	}

	expecting := 1
	if len(usersArr) != expecting {
		t.Errorf("did not receive correct number of users, expected %d, got %d", expecting, len(usersArr))
		return
	}
}

func TestListUsersMultipleUsers(t *testing.T) {
	testPassword, _ := authentication.CompileLocksmithPassword("securepassword123")
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"password": testPassword,
					"sessions": []interface{}{},
				},
				"a2bHHs4L-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "a2bHHs4L-22a7-493f-b228-028842e09a05",
					"username": "bob",
					"password": testPassword,
					"sessions": []interface{}{},
				},
			},
		},
	}

	usersArr, err := ListUsers(testDb)

	if err != nil {
		t.Errorf("unexpected listing error: %s", err)
		return
	}

	expecting := 2
	if len(usersArr) != expecting {
		t.Errorf("did not receive correct number of users, expected %d, got %d", expecting, len(usersArr))
		return
	}
}

type customUserInterface interface {
	users.LocksmithUserStruct
}

type customUser struct {
	users.LocksmithUser

	customObject string
}

func (c customUser) ReadFromMap(writeTo *users.LocksmithUserStruct, u map[string]interface{}) {
	// Load initial locksmith data
	var user users.LocksmithUserStruct
	c.LocksmithUser.ReadFromMap(&user, u)
	lsu := user.(users.LocksmithUser)

	converted := customUser{
		LocksmithUser: lsu,
	}

	converted.customObject = u["customObject"].(string)

	*writeTo = converted
}

func TestListUsersOneUserCustomStruct(t *testing.T) {
	testPassword, _ := authentication.CompileLocksmithPassword("securepassword123")
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":           "c8531661-22a7-493f-b228-028842e09a05",
					"username":     "kenton",
					"password":     testPassword,
					"sessions":     []interface{}{},
					"customObject": "helloworld",
				},
			},
		},
	}

	usersArr, err := ListUsers(testDb, customUser{})

	if err != nil {
		t.Errorf("unexpected listing error: %s", err)
		return
	}

	expecting := 1
	if len(usersArr) != expecting {
		t.Errorf("did not receive correct number of users, expected %d, got %d", expecting, len(usersArr))
		return
	}

	value, ok := usersArr[0].(customUser)

	if !ok {
		t.Errorf("failed to convert to custom user object")
		return
	}

	if value.customObject != "helloworld" {
		t.Errorf("customObject incorrect: %s\n", value.customObject)
	}
}
