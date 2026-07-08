package register_domain

import (
	"net/http"

	"github.com/kvizdos/locksmith/database"
)

type Handler interface {
	CanHandle(r *http.Request) bool
	Name() string
	Session(db database.DatabaseAccessor) Session
}

type Session interface {
	LoadRequest(r *http.Request) error
	RegistrationRequest() Request
}
