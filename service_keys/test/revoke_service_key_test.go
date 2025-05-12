package service_keys_test

import (
	"testing"
	"time"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/service_keys"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRevokeServiceKey_Success(t *testing.T) {
	now := time.Now().UTC()

	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"demo-key": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-test-client",
					Name:      "My Service",
					CreatedAt: now,
					CreatedBy: "admin",
				}.ToMap(),
			},
		},
	}

	err := service_keys.RevokeServiceKey(db, "svc-test-client", "admin")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	raw, found := db.FindOne("service_keys", map[string]any{"client_id": "svc-test-client"})
	if !found {
		t.Fatal("expected key to still exist")
	}

	updated := service_keys.ServiceKeyFromMap(raw)
	if updated.RevokedAt == nil {
		t.Error("expected RevokedAt to be set")
	}
	if updated.RevokedBy != "admin" {
		t.Errorf("expected RevokedBy to be 'admin', got '%s'", updated.RevokedBy)
	}
}

func TestRevokeServiceKey_AlreadyRevoked(t *testing.T) {
	now := time.Now().UTC()
	revoked := now.Add(-1 * time.Hour)

	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"revoked-key": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-already-dead",
					Name:      "Dead Service",
					CreatedAt: now,
					CreatedBy: "admin",
					RevokedAt: &revoked,
					RevokedBy: "admin",
				}.ToMap(),
			},
		},
	}

	err := service_keys.RevokeServiceKey(db, "svc-already-dead", "admin")
	if err == nil {
		t.Fatal("expected error for already revoked key")
	}
}

func TestRevokeServiceKey_ClientNotFound(t *testing.T) {
	db := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {},
		},
	}

	err := service_keys.RevokeServiceKey(db, "does-not-exist", "admin")
	if err == nil {
		t.Fatal("expected error for missing service key")
	}
}
