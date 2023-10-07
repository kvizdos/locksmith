package magic

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
)

func TestValidateMagicInvalidBase64(t *testing.T) {
	testDb := database.TestDatabase{}

	macIdentifer := "dGVzdC11aWQ6Y-----!-[]TlpQ2lKWjhGelZjZHVIU1c4ZTlaUE1LUmNhNVlRd1h6Ulkydy1zVTBCR3d4MkFNeURma3kxT3V6YnkyRHJLZFJpeVJ5d2dlRTJFYVBxOHBBQ1k3ZThnRkdxVmVvQi1KcWFuS3hoM0dKMFhZd2QyUjZUOVMyMVRuR3FNQXJDU0g6UFdTU1RqTzhSNXNraGVJdUs1NUJpYTVWMUNjZ240MHBZdndVZHNmWk5zdkMxOVhBTWxYcitodVk4TFpMTWR0MkxIWDJ1VU5ITG90N2RpNE5zeGlYVkE9Pz=="

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "bad identifier" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicIncorrectNumberOfParts(t *testing.T) {
	testDb := database.TestDatabase{}

	macIdentifer, _ := MagicSigningPackage.Sign("a:b:c")

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "incorrect number of parts" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicInvalidTokenLength(t *testing.T) {
	testDb := database.TestDatabase{}

	userID := uuid.New().String()
	tok := "token-here"
	expires := time.Now().UTC().Add(time.Hour).Unix()
	info := fmt.Sprintf("%s:%s:%d", userID, tok, expires)
	sig, _ := MagicSigningPackage.Sign(info)
	rawMac := fmt.Sprintf("%s:%s", info, sig)
	macIdentifer := base64.StdEncoding.EncodeToString([]byte(rawMac))

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "incorrect token length" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicInvalidUserIDLength(t *testing.T) {
	testDb := database.TestDatabase{}

	authToken, _ := authentication.GenerateRandomString(128)

	userID := "bad-length"
	tok := authToken
	expires := "123456"
	info := fmt.Sprintf("%s:%s:%s", userID, tok, expires)
	sig, _ := MagicSigningPackage.Sign(info)
	rawMac := fmt.Sprintf("%s:%s", info, sig)
	macIdentifer := base64.StdEncoding.EncodeToString([]byte(rawMac))

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "incorrect uid length" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicAlreadyExpired(t *testing.T) {
	testDb := database.TestDatabase{}

	_, macIdentifer, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uuid.NewString(),
		AllowedPermissions: []string{"xyz", "abc"},
		TTL:                -1 * time.Hour,
	})
	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "expired" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicInvalidSignature(t *testing.T) {
	testDb := database.TestDatabase{}

	authToken, _ := authentication.GenerateRandomString(128)

	userID := uuid.New().String()
	tok := authToken
	expires := time.Now().UTC().Add(time.Hour).Unix()
	info := fmt.Sprintf("%s:%s:%d", userID, tok, expires)
	sig, _ := MagicSigningPackage.Sign(info + "z")
	rawMac := fmt.Sprintf("%s:%s", info, sig)
	macIdentifer := base64.StdEncoding.EncodeToString([]byte(rawMac))

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "invalid signature" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicValidButUserWasDeleted(t *testing.T) {
	uid := uuid.NewString()
	_, macIdentifer, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"xyz", "abc"},
		TTL:                1 * time.Hour,
	})

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "rand",
					"username": "kenton",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role":  "user",
					"magic": MagicAuthentications{}.ToMap(),
				},
			},
		},
	}

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "uid not found" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicValidButUserHas0ActiveMagics(t *testing.T) {
	uid := uuid.NewString()
	_, macIdentifer, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"xyz", "abc"},
		TTL:                1 * time.Hour,
	})

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       uid,
					"username": "kenton",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role":  "user",
					"magic": MagicAuthentications{}.ToMap(),
				},
			},
		},
	}

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "no magics found" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicValidButDoesNotExistInDatabase(t *testing.T) {
	uid := uuid.NewString()
	randomMac, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"bac", "abc"},
		TTL:                1 * time.Hour,
	})
	_, macIdentifer, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"xyz", "abc"},
		TTL:                1 * time.Hour,
	})

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       uid,
					"username": "kenton",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role": "user",
					"magic": MagicAuthentications{
						randomMac,
					}.ToMap(),
				},
			},
		},
	}

	_, err := Validate(testDb, macIdentifer)

	if err == nil {
		t.Error("expected to find an error!")
		return
	}

	if err.Error() != "invalidated" {
		t.Errorf("got a weird error: %s", err)
	}
}

func TestValidateMagicValidSuccess(t *testing.T) {
	uid := uuid.NewString()
	randomMac, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"bac"},
		TTL:                1 * time.Hour,
	})
	correctMac, macIdentifer, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"xyz", "abc"},
		TTL:                1 * time.Hour,
	})

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       uid,
					"username": "kenton",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role": "user",
					"magic": MagicAuthentications{
						randomMac,
						correctMac,
					}.ToMap(),
				},
			},
		},
	}

	foundMac, err := Validate(testDb, macIdentifer)

	if err != nil {
		t.Errorf("got a weird error: %s", err)
		return
	}

	if foundMac.InheritRole != "user" {
		t.Errorf("got incorrect role: %s", foundMac.InheritRole)
	}

	if len(foundMac.AllowedPermissions) != 2 {
		t.Errorf("got incorrect number of permissions: %d", len(foundMac.AllowedPermissions))
	}
}

func TestValidateMagicOldCodesAreExpired(t *testing.T) {
	uid := uuid.NewString()
	randomMac, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"bac"},
		TTL:                -1 * time.Hour,
	})
	correctMac, macIdentifer, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          uid,
		AllowedPermissions: []string{"xyz", "abc"},
		TTL:                1 * time.Hour,
	})

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       uid,
					"username": "kenton",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role": "user",
					"magic": MagicAuthentications{
						randomMac,
						correctMac,
					}.ToMap(),
				},
			},
		},
	}

	_, err := Validate(testDb, macIdentifer)

	if err != nil {
		t.Errorf("got a weird error: %s", err)
		return
	}

	// Wait to make sure the ExpireOld() go routine
	// completes. (i don't want to add a channel / waitgroup
	// because it'd never be necessary to know in other parts
	// of code)
	time.Sleep(time.Second * 2)

	rawUser, _ := testDb.FindOne("users", map[string]interface{}{
		"id": uid,
	})
	user := rawUser.(map[string]interface{})

	magics := MagicsFromMap(user["magic"].([]map[string]interface{}))

	if len(magics) != 1 {
		t.Errorf("got incorrect number of magics: %d", len(magics))
		return
	}
}
