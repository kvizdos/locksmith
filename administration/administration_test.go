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
					"email":    "email@email.com",
					"password": testPassword,
					"sessions": []interface{}{},
					"role":     "user",
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
					"email":    "email@email.com",
					"password": testPassword,
					"sessions": []interface{}{},
					"role":     "user",
				},
				"a2bHHs4L-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "a2bHHs4L-22a7-493f-b228-028842e09a05",
					"username": "bob",
					"email":    "email@email.com",
					"password": testPassword,
					"sessions": []interface{}{},
					"role":     "user",
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

type customUser struct {
	users.LocksmithUser

	CustomObject string `json:"customObject"`
}

func (u customUser) ToPublic() (users.PublicLocksmithUserInterface, error) {
	publicUser, err := publicCustomUser{}.FromRegular(u)

	return publicUser, err
}

type publicCustomUser struct {
	users.PublicLocksmithUser

	CustomObject string `json:"customObject"`
}

func (u publicCustomUser) FromRegular(user users.LocksmithUserInterface) (users.PublicLocksmithUserInterface, error) {
	lsPub, err := u.PublicLocksmithUser.FromRegular(user)

	if err != nil {
		return publicCustomUser{}, nil
	}

	publicUser := publicCustomUser{
		PublicLocksmithUser: lsPub.(users.PublicLocksmithUser),
	}

	publicUser.CustomObject = user.(customUser).CustomObject

	return publicUser, nil
}

func (c customUser) ReadFromMap(writeTo *users.LocksmithUserInterface, u map[string]interface{}) {
	// Load initial locksmith data
	var user users.LocksmithUserInterface
	c.LocksmithUser.ReadFromMap(&user, u)
	lsu := user.(users.LocksmithUser)

	converted := customUser{
		LocksmithUser: lsu,
	}

	converted.CustomObject = u["customObject"].(string)

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
					"email":        "email@email.com",
					"password":     testPassword,
					"sessions":     []interface{}{},
					"role":         "user",
					"customObject": "helloworld",
				},
			},
		},
	}

	publicUser := customUser{}

	usersArr, err := ListUsers(testDb, publicUser)

	if err != nil {
		t.Errorf("unexpected listing error: %s", err)
		return
	}

	expecting := 1
	if len(usersArr) != expecting {
		t.Errorf("did not receive correct number of users, expected %d, got %d", expecting, len(usersArr))
		return
	}

	value, ok := usersArr[0].(publicCustomUser)

	if !ok {
		t.Errorf("failed to convert to custom user object")
		return
	}

	if value.CustomObject != "helloworld" {
		t.Errorf("customObject incorrect: %s\n", value.CustomObject)
	}
}

func TestDeleteUserNoUser(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	deleted, err := DeleteUser(testDb, "kvizdos")

	if err != nil {
		t.Errorf("unexpected listing error: %s", err)
		return
	}

	if deleted {
		t.Errorf("user should not have been deleted as none exist")
	}
}

func TestDeleteUser(t *testing.T) {
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

	deleted, err := DeleteUser(testDb, "kenton")

	if err != nil {
		t.Errorf("unexpected listing error: %s", err)
		return
	}

	if !deleted {
		t.Errorf("user should have been deleted")
		return
	}

	if _, exists := testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"]; exists {
		t.Errorf("user not deleted")
	}
}
