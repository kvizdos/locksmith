package validation

import (
	"fmt"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

func ValidateToken(token authentication.Token, db database.DatabaseAccessor) (users.LocksmithUserStruct, bool, error) {
	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": token.Username,
	})

	if !usernameExists {
		return users.LocksmithUser{}, false, fmt.Errorf("invalid username")
	}

	user := users.LocksmithUserFromMap(dbUser.(map[string]interface{}))

	validated := user.ValidateSessionToken(token.Token, db)

	if !validated {
		return users.LocksmithUser{}, false, nil
	}

	return user, true, nil
}
