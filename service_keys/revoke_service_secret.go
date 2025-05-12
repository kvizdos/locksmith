package service_keys

import (
	"fmt"
	"time"

	"github.com/kvizdos/locksmith/database"
)

func RevokeServiceSecret(db database.DatabaseAccessor, clientID, secretID, revokedBy string) error {
	// Lookup existing ServiceKey
	rawKey, found := db.FindOne("service_keys", map[string]any{
		"client_id": clientID,
	})

	if !found {
		return fmt.Errorf("service key not found")
	}

	key := ServiceKeyFromMap(rawKey)

	// Find and revoke the secret
	now := time.Now().UTC()
	foundSecret := false

	for i, secret := range key.Secrets {
		if secret.ID == secretID {
			if secret.RevokedAt != nil {
				return fmt.Errorf("secret already revoked")
			}
			key.Secrets[i].RevokedAt = &now
			key.Secrets[i].RevokedBy = revokedBy
			foundSecret = true
			break
		}
	}

	if !foundSecret {
		return fmt.Errorf("secret not found")
	}

	// Update database
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
