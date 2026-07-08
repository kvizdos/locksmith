package tokens

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
	"github.com/kvizdos/locksmith/users"
)

type TokenManager interface {
	Read(r *http.Request) (*authorizer_domain.Token, error)

	PassToClient(w http.ResponseWriter, r *http.Request, token *authorizer_domain.Token) error
	CreateAuthToken(users.LocksmithUserInterface) (*authorizer_domain.Token, error)
}
