package tokens

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
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

func (c *cookieManager) Read(r *http.Request) (*authenticator_domain.Token, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *cookieManager) CreateAuthToken(user users.LocksmithUserInterface) (*authenticator_domain.Token, error) {
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

	return &authenticator_domain.Token{
		User:      u,
		AuthToken: token.Token,
		ExpiresAt: time.Unix(token.ExpiresAt, 0).UTC(),
	}, nil
}

func SetBaseCookies(w http.ResponseWriter, token *authenticator_domain.Token) {
	sessionExpiresAtCookie := http.Cookie{Name: "ls_expires_at", Value: fmt.Sprintf("%d", token.ExpiresAt.Unix()), Expires: time.Unix(token.ExpiresAt.Unix(), 0), HttpOnly: false, Secure: true, Path: "/"}

	if token.OAuthProvider == "" {
		oauthProviderCookie := http.Cookie{Name: "ls_oauth_provider", Value: token.OAuthProvider, Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"}
		http.SetCookie(w, &oauthProviderCookie)
	} else {
		oauthProviderCookie := http.Cookie{Name: "ls_oauth_provider", Value: token.OAuthProvider, Expires: time.Now().UTC().AddDate(10, 0, 0), HttpOnly: false, Secure: true, Path: "/"}
		oauthHintCookie := http.Cookie{Name: "ls_oauth_hint", Value: token.OAuthHint, Expires: time.Now().UTC().AddDate(10, 0, 0), HttpOnly: true, Secure: true, Path: "/"}
		http.SetCookie(w, &oauthProviderCookie)
		http.SetCookie(w, &oauthHintCookie)
	}

	http.SetCookie(w, &sessionExpiresAtCookie)
}

func (c *cookieManager) PassToClient(w http.ResponseWriter, r *http.Request, token *authenticator_domain.Token) error {
	cookieValue := token.User.GenerateCookieValueFromSession(authentication.PasswordSession{
		Token:     token.AuthToken,
		ExpiresAt: token.ExpiresAt.Unix(),
	})

	cookie := http.Cookie{Name: "token", Value: cookieValue, Expires: time.Unix(token.ExpiresAt.Unix(), 0), HttpOnly: true, Secure: true, Path: "/", SameSite: http.SameSiteLaxMode}

	http.SetCookie(w, &cookie)

	redirectPath := c.redirectPath
	if token.RedirectPath != "" {
		redirectPath = token.RedirectPath
	}

	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	return nil
}
