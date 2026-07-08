package method_password

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_methods"
	"github.com/kvizdos/locksmith/database"
)

func NewPasswordValidator(opts ...authorizer_methods.PasswordValidatorOption) authorizer_domain.Handler {
	pv := passwordValidator{}
	for _, opt := range opts {
		opt(&pv.options)
	}
	return pv
}

type passwordValidator struct {
	options authorizer_methods.PasswordValidatorOptions
}

func (pv passwordValidator) CanHandle(r *http.Request) bool {
	return r.Method == http.MethodPost
}

func (pv passwordValidator) Name() string {
	return "password"
}

func (pv passwordValidator) Passwordless() bool {
	return false
}

func (pv passwordValidator) Session(db database.DatabaseAccessor) authorizer_domain.Session {
	return newPasswordValidatorSession(db, pv.options)
}
