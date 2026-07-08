package method_oidc

import (
	"github.com/coreos/go-oidc"
	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_methods"
	"github.com/kvizdos/locksmith/authentication/authorizer/internal/sessions"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func newOIDCValidationSession(db database.DatabaseAccessor, opts authorizer_methods.OIDCValidatorOptions) *oidcValidationSession {
	return &oidcValidationSession{
		db: db,
		LinkedSession: sessions.NewLinkedSession(
			sessions.WithProvider(opts.ProviderName),
		),
		options: opts,
	}
}

type oidcValidationSession struct {
	sessions.LinkedSession

	// Common
	db      database.DatabaseAccessor
	options authorizer_methods.OIDCValidatorOptions

	// Identity Resolution Context
	flow                     oidcFlow
	untrustedParsedCode      string
	pkceVerifier             string
	untrustedCredentialToken string

	// Authorization Context
	authoritativeToken *oidc.IDToken
	displayName        string
}

func (pv oidcValidationSession) IsAuthorized(user users.LocksmithUserInterface) error {

	return nil
}
