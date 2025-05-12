package roles

import (
	"fmt"
	"slices"
)

var AVAILABLE_PERMISSIONS = map[string]Permission{}

func GetExposable(perms []string) []string {
	out := []string{}
	for _, perm := range perms {
		if v, ok := AVAILABLE_PERMISSIONS[perm]; ok {
			if !v.DontExposeFrontend {
				out = append(out, perm)
			}
		}
	}
	return out
}

// Define required admin role for Locksmith to work
func init() {
	CreatePermissionSet([]Permission{
		{
			Permission:         "human",
			IsElevated:         true,
			JWTOnly:            false,
			AvailableRoles:     []string{"admin", "user"},
			DontExposeFrontend: true,
		},
		{
			Permission:     "user.delete.self",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin", "user"},
		},
		{
			Permission:     "view.ls-admin",
			IsElevated:     true,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
		{
			Permission:     "user.invite",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
		{
			Permission:     "user.delete.other",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
		{
			Permission:     "users.list.all",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
		{
			Permission:     "users.lock",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
		{
			Permission:     "users.lock.manage",
			IsElevated:     false,
			JWTOnly:        false,
			AvailableRoles: []string{"admin"},
		},
	})
}

type Permission struct {
	// Action Scope, e.g. "users.read", "users.create", etc.
	// Best practice to use this format, so the UI can be easily
	// visualized per resource.
	Permission string
	// Will cause a warning to show when adding to a service key.
	IsElevated bool
	// Specialized only for when you want a permission to be retrieved
	// by a JWT. Users with this permission may request this scope and
	// pass the Authentication Bearer header, but it will not be
	// included by default.
	JWTOnly bool
	// Which roles is this available to?
	AvailableRoles []string
	// Disable the permission from being exposed on the /me endpoint.
	DontExposeFrontend bool
}

func CreatePermissionSet(permissions []Permission) {
	if AVAILABLE_ROLES == nil {
		AVAILABLE_ROLES = map[string][]string{}
	}
	if AVAILABLE_PERMISSIONS == nil {
		AVAILABLE_PERMISSIONS = map[string]Permission{}
	}
	for _, permission := range permissions {
		if !permission.JWTOnly {
			for _, role := range permission.AvailableRoles {
				if _, ok := AVAILABLE_ROLES[role]; !ok {
					AVAILABLE_ROLES[role] = []string{permission.Permission}
					continue
				}

				if slices.Contains(AVAILABLE_ROLES[role], permission.Permission) {
					continue
				}
				AVAILABLE_ROLES[role] = append(AVAILABLE_ROLES[role], permission.Permission)
			}
		}

		if _, ok := AVAILABLE_ROLES[permission.Permission]; ok {
			panic(fmt.Sprintf("Duplicate setting of %s", permission.Permission))
		}
		AVAILABLE_PERMISSIONS[permission.Permission] = permission
	}
}
