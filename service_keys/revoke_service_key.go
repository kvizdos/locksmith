package service_keys

import (
	"fmt"
	"time"

	"github.com/kvizdos/locksmith/database"
)

func RevokeServiceKey(db database.DatabaseAccessor, clientID, revokedBy string) error {
	rawKey, found := db.FindOne("service_keys", map[string]any{
		"client_id": clientID,
	})

	if !found {
		return fmt.Errorf("service key not found")
	}

	key := ServiceKeyFromMap(rawKey)

	if key.RevokedAt != nil {
		return fmt.Errorf("service key already revoked")
	}

	now := time.Now().UTC()
	key.RevokedAt = &now
	key.RevokedBy = revokedBy

	_, err := db.UpdateOne("service_keys", map[string]any{
		"client_id": clientID,
	}, map[database.DatabaseUpdateActions]map[string]any{
		database.SET: key.ToMap(),
	})

	if err != nil {
		return err
	}

	return nil
}
