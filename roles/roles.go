package roles

import "fmt"

type RoleInfo struct {
	FrontendPermissions []string
	BackendPermissions  []string
}

var AVAILABLE_ROLES map[string]RoleInfo

// Define required admin role for Locksmith to work
func init() {
	adminRole := Role{
		Name: "AUTHENTICATION.admin",
		Permissions: []string{
			"AUTHENTICATION.view.ls-admin",
			"AUTHENTICATION.user.invite",
			"AUTHENTICATION.user.delete.self",
			"AUTHENTICATION.user.delete.other",
			"AUTHENTICATION.users.list.all",
			"AUTHENTICATION.users.lock",        // get lock status
			"AUTHENTICATION.users.lock.manage", // set lock state
		},
	}
	AddRole(adminRole)

	userRole := Role{
		Name: "AUTHENTICATION.user",
		Permissions: []string{
			"AUTHENTICATION.user.delete.self",
		},
	}
	AddRole(userRole)
}

type Role struct {
	Name                string   `json:"name" bson:"name"`
	Permissions         []string `json:"permissions" bson:"permissions"`
	FrontendPermissions []string `json:"frontend" bson:"frontend"`
}

func (r Role) HasPermission(permission string) bool {
	for _, perm := range r.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

func (r Role) HasFrontendPermission(permission string) bool {
	for _, perm := range r.FrontendPermissions {
		if perm == permission {
			return true
		}
	}
	return false
}

func AddRole(role Role) {
	if AVAILABLE_ROLES == nil {
		AVAILABLE_ROLES = map[string]RoleInfo{}
	}
	AVAILABLE_ROLES[role.Name] = RoleInfo{
		FrontendPermissions: role.FrontendPermissions,
		BackendPermissions:  role.Permissions,
	}
}

func RoleExists(roleName string) bool {
	_, exists := AVAILABLE_ROLES[roleName]
	return exists
}

func GetRole(roleName string) (Role, error) {
	perms, exists := AVAILABLE_ROLES[roleName]

	if !exists {
		return Role{}, fmt.Errorf("invalid role")
	}

	return Role{
		Name:                roleName,
		Permissions:         perms.BackendPermissions,
		FrontendPermissions: perms.FrontendPermissions,
	}, nil
}

// Useful to add roles to the default Locksmith Admin role
// not many other uses..
func AddPermissionsToRole(roleName string, permissions []string) error {
	role, err := GetRole(roleName)

	if err != nil {
		return err
	}

	newPerms := role.Permissions

	for _, permission := range permissions {
		if role.HasPermission(permission) {
			continue
		}

		newPerms = append(newPerms, permission)
	}

	AVAILABLE_ROLES[roleName] = RoleInfo{
		FrontendPermissions: role.FrontendPermissions,
		BackendPermissions:  newPerms,
	}

	return nil
}

func AddFrontendPermissionsToRole(roleName string, permissions []string) error {
	role, err := GetRole(roleName)

	if err != nil {
		return err
	}

	newPerms := role.FrontendPermissions

	for _, permission := range permissions {
		if role.HasFrontendPermission(permission) {
			continue
		}

		newPerms = append(newPerms, permission)
	}

	AVAILABLE_ROLES[roleName] = RoleInfo{
		FrontendPermissions: newPerms,
		BackendPermissions:  role.Permissions,
	}

	return nil
}

func GetPermissionsForRole(roleName string) ([]string, error) {
	permissions, exists := AVAILABLE_ROLES[roleName]

	if !exists {
		return []string{}, fmt.Errorf("invalid role name")
	}

	return permissions.BackendPermissions, nil
}

func GetFrontendPermissionsForRole(roleName string) ([]string, error) {
	permissions, exists := AVAILABLE_ROLES[roleName]

	if !exists {
		return []string{}, fmt.Errorf("invalid role name")
	}

	return permissions.FrontendPermissions, nil
}
