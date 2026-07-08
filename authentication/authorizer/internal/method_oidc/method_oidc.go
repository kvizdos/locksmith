package method_oidc

import (
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_methods"
	"github.com/kvizdos/locksmith/database"
)

func NewOIDCValidator(opts ...authorizer_methods.OIDCValidatorOption) authorizer_domain.Handler {
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
	options authorizer_methods.OIDCValidatorOptions
}

func (pv oidcHandler) CanHandle(r *http.Request) bool {
	return detectFlow(r) != flowNone
}

func (pv oidcHandler) Name() string {
	return fmt.Sprintf("oidc-%s", pv.options.ProviderName)
}

func (pv oidcHandler) Passwordless() bool {
	return true
}

func (pv oidcHandler) Session(db database.DatabaseAccessor) authorizer_domain.Session {
	return newOIDCValidationSession(db, pv.options)
}
