package tokens

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type cookieManager struct {
	db           database.DatabaseAccessor
	redirectPath string
}

func NewCookieManager(db database.DatabaseAccessor, redirectPath string) *cookieManager {
	return &cookieManager{db: db, redirectPath: redirectPath}
}

func (c *cookieManager) Read(r *http.Request) (*authorizer_domain.Token, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *cookieManager) CreateAuthToken(user users.LocksmithUserInterface) (*authorizer_domain.Token, error) {
	token, err := user.GeneratePasswordSession()
	if err != nil {
		return nil, fmt.Errorf("failed to generate password session: %w", err)
	}
	var u *users.LocksmithUser
	if us, ok := user.(users.LocksmithUser); ok {
		u = &us
	}

	err = user.SavePasswordSession(token, c.db)

	if err != nil {
		return nil, fmt.Errorf("failed to save password session: %w", err)
	}

	return &authorizer_domain.Token{
		User:      u,
		AuthToken: token.Token,
		ExpiresAt: time.Unix(token.ExpiresAt, 0).UTC(),
	}, nil
}

func (c *cookieManager) PassToClient(w http.ResponseWriter, r *http.Request, token *authorizer_domain.Token) error {
	cookieValue := token.User.GenerateCookieValueFromSession(authentication.PasswordSession{
		Token:     token.AuthToken,
		ExpiresAt: token.ExpiresAt.Unix(),
	})

	cookie := http.Cookie{Name: "token", Value: cookieValue, Expires: time.Unix(token.ExpiresAt.Unix(), 0), HttpOnly: true, Secure: true, Path: "/"}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, c.redirectPath, http.StatusSeeOther)
	return nil
}
