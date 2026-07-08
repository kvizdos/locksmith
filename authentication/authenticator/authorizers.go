package authenticator

import (
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"github.com/kvizdos/locksmith/authentication/authenticator/internal/method_oidc"
	"github.com/kvizdos/locksmith/authentication/authenticator/internal/method_password"
)

type AuthorizerHandler interface {
	ServeLoginAPI(w http.ResponseWriter, r *http.Request)
	ServeProviderStartAPI(w http.ResponseWriter, r *http.Request)
}

func AllowMethodPassword(opts ...authenticator_methods.PasswordValidatorOption) authenticator_domain.Handler {
	return method_password.NewPasswordValidator(opts...)
}
func AllowMethodOIDC(opts ...authenticator_methods.OIDCValidatorOption) authenticator_domain.Handler {
	return method_oidc.NewOIDCValidator(opts...)
}

func (a *authorizers) getHandler(r *http.Request) (authenticator_domain.Handler, error) {
	providerMatch := r.PathValue("provider")
	handlerHint := r.Header.Get("X-Handler-Hint")

	if providerMatch == "" && r.URL.Query().Get("provider") != "" {
		providerMatch = r.URL.Query().Get("provider")
	}
	if providerMatch == "" && handlerHint != "" {
		providerMatch = handlerHint
	}

	for _, method := range a.methods {
		fmt.Println(providerMatch, method.Name())
		if providerMatch != "" && method.Name() == providerMatch {
			return method, nil
		}
		if providerMatch == "" && method.CanHandle(r) {
			return method, nil
		}
	}

	return nil, authenticator_domain.ErrUnhandleableRequest
}
