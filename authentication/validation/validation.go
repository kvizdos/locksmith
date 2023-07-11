package validation

import (
	"fmt"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func ValidateToken(token authentication.Token, db database.DatabaseAccessor) (users.LocksmithUserInterface, bool, error) {
	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": token.Username,
	})

	if !usernameExists {
		return users.LocksmithUser{}, false, fmt.Errorf("invalid username")
	}

	var tmpUser users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

	validated := user.ValidateSessionToken(token.Token, db)

	if !validated {
		return users.LocksmithUser{}, false, nil
	}

	return user, true, nil
}
