package validation

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type MagicValidation struct {
	Token      string
	Prioritize bool
}

type HTTPValidationOptions struct {
	UserType       users.LocksmithUserInterface
	ValidationKeys ValidationSigningKeys
	Claims         jwt.Claims
}

func ValidateHTTPUserToken(r *http.Request, db database.DatabaseAccessor, magic MagicValidation, opts HTTPValidationOptions) (users.LocksmithUserInterface, error) {
	var userInterface users.LocksmithUserInterface

	if opts.UserType != nil {
		userInterface = opts.UserType
	} else {
		userInterface = users.LocksmithUser{}
	}

	// Validate token
	token, tokenErr := r.Cookie("token")
	profileToken, profileErr := r.Cookie("profile")

	if tokenErr != nil || profileErr != nil {
		return userInterface, fmt.Errorf("missing cookies")
	}

	var claims jwt.Claims

	if opts.Claims != nil {
		claims = opts.Claims
	} else {
		claims = &users.BaseValidationClaims{}
	}

	// if magic.Token == "" || (err == nil && !magic.Prioritize) {
	// 	if err != nil {
	// 		return userInterface, fmt.Errorf("no cookie present")
	// 	}

	// 	parsedToken, err = authentication.ParseToken(token.Value)

	// 	if err != nil {
	// 		return userInterface, fmt.Errorf("token could not be parsed")
	// 	}
	// }

	user, validated, err := ValidateToken(authentication.Token{
		Token:        token.Value,
		ProfileToken: profileToken.Value,
		Username:     "",
	}, db, magic.Token, ValidationOptions{
		SigningKeys: opts.ValidationKeys,
		UserType:    userInterface,
		Claims:      claims,
	})

	if err != nil {
		return userInterface, fmt.Errorf("token could not be validated")
	}

	if !validated {
		return userInterface, fmt.Errorf("token did not validate")
	}

	return user, nil
}
