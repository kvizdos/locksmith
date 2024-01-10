package validation

import (
	"fmt"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func ValidateToken(token authentication.Token, db database.DatabaseAccessor, magicToken string, userType ...users.LocksmithUserInterface) (users.LocksmithUserInterface, bool, error) {
	var lsu users.LocksmithUserInterface
	if len(userType) > 0 {
		lsu = userType[0]
	} else {
		lsu = users.LocksmithUser{}
	}

	var dbUser map[string]interface{}
	var mac magic.MagicAuthentication

	if magicToken != "" && token.Token == "" {
		rm, userInfo, err := magic.Validate(db, magicToken)

		if err != nil {
			return lsu, false, fmt.Errorf("bad magic token")
		}
		mac = rm
		token.Username = mac.Username
		dbUser = userInfo
	} else {
		rawUser, usernameExists := db.FindOne("users", map[string]interface{}{
			"id": token.Username,
		})

		if !usernameExists {
			return lsu, false, fmt.Errorf("invalid username")
		}

		dbUser = rawUser.(map[string]interface{})
	}

	var tmpUser users.LocksmithUserInterface
	lsu.ReadFromMap(&tmpUser, dbUser)
	user := tmpUser

	if magicToken == "" {
		validated := user.ValidateSessionToken(token.Token, db)

		if !validated {
			return lsu, false, nil
		}
	} else {
		user = user.SetMagicPermissions(mac.AllowedPermissions)
		if token.Token == "" {
			user = user.SetMagic()
		}
	}

	go user.CleanupOldMagicTokens(db)

	return user, true, nil
}
