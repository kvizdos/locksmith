package service_keys

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
)

func AddServiceSecret(db database.DatabaseAccessor, clientID string, createdByUserID string) (*ServiceKeyCredentials, error) {
	// Lookup existing service key
	rawKey, found := db.FindOne("service_keys", map[string]any{
		"client_id": clientID,
	})

	if !found {
		return nil, fmt.Errorf("service key not found")
	}

	// Deserialize
	key := ServiceKeyFromMap(rawKey)

	// Create new secret
	secretKey, err := authentication.GenerateRandomStringURLSafe(256)
	if err != nil {
		return nil, err
	}

	hash, err := argon2id.CreateHash(secretKey, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	newSecret := ServiceSecret{
		ID:         uuid.NewString(),
		SecretHash: hash,
		CreatedAt:  now,
		CreatedBy:  createdByUserID,
	}

	// Append
	key.Secrets = append(key.Secrets, newSecret)

	// Update DB
	_, err = db.UpdateOne("service_keys", map[string]any{
		"client_id": clientID,
	}, map[database.DatabaseUpdateActions]map[string]any{
		database.SET: key.ToMap(),
	})

	if err != nil {
		return nil, err
	}

	return &ServiceKeyCredentials{
		ClientID: clientID,
		SecretID: newSecret.ID,
		Secret:   secretKey,
	}, nil
}
