package authenticator_domain

import (
	"time"

	"github.com/kvizdos/locksmith/users"
)

type Token struct {
	AuthToken string
	ExpiresAt time.Time

	HandlerName string

	// For OAuth auto-login
	OAuthProvider string
	OAuthHint     string

	// RedirectPath, if set, overrides the TokenManager's default post-login
	// redirect target (see RedirectSource) for this specific login.
	RedirectPath string

	// Optional..
	User *users.LocksmithUser
}
