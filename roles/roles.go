package roles

import "fmt"

var AVAILABLE_ROLES map[string][]string

// Define required admin role for Locksmith to work
func init() {
	adminRole := Role{
		Name: "admin",
		Permissions: []string{
			"view.ls-admin",
			"user.invite",
			"user.delete.self",
			"user.delete.other",
			"users.list.all",
		},
	}
	AddRole(adminRole)

	userRole := Role{
		Name: "user",
		Permissions: []string{
			"user.delete.self",
		},
	}
	AddRole(userRole)
}

type Role struct {
	Name        string   `json:"name" bson:"name"`
	Permissions []string `json:"permissions" bson:"permissions"`
}

func (r Role) HasPermission(permission string) bool {
	for _, perm := range r.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

func AddRole(role Role) {
	if AVAILABLE_ROLES == nil {
		AVAILABLE_ROLES = map[string][]string{}
	}
	AVAILABLE_ROLES[role.Name] = role.Permissions
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
		Name:        roleName,
		Permissions: perms,
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

	AVAILABLE_ROLES[roleName] = newPerms

	return nil
}

func GetPermissionsForRole(roleName string) ([]string, error) {
	permissions, exists := AVAILABLE_ROLES[roleName]

	if !exists {
		return []string{}, fmt.Errorf("invalid role name")
	}

	return permissions, nil
}
