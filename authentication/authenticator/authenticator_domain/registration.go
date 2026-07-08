package authenticator_domain

import (
	"net/http"
)

type RegistrationHandler interface {
	CanHandle(r *http.Request) bool
	Register(r *http.Request) bool
}
