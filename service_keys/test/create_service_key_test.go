package service_keys_test

import (
	"testing"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/service_keys"
)

func TestCreateServiceKey_Success(t *testing.T) {
	roles.AVAILABLE_PERMISSIONS = map[string]roles.Permission{
		"students.read": {
			Permission:     "students.read",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
	}

	testDb := &database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {},
		},
	}

	creds, err := service_keys.CreateServiceKey(testDb, service_keys.CreateServiceKeyOptions[any]{
		FriendlyName:    "Demo Key",
		Scopes:          []string{"students.read"},
		CreatedByUserID: "test_user_id",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if creds.ClientID == "" {
		t.Errorf("expected clientID to be set")
	}
	if len(creds.Secret) < 50 {
		t.Errorf("expected long secret, got: %d characters", len(creds.Secret))
	}
}

func TestCreateServiceKey_DuplicateName(t *testing.T) {
	roles.AVAILABLE_PERMISSIONS = map[string]roles.Permission{
		"students.read": {
			Permission:     "students.read",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
	}

	testDb := &database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {
				"1": map[string]interface{}{"name": "Demo Key"},
			},
		},
	}

	_, err := service_keys.CreateServiceKey(testDb, service_keys.CreateServiceKeyOptions[any]{
		FriendlyName:    "Demo Key",
		Scopes:          []string{"students.read"},
		CreatedByUserID: "test_user_id",
	})
	if err == nil {
		t.Fatal("expected error due to duplicate name")
	}
}

func TestCreateServiceKey_InvalidScope(t *testing.T) {
	roles.AVAILABLE_PERMISSIONS = map[string]roles.Permission{}

	testDb := &database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {},
		},
	}

	_, err := service_keys.CreateServiceKey(testDb, service_keys.CreateServiceKeyOptions[any]{
		FriendlyName:    "Key with bad scope",
		Scopes:          []string{"invalid.scope"},
		CreatedByUserID: "test_user_id",
	})
	if err == nil {
		t.Fatal("expected error due to invalid scope")
	}
}

func TestCreateServiceKey_DuplicateScopes(t *testing.T) {
	roles.AVAILABLE_PERMISSIONS = map[string]roles.Permission{
		"students.read": {
			Permission:     "students.read",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
	}

	testDb := &database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"service_keys": {},
		},
	}

	_, err := service_keys.CreateServiceKey(testDb, service_keys.CreateServiceKeyOptions[any]{
		FriendlyName:    "Key with dups",
		Scopes:          []string{"students.read", "students.read"},
		CreatedByUserID: "test_user_id",
	})
	if err == nil {
		t.Fatal("expected error due to duplicate scope")
	}
}
