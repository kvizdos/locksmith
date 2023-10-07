package magic

import (
	"time"

	"github.com/kvizdos/locksmith/authentication/signing"
)

var MagicSigningPackage signing.SigningPackage

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
}

func (m MagicAuthentication) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":        m.Code,
		"permissions": m.AllowedPermissions,
		"expires":     m.ExpiresAt,
	}
}

type MagicAuthentications []MagicAuthentication

func (m MagicAuthentications) ToMap() []map[string]interface{} {
	out := make([]map[string]interface{}, len(m))

	for i, magic := range m {
		out[i] = magic.ToMap()
	}

	return out
}
