package authorizer

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_methods"
	"github.com/kvizdos/locksmith/authentication/authorizer/internal/method_oidc"
	"github.com/kvizdos/locksmith/authentication/authorizer/internal/method_password"
)

type AuthorizerHandler interface {
	ServeLoginAPI(w http.ResponseWriter, r *http.Request)
	ServeProviderStartAPI(w http.ResponseWriter, r *http.Request)
}

func AllowMethodPassword(opts ...authorizer_methods.PasswordValidatorOption) authorizer_domain.Handler {
	return method_password.NewPasswordValidator(opts...)
}
func AllowMethodOIDC(opts ...authorizer_methods.OIDCValidatorOption) authorizer_domain.Handler {
	return method_oidc.NewOIDCValidator(opts...)
}

func (a *authorizers) getHandler(r *http.Request) (authorizer_domain.Handler, error) {
	providerMatch := r.PathValue("provider")
	handlerHint := r.Header.Get("X-Handler-Hint")

	if providerMatch == "" && handlerHint != "" {
		providerMatch = handlerHint
	}

	for _, method := range a.methods {
		if providerMatch != "" && method.Name() == providerMatch {
			return method, nil
		}
		if providerMatch == "" && method.CanHandle(r) {
			return method, nil
		}
	}

	return nil, authorizer_domain.ErrUnhandleableRequest
}
