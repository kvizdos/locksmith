package validation

import (
	"crypto/ecdsa"
	"fmt"
	"reflect"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type ValidationSigningKeys struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

type ValidationOptions struct {
	SigningKeys ValidationSigningKeys
	UserType    users.LocksmithUserInterface
	Claims      jwt.Claims
}

func getIDFromClaims(claims interface{}) (string, error) {
	val := reflect.ValueOf(claims)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("claims is not a struct")
	}

	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("ID field not found in claims")
	}

	if idField.Kind() != reflect.String {
		return "", fmt.Errorf("ID field is not a string")
	}

	return idField.String(), nil
}

func ValidateToken(token authentication.Token, db database.DatabaseAccessor, magicToken string, opts ValidationOptions) (users.LocksmithUserInterface, bool, error) {
	var lsu users.LocksmithUserInterface
	if opts.UserType != nil {
		lsu = opts.UserType
	} else {
		lsu = users.LocksmithUser{}
	}

	var user users.LocksmithUserInterface

	if magicToken == "" {
		validated, claims := lsu.ValidateAccessJWT(token.Token, opts.SigningKeys.PublicKey, opts.Claims)

		if !validated {
			return lsu, false, nil
		}

		jwtClaimID, err := getIDFromClaims(claims)

		if err != nil {
			return lsu, false, nil
		}
		validated = lsu.ValidateProfileJWT(token.ProfileToken, jwtClaimID, opts.SigningKeys.PublicKey)

		if !validated {
			return lsu, false, nil
		}

		user = lsu.FromAccessJWTClaims(claims)
	} else {
		panic("magic tokens not supported!")
		// user = user.SetMagicPermissions(mac.AllowedPermissions)
		// if token.Token == "" {
		// 	user = user.SetMagic()
		// }
	}

	// go user.CleanupOldMagicTokens(db)

	return user, true, nil
}
