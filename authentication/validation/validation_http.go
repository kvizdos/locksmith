package validation

import (
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type MagicValidation struct {
	Token      string
	Prioritize bool
}

func ValidateHTTPUserToken(r *http.Request, db database.DatabaseAccessor, magic MagicValidation, userType ...users.LocksmithUserInterface) (users.LocksmithUserInterface, error) {
	// Validate token
	token, err := r.Cookie("token")

	var userInterface users.LocksmithUserInterface

	if len(userType) > 0 {
		userInterface = userType[0]
	} else {
		userInterface = users.LocksmithUser{}
	}

	var parsedToken authentication.Token
	if magic.Token == "" || (err == nil && !magic.Prioritize) {
		if err != nil {
			return userInterface, fmt.Errorf("no cookie present")
		}

		parsedToken, err = authentication.ParseToken(token.Value)

		if err != nil {
			return userInterface, fmt.Errorf("token could not be parsed")
		}
	}

	user, validated, err := ValidateToken(parsedToken, db, magic.Token, userInterface)

	if err != nil {
		return userInterface, fmt.Errorf("token could not be validated")
	}

	if !validated {
		return userInterface, fmt.Errorf("token did not validate")
	}

	return user, nil
}
