package magic

import "github.com/kvizdos/locksmith/database"

func ExpireOld(db database.DatabaseAccessor, lookupUserID string, manuallyExpireTokenID ...string) {
	rawUser, found := db.FindOne("users", map[string]interface{}{
		"id": lookupUserID,
	})

	if !found {
		return
	}

	user := rawUser.(map[string]interface{})

	magics := MagicsFromMap(user["magic"].([]map[string]interface{}))

	active := make(chan MagicAuthentications)
	go FilterActive(active, magics, manuallyExpireTokenID...)
	keep := <-active

	if len(keep) != len(magics) {
		db.UpdateOne("users", map[string]interface{}{
			"id": lookupUserID,
		}, map[database.DatabaseUpdateActions]map[string]interface{}{
			database.SET: {
				"magic": keep.ToMap(),
			},
		})
	}
}
