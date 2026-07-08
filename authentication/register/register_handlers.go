package register

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/register/internal/method_hint"
	"github.com/kvizdos/locksmith/authentication/register/internal/method_password"
	"github.com/kvizdos/locksmith/authentication/register/register_domain"
	"github.com/kvizdos/locksmith/authentication/register/register_methods"
)

// RegistrarHandler is the public interface exposed by a *registrar,
// mirroring authenticator.AuthorizerHandler.
type RegistrarHandler interface {
	ServeRegisterAPI(w http.ResponseWriter, r *http.Request)
}

func AllowMethodPassword(opts ...register_methods.PasswordOption) register_domain.Handler {
	return method_password.NewPasswordRegistration(opts...)
}

func AllowMethodHint(opts ...register_methods.HintOption) register_domain.Handler {
	return method_hint.NewHintRegistration(opts...)
}

// getHandler finds the first registered method that can handle the given
// request. Order matters: methods that key off of specific request
// characteristics (e.g. hint cookie/header) should be registered before
// broader fallbacks (e.g. "any POST") via WithMethods.
func (r *registrar) getHandler(req *http.Request) (register_domain.Handler, error) {
	for _, method := range r.methods {
		if method.CanHandle(req) {
			return method, nil
		}
	}
	return nil, register_domain.ErrUnhandleableRequest
}
