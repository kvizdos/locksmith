package magic

import (
	"github.com/kvizdos/locksmith/database"
)

func ExpireOld(db database.DatabaseAccessor, lookupUserID string, manuallyExpireTokenID ...string) {
	rawUser, found := db.FindOne("users", map[string]any{
		"id": lookupUserID,
	})

	if !found {
		return
	}

	user := rawUser.(map[string]any)

	magics := MagicsFromMap(user["magic"].([]any))

	active := make(chan MagicAuthentications)
	go FilterActive(active, magics, manuallyExpireTokenID...)
	keep := <-active

	if len(keep) != len(magics) {
		db.UpdateOne("users", map[string]any{
			"id": lookupUserID,
		}, map[database.DatabaseUpdateActions]map[string]any{
			database.SET: {
				"magic": keep.ToMap(),
			},
		})
	}
}
