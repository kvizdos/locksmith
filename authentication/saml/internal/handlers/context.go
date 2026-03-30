package saml_handlers

import saml_auth "github.com/kvizdos/locksmith/authentication/saml/internal/auth"

type SAMLCtxKey struct{}

type SAMLContext struct {
	AuthnRequest *saml_auth.AuthnRequest
	RelayState   string
	Validated    *saml_auth.ValidatedRequest
}
