package method_oidc

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"github.com/kvizdos/locksmith/database"
)

func NewOIDCValidator(opts ...authenticator_methods.OIDCValidatorOption) authenticator_domain.Handler {
	pv := oidcHandler{}
	for _, opt := range opts {
		opt(&pv.options)
	}

	if pv.options.ProviderName == "" {
		panic("config is required")
	}

	return pv
}

type oidcHandler struct {
	options authenticator_methods.OIDCValidatorOptions
}

func (pv oidcHandler) CanHandle(r *http.Request) bool {
	return detectFlow(r) != flowNone
}

func (pv oidcHandler) GetIconBytes() []byte {
	return pv.options.LogoBytes
}

func (pv oidcHandler) Name() string {
	return pv.options.ProviderName
}

func (pv oidcHandler) Passwordless() bool {
	return true
}

func (pv oidcHandler) Session(db database.DatabaseAccessor) authenticator_domain.Session {
	return newOIDCValidationSession(db, pv.options)
}
