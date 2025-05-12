package service_keys_test

import (
	"testing"
	"time"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/service_keys"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestListServiceKeysByCreatedBy(t *testing.T) {
	now := time.Now().UTC()

	db := &database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"key-1": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-a",
					Name:      "Key A",
					CreatedBy: "admin-1",
					CreatedAt: now,
				}.ToMap(),
				"key-2": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-b",
					Name:      "Key B",
					CreatedBy: "admin-1",
					CreatedAt: now,
				}.ToMap(),
				"key-3": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-c",
					Name:      "Key C",
					CreatedBy: "admin-2",
					CreatedAt: now,
				}.ToMap(),
			},
		},
	}

	keys := service_keys.ListServiceKeys(db, "created_by", "admin-1")

	if len(keys) != 2 {
		t.Errorf("expected 2 service keys, got %d", len(keys))
	}

	for _, key := range keys {
		if key.CreatedBy != "admin-1" {
			t.Errorf("expected CreatedBy to be admin-1, got %s", key.CreatedBy)
		}
	}
}

func TestListServiceKeys_Empty(t *testing.T) {
	db := &database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {},
		},
	}

	keys := service_keys.ListServiceKeys(db, "created_by", "nobody")

	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestListServiceKeys_AllKeys(t *testing.T) {
	now := time.Now().UTC()

	db := &database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"a": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-a",
					Name:      "Service A",
					CreatedAt: now,
					CreatedBy: "admin-1",
				}.ToMap(),
				"b": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-b",
					Name:      "Service B",
					CreatedAt: now,
					CreatedBy: "admin-2",
				}.ToMap(),
				"c": service_keys.ServiceKey{
					ID:        primitive.NewObjectID(),
					ClientID:  "svc-c",
					Name:      "Service C",
					CreatedAt: now,
					CreatedBy: "admin-3",
				}.ToMap(),
			},
		},
	}

	keys := service_keys.ListServiceKeys(db, "", nil)

	if len(keys) != 3 {
		t.Errorf("expected 3 service keys, got %d", len(keys))
	}

	expectedClients := map[string]bool{
		"svc-a": true,
		"svc-b": true,
		"svc-c": true,
	}

	for _, key := range keys {
		if !expectedClients[key.ClientID] {
			t.Errorf("unexpected key: %s", key.ClientID)
		}
	}
}
