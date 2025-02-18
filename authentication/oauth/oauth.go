package oauth

import (
	"net/http"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type OAuthProviders []OAuthProvider

func (o OAuthProviders) GetNames() []string {
	out := make([]string, len(o))
	for i, oauth := range o {
		out[i] = oauth.GetName()
	}
	return out
}

type OAuthProvider interface {
	RegisterRoutes(apiMux *http.ServeMux)
	GetName() string
	GetUserByEmail(email string, r *http.Request, customUserInterface ...users.LocksmithUserInterface) users.LocksmithUserInterface
}

type BaseOAuthProvider struct {
	ClientID                     string
	ClientSecret                 string
	BaseURL                      string
	Scopes                       string
	Database                     database.DatabaseAccessor
	RedirectToRegisterOnNotFound bool
}

func (g BaseOAuthProvider) GetName() string {
	panic("OAuthProvider requires a GetName() function!")
}

func (g BaseOAuthProvider) RegisterRoutes(apiMux *http.ServeMux) {
	panic("OAuthProvider requires a RegisterRoutes(mux) function!")
}

func (g BaseOAuthProvider) GetUserByEmail(email string, r *http.Request, customUserInterface ...users.LocksmithUserInterface) users.LocksmithUserInterface {
	rawUser, found := g.Database.FindOne("users", map[string]interface{}{
		"email": email,
	})

	if !found {
		return nil
	}

	var lsu users.LocksmithUserInterface
	if len(customUserInterface) > 0 {
		lsu = customUserInterface[0]
	} else {
		lsu = users.LocksmithUser{}
	}

	var tmpUser users.LocksmithUserInterface
	lsu.ReadFromMap(&tmpUser, rawUser.(map[string]interface{}))

	return tmpUser

}
