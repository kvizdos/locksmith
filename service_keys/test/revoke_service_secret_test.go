package service_keys_test

import (
	"testing"
	"time"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/service_keys"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRevokeServiceSecret_Success(t *testing.T) {
	now := time.Now().UTC()

	secretID := "active-secret-id"
	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"svc-id": service_keys.ServiceKey{
					ID:       primitive.NewObjectID(),
					ClientID: "svc-test-client",
					Name:     "Test Service",
					Scopes:   []string{"records.read"},
					Secrets: []service_keys.ServiceSecret{
						{
							ID:         secretID,
							SecretHash: "hash",
							CreatedAt:  now,
							CreatedBy:  "test-user",
							RevokedAt:  nil,
						},
					},
					CreatedAt: now,
					CreatedBy: "test-user",
				}.ToMap(),
			},
		},
	}

	err := service_keys.RevokeServiceSecret(db, "svc-test-client", secretID, "admin")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	raw, found := db.FindOne("service_keys", map[string]any{"client_id": "svc-test-client"})
	if !found {
		t.Fatal("expected service key to still exist")
	}

	updated := service_keys.ServiceKeyFromMap(raw)

	if len(updated.Secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(updated.Secrets))
	}

	if updated.Secrets[0].RevokedAt == nil {
		t.Error("expected RevokedAt to be set")
	}
}

func TestRevokeServiceSecret_AlreadyRevoked(t *testing.T) {
	now := time.Now().UTC()
	revokedAt := now.Add(-1 * time.Hour)

	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"svc-id": service_keys.ServiceKey{
					ID:       primitive.NewObjectID(),
					ClientID: "svc-test-client",
					Name:     "Test Service",
					Secrets: []service_keys.ServiceSecret{
						{
							ID:         "revoked-secret",
							SecretHash: "hash",
							CreatedAt:  now,
							RevokedAt:  &revokedAt,
						},
					},
					CreatedAt: now,
				}.ToMap(),
			},
		},
	}

	err := service_keys.RevokeServiceSecret(db, "svc-test-client", "revoked-secret", "admin")
	if err == nil {
		t.Fatal("expected error for already revoked secret")
	}
}

func TestRevokeServiceSecret_SecretNotFound(t *testing.T) {
	now := time.Now().UTC()

	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"svc-id": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-test-client",
					Name:      "Test Service",
					Secrets:   []service_keys.ServiceSecret{},
					CreatedAt: now,
				}.ToMap(),
			},
		},
	}

	err := service_keys.RevokeServiceSecret(db, "svc-test-client", "non-existent-secret", "admin")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestRevokeServiceSecret_ClientNotFound(t *testing.T) {
	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {},
		},
	}

	err := service_keys.RevokeServiceSecret(db, "missing-client", "any-secret", "admin")
	if err == nil {
		t.Fatal("expected error for missing client ID")
	}
}
