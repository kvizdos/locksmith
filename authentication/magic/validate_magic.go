package magic

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kvizdos/locksmith/database"
)

func Validate(db database.DatabaseAccessor, tokenIdentifer string) (MagicAuthentication, error) {
	identifier, err := base64.StdEncoding.DecodeString(tokenIdentifer)

	if err != nil {
		return MagicAuthentication{}, fmt.Errorf("bad identifier")
	}

	identifierParts := strings.Split(string(identifier), ":")

	if len(identifierParts) != 4 {
		return MagicAuthentication{}, fmt.Errorf("incorrect number of parts")
	}

	userID := identifierParts[0]
	token := identifierParts[1]
	expiresAtStr := identifierParts[2]
	signature := identifierParts[3]

	if len(token) != 128 {
		return MagicAuthentication{}, fmt.Errorf("incorrect token length")
	}

	if len(userID) != 36 {
		return MagicAuthentication{}, fmt.Errorf("incorrect uid length")
	}

	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return MagicAuthentication{}, fmt.Errorf("invalid expiresAt")
	}

	if expiresAt <= time.Now().UTC().Unix() {
		return MagicAuthentication{}, fmt.Errorf("expired")
	}

	if !MagicSigningPackage.Validate(fmt.Sprintf("%s:%s:%s", userID, token, expiresAtStr), signature) {
		return MagicAuthentication{}, fmt.Errorf("invalid signature")
	}

	rawUser, found := db.FindOne("users", map[string]interface{}{
		"id": userID,
	})

	if !found {
		return MagicAuthentication{}, fmt.Errorf("uid not found")
	}

	magics := MagicsFromMap(rawUser.(map[string]interface{})["magic"].([]map[string]interface{}))

	if len(magics) == 0 {
		return MagicAuthentication{}, fmt.Errorf("no magics found")
	}

	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedToken := fmt.Sprintf("%x", hasher.Sum(nil))

	var foundMagic MagicAuthentication
	var didFindMagic bool
	for _, m := range magics {
		if m.Code == hashedToken {
			foundMagic = m
			didFindMagic = true
			break
		}
	}

	if !didFindMagic {
		return MagicAuthentication{}, fmt.Errorf("invalidated")
	}

	foundMagic.InheritRole = rawUser.(map[string]interface{})["role"].(string)

	go ExpireOld(db, userID)

	return foundMagic, nil
}
