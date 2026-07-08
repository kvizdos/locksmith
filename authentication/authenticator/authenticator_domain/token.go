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

	// Optional..
	User *users.LocksmithUser
}
