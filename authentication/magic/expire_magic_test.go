package magic

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/database"
)

func TestExpireOldAutomatically(t *testing.T) {
	userID := uuid.New().String()

	mac1, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          userID,
		AllowedPermissions: []string{},
		TTL:                -1 * time.Hour,
	})

	mac2, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          userID,
		AllowedPermissions: []string{"abc"},
		TTL:                time.Hour,
	})
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id": userID,
					"magic": MagicAuthentications{
						mac1,
						mac2,
					}.ToMap(),
				},
			},
		},
	}

	ExpireOld(testDb, userID)

	rawUser, _ := testDb.FindOne("users", map[string]interface{}{
		"id": userID,
	})

	user := rawUser.(map[string]interface{})

	magics := MagicsFromMap(user["magic"].([]interface{}))
	if len(magics) != 1 {
		t.Errorf("got incorrect number of magics: %d", len(magics))
		return
	}

	if len(magics[0].AllowedPermissions) != 1 {
		t.Errorf("got incorrect number of magic permissions: %d", len(magics[0].AllowedPermissions))
		return
	}

	if magics[0].AllowedPermissions[0] != "abc" {
		t.Errorf("got bad permission: %s", magics[0].AllowedPermissions[0])
	}
}

func TestExpireOldWithManual(t *testing.T) {
	userID := uuid.New().String()

	mac1, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          userID,
		AllowedPermissions: []string{},
		TTL:                -1 * time.Hour,
	})

	mac2, manualID, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          userID,
		AllowedPermissions: []string{"abc"},
		TTL:                time.Hour,
	})
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id": userID,
					"magic": MagicAuthentications{
						mac1,
						mac2,
					}.ToMap(),
				},
			},
		},
	}

	ExpireOld(testDb, userID, manualID)

	rawUser, _ := testDb.FindOne("users", map[string]interface{}{
		"id": userID,
	})
	user := rawUser.(map[string]interface{})

	magics := MagicsFromMap(user["magic"].([]interface{}))

	if len(magics) != 0 {
		t.Errorf("got incorrect number of magics: %d", len(magics))
		return
	}
}
