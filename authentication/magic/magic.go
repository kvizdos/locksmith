package magic

import (
	"time"

	"github.com/kvizdos/locksmith/authentication/signing"
)

/*
The "magic.go" file serves as the entrypoint for the Magic Tokens functionality within Locksmith.
Its primary purpose is to manage Magic Access Codes (MACs) that grant users temporary, scoped access to specific app areas.
Key Use Cases:
- Providing users with seamless access for password resets or notification-based URLs.
- Ensuring secure yet convenient access for users, especially for actions requiring limited permissions.
- Implementing Time-To-Live (TTL) for tokens, adding a layer of security by ensuring time-bound accessibility.

Note: It's crucial to understand the security implications of the Magic Tokens feature and ensure its proper configuration and deployment.
*/

var MagicSigningPackage signing.SigningPackageInterface

type MagicAuthenticationVariables struct {
	ForUserID          string
	AllowedPermissions []string      // what the user can DO
	TTL                time.Duration // how long the token should live for
}

type MagicAuthentication struct {
	Code               string   // authentication code (hashed)
	AllowedPermissions []string // what the user can DO
	ExpiresAt          int64    // when the token expires
	InheritRole        string   // filled in at validation time
	Username           string   // filled in at validation time
}

func (m MagicAuthentication) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":        m.Code,
		"permissions": m.AllowedPermissions,
		"expires":     m.ExpiresAt,
	}
}

type MagicAuthentications []MagicAuthentication

func (m MagicAuthentications) ToMap() []interface{} {
	out := make([]interface{}, len(m))

	for i, magic := range m {
		out[i] = magic.ToMap()
	}

	return out
}
