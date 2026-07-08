package method_oidc

import (
	"github.com/coreos/go-oidc"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"github.com/kvizdos/locksmith/authentication/authenticator/internal/sessions"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func newOIDCValidationSession(db database.DatabaseAccessor, opts authenticator_methods.OIDCValidatorOptions) *oidcValidationSession {
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
	options authenticator_methods.OIDCValidatorOptions

	// Identity Resolution Context
	flow                     oidcFlow
	untrustedParsedCode      string
	pkceVerifier             string
	untrustedCredentialToken string
	selectBy                 string
	redirectTarget           string

	// Authorization Context
	authoritativeToken *oidc.IDToken
	displayName        string
}

// Identity is authenticated via ResolveIdentity; nothing further to check here.
func (pv oidcValidationSession) IsAuthorized(_ users.LocksmithUserInterface) error {
	return nil
}
