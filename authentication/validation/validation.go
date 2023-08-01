package validation

import (
	"fmt"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func ValidateToken(token authentication.Token, db database.DatabaseAccessor, userType ...users.LocksmithUserInterface) (users.LocksmithUserInterface, bool, error) {
	var lsu users.LocksmithUserInterface

	if len(userType) > 0 {
		lsu = userType[0]
	} else {
		lsu = users.LocksmithUser{}
	}

	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": token.Username,
	})

	if !usernameExists {
		return lsu, false, fmt.Errorf("invalid username")
	}

	var tmpUser users.LocksmithUserInterface
	lsu.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser

	validated := user.ValidateSessionToken(token.Token, db)

	if !validated {
		return lsu, false, nil
	}

	return user, true, nil
}
