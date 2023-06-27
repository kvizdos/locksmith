package validation

import (
	"fmt"
	"net/http"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

func ValidateHTTPUserToken(r *http.Request, db database.DatabaseAccessor) (users.LocksmithUser, error) {
	// Validate token
	token, err := r.Cookie("token")

	if err != nil {
		return users.LocksmithUser{}, fmt.Errorf("no cookie present")
	}

	parsedToken, err := authentication.ParseToken(token.Value)

	if err != nil {
		return users.LocksmithUser{}, fmt.Errorf("token could not be parsed")
	}

	user, validated, err := ValidateToken(parsedToken, db)

	if err != nil {
		return users.LocksmithUser{}, fmt.Errorf("token could not be validated")
	}

	if !validated {
		return users.LocksmithUser{}, fmt.Errorf("token did not validate")
	}

	return user.(users.LocksmithUser), nil
}
