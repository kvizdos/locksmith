package service_keys_test

import (
	"testing"
	"time"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/service_keys"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddServiceSecret_Success(t *testing.T) {
	now := time.Now().UTC()

	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"rand-id": service_keys.ServiceKey{
					ID:       primitive.NewObjectID(),
					ClientID: "test-client-id",
					Name:     "Test Name",
					Scopes:   []string{"test"},
					Secrets: []service_keys.ServiceSecret{
						{
							ID:         "test-id",
							SecretHash: "existinghash",
							CreatedAt:  now,
							CreatedBy:  "test-user",
							LastUsedAt: nil,
							RevokedAt:  nil,
						},
					},
					CreatedAt: now,
					CreatedBy: "test-user",
				}.ToMap(),
			},
		},
	}

	creds, err := service_keys.AddServiceSecret(db, "test-client-id", "test-user-2")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if creds.ClientID != "test-client-id" {
		t.Errorf("expected clientID to match, got %s", creds.ClientID)
	}

	if creds.Secret == "" {
		t.Error("expected secret to be returned")
	}

	if creds.SecretID == "" {
		t.Error("expected secret ID to be returned")
	}

	// Check DB update
	rawRes, found := db.FindOne("service_keys", map[string]any{
		"client_id": "test-client-id",
	})
	if !found {
		t.Fatal("expected service key to still exist after insert")
	}

	key := service_keys.ServiceKeyFromMap(rawRes)

	if len(key.Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(key.Secrets))
	}

	var foundNewSecret bool
	for _, secret := range key.Secrets {
		if secret.ID == creds.SecretID {
			foundNewSecret = true

			if secret.RevokedAt != nil {
				t.Error("new secret should not be revoked")
			}

			if secret.CreatedBy != "test-user-2" {
				t.Errorf("expected CreatedBy to be 'test-user-2', got %s", secret.CreatedBy)
			}

			if secret.CreatedAt.IsZero() {
				t.Error("CreatedAt not set")
			}
		}
	}

	if !foundNewSecret {
		t.Error("new secret not found in updated secrets list")
	}
}

func TestAddServiceSecret_MissingClientID(t *testing.T) {
	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {},
		},
	}

	_, err := service_keys.AddServiceSecret(db, "missing-client-id", "admin")
	if err == nil {
		t.Fatal("expected error when service key not found")
	}
}
