package service_keys

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateClientID(name string) (string, error) {
	randomPart, err := authentication.GenerateRandomStringURLSafe(8)
	if err != nil {
		return "", err
	}

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	return fmt.Sprintf("svc-%s-%s", name, randomPart), nil
}

type CreateServiceKeyOptions[T any] struct {
	FriendlyName    string
	Scopes          []string
	AllowListIPs    []string
	CreatedByUserID string
	Claims          T
}

func CreateServiceKey[T any](db database.DatabaseAccessor, opts CreateServiceKeyOptions[T]) (*ServiceKeyCredentials, error) {
	_, found := db.FindOne("service_keys", map[string]any{
		"name": opts.FriendlyName,
	})

	if found {
		return nil, fmt.Errorf("service key name already used")
	}

	seen := map[string]bool{}

	for _, scope := range opts.Scopes {
		if seen[scope] {
			return nil, fmt.Errorf("duplicate scope: %s", scope)
		}
		seen[scope] = true

		if _, ok := roles.AVAILABLE_PERMISSIONS[scope]; !ok {
			return nil, fmt.Errorf("scope not available: %s", scope)
		}
	}

	clientID, err := generateClientID(opts.FriendlyName)

	if err != nil {
		return nil, err
	}

	secretKey, err := authentication.GenerateRandomStringURLSafe(256)

	if err != nil {
		return nil, err
	}

	hash, err := argon2id.CreateHash(secretKey, argon2id.DefaultParams)

	if err != nil {
		return nil, err
	}

	secretID := uuid.NewString()

	now := time.Now().UTC()
	key := ServiceKey{
		ID:           primitive.NewObjectID(), // used for db only
		ClientID:     clientID,
		Name:         opts.FriendlyName,
		Scopes:       opts.Scopes,
		CreatedAt:    now,
		CreatedBy:    opts.CreatedByUserID,
		AllowListIPs: opts.AllowListIPs,
		Secrets: []ServiceSecret{
			{
				ID:         secretID,
				SecretHash: hash,
				CreatedAt:  now,
			},
		},
	}

	_, err = db.InsertOne("service_keys", key.ToMap())

	if err != nil {
		return nil, err
	}

	return &ServiceKeyCredentials{
		ClientID: clientID,
		SecretID: secretID,
		Secret:   secretKey,
	}, nil
}
