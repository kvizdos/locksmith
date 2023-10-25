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

func Validate(db database.DatabaseAccessor, tokenIdentifer string) (MagicAuthentication, map[string]interface{}, error) {
	identifier, err := base64.StdEncoding.DecodeString(tokenIdentifer)

	if err != nil {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("bad identifier")
	}

	identifierParts := strings.Split(string(identifier), ":")

	if len(identifierParts) != 4 {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("incorrect number of parts")
	}

	userID := identifierParts[0]
	token := identifierParts[1]
	expiresAtStr := identifierParts[2]
	signature := identifierParts[3]

	if len(token) != 128 {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("incorrect token length")
	}

	if len(userID) != 36 {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("incorrect uid length")
	}

	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("invalid expiresAt")
	}

	if expiresAt <= time.Now().UTC().Unix() {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("expired")
	}

	if !MagicSigningPackage.Validate(fmt.Sprintf("%s:%s:%s", userID, token, expiresAtStr), signature) {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("invalid signature")
	}

	rawUser, found := db.FindOne("users", map[string]interface{}{
		"id": userID,
	})

	if !found {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("uid not found")
	}

	if rawUser.(map[string]interface{})["magic"] == nil {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("no magics found")
	}

	magics := MagicsFromMap(rawUser.(map[string]interface{})["magic"].([]interface{}))

	if len(magics) == 0 {
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("no magics found")
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
		return MagicAuthentication{}, map[string]interface{}{}, fmt.Errorf("invalidated")
	}

	foundMagic.InheritRole = rawUser.(map[string]interface{})["role"].(string)
	foundMagic.Username = rawUser.(map[string]interface{})["username"].(string)

	go ExpireOld(db, userID)

	return foundMagic, rawUser.(map[string]interface{}), nil
}
