package roles

import (
	"testing"
)

func TestCreatePermissionSet_Basic(t *testing.T) {
	AVAILABLE_ROLES = map[string][]string{}
	AVAILABLE_PERMISSIONS = map[string]Permission{}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string][]string{}
		AVAILABLE_PERMISSIONS = map[string]Permission{}
	})

	permissions := []Permission{
		{
			Permission:     "users.read",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin", "viewer"},
		},
		{
			Permission:     "users.create",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
		{
			Permission:     "billing.export",
			IsElevated:     true,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
		{
			Permission:     "token.issue",
			IsElevated:     false,
			JWTOnly:        true,
			AvailableRoles: []string{"admin"},
		},
	}

	CreatePermissionSet(permissions)

	if len(AVAILABLE_PERMISSIONS) != 4 {
		t.Errorf("expected 4 permissions, got %d", len(AVAILABLE_PERMISSIONS))
	}

	adminPerms, ok := AVAILABLE_ROLES["admin"]
	if !ok {
		t.Fatal("admin role not found")
	}

	expected := []string{"users.read", "users.create", "billing.export"}
	for _, exp := range expected {
		found := false
		for _, p := range adminPerms {
			if p == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected permission %s in admin role", exp)
		}
	}

	for _, p := range adminPerms {
		if p == "token.issue" {
			t.Error("jwtOnly permission should not be added to role")
		}
	}

	viewerPerms := AVAILABLE_ROLES["viewer"]
	if len(viewerPerms) != 1 || viewerPerms[0] != "users.read" {
		t.Errorf("viewer role should only have users.read, got %+v", viewerPerms)
	}
}
