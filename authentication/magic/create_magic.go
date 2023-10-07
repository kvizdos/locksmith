package magic

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/kvizdos/locksmith/authentication"
)

// Returns
// The MagicAuth for saving in the Database
// The Identifier String (userID:code:signature)
// Error, if any.
func CreateMagicAuthentication(macVariables MagicAuthenticationVariables) (MagicAuthentication, string, error) {
	authToken, err := authentication.GenerateRandomString(128)

	if err != nil {
		return MagicAuthentication{}, "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(authToken))
	hashedToken := fmt.Sprintf("%x", hasher.Sum(nil))

	expiresAt := time.Now().UTC().Add(macVariables.TTL)

	mac := MagicAuthentication{
		Code:               hashedToken,
		AllowedPermissions: macVariables.AllowedPermissions,
		ExpiresAt:          expiresAt.Unix(),
	}

	tokenIdentifier := fmt.Sprintf("%s:%s:%d", macVariables.ForUserID, authToken, expiresAt.Unix())

	signature, err := MagicSigningPackage.Sign(tokenIdentifier)

	if err != nil {
		return MagicAuthentication{}, "", err
	}

	signedToken := fmt.Sprintf("%s:%s", tokenIdentifier, signature)

	base64d := base64.StdEncoding.EncodeToString([]byte(signedToken))

	return mac, base64d, nil
}
