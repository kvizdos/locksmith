package tokens

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/users"
)

type TokenManager interface {
	Read(r *http.Request) (*authenticator_domain.Token, error)

	PassToClient(w http.ResponseWriter, r *http.Request, token *authenticator_domain.Token) error
	CreateAuthToken(users.LocksmithUserInterface) (*authenticator_domain.Token, error)
}
