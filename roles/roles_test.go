package roles

import (
	"testing"
)

func TestRoleExistsDoesNot(t *testing.T) {
	exists := RoleExists("test")

	if exists == true {
		t.Error("role does not exist")
	}
}

func TestRoleExistsDoes(t *testing.T) {
	AVAILABLE_ROLES = map[string]RoleInfo{
		"admin": {
			BackendPermissions: []string{"view.admin", "user.delete.self"},
		},
	}

	exists := RoleExists("admin")

	if exists == false {
		t.Error("role does exist")
	}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	})
}

func TestRoleHasPermissionDoes(t *testing.T) {
	role := Role{
		Name: "admin",
		Permissions: []string{
			"view.admin",
			"user.delete.self",
			"user.delete.other",
		},
	}

	hasPermission := role.HasPermission("user.delete.other")

	if !hasPermission {
		t.Error("role should have permission")
	}
}

func TestRoleHasPermissionDoesNot(t *testing.T) {
	role := Role{
		Name: "admin",
		Permissions: []string{
			"view.admin",
			"user.delete.self",
		},
	}

	hasPermission := role.HasPermission("user.delete.other")

	if hasPermission {
		t.Error("role should not have permission")
	}
}

func TestAddRole(t *testing.T) {
	AddRole(Role{
		Name: "admin",
		Permissions: []string{
			"view.admin",
			"user.delete.others",
			"user.delete.self",
		},
	})

	permissions, exists := AVAILABLE_ROLES["admin"]

	if !exists {
		t.Errorf("could not find role name")
		return
	}

	if len(permissions.BackendPermissions) != 3 {
		t.Errorf("could not find permissions, found %d expected %d", 2, len(permissions.BackendPermissions))
	}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	})
}

func TestGetRoleInvalidRole(t *testing.T) {
	AVAILABLE_ROLES = map[string]RoleInfo{
		"admin": {
			BackendPermissions: []string{"view.admin", "user.delete.self"},
		},
	}

	_, err := GetRole("user")

	if err == nil {
		t.Error("expected to receive an error")
		return
	}

	if err.Error() != "invalid role" {
		t.Errorf("received unexpected error: %s", err.Error())
	}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	})
}

func TestGetRole(t *testing.T) {
	AVAILABLE_ROLES = map[string]RoleInfo{
		"admin": {
			BackendPermissions: []string{"view.admin", "user.delete.self"},
		},
	}

	role, err := GetRole("admin")

	if err != nil {
		t.Errorf("received unexpected error: %s", err.Error())
		return
	}

	if role.Name != "admin" {
		t.Errorf("invalid role name: %s", role.Name)
	}

	if len(role.Permissions) != 2 {
		t.Errorf("received invalid number of permissions: %d", len(role.Permissions))
	}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	})
}

func TestAddPermissionToRoleFailsIfInvalidRole(t *testing.T) {
	AVAILABLE_ROLES = map[string]RoleInfo{
		"admin": RoleInfo{
			BackendPermissions: []string{"view.admin"},
		},
	}

	err := AddPermissionsToRole("admins", []string{"users.delete.other", "user.delete.self"})

	if err == nil {
		t.Error("should have thrown an error")
		return
	}

	if err.Error() != "invalid role" {
		t.Errorf("received invalid error: %s", err.Error())
	}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	})
}

func TestAddPermissionToRoleContinuesIfPermissionExists(t *testing.T) {
	AVAILABLE_ROLES = map[string]RoleInfo{
		"admin": RoleInfo{
			BackendPermissions: []string{"view.admin"},
		},
	}

	err := AddPermissionsToRole("admin", []string{"view.admin", "users.delete.other", "user.delete.self"})

	if err != nil {
		t.Errorf("shouldn't have thrown an error: %s", err.Error())
		return
	}

	perms := AVAILABLE_ROLES["admin"]

	if len(perms.BackendPermissions) != 3 {
		t.Errorf("found invalid amount of permissions, found %d expected %d", len(perms.BackendPermissions), 3)
	}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	})
}

func TestAddPermissionToRole(t *testing.T) {
	AVAILABLE_ROLES = map[string]RoleInfo{
		"admin": RoleInfo{
			BackendPermissions: []string{"view.admin"},
		},
	}

	err := AddPermissionsToRole("admin", []string{"users.delete.other", "user.delete.self"})

	if err != nil {
		t.Error(err)
		return
	}

	perms := AVAILABLE_ROLES["admin"]

	if len(perms.BackendPermissions) != 3 {
		t.Errorf("found invalid amount of permissions, found %d expected %d", len(perms.BackendPermissions), 3)
	}

	t.Cleanup(func() {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	})
}
