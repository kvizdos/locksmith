package authorizer

import (
	"log/slog"
	"time"

	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
	"github.com/kvizdos/locksmith/authentication/authorizer/tokens"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/database"
)

type authorizers struct {
	db  database.DatabaseAccessor
	log *slog.Logger

	methods []authorizer_domain.Handler
	tm      tokens.TokenManager
	sp      signing.SigningPackageInterface

	// Options
	redirectPath                     string
	minimumResponseTime              time.Duration
	disableUserEnumerationProtection bool
	emailAsUsername                  bool
}

type Option func(*authorizers)

func NewAuthorizer(db database.DatabaseAccessor, opts ...Option) *authorizers {
	a := &authorizers{
		db:                  db,
		log:                 slog.Default(),
		redirectPath:        "/app",
		minimumResponseTime: 0, // defaults
	}

	for _, opt := range opts {
		opt(a)
	}

	if a.tm == nil {
		panic("token manager is required")
	}

	if len(a.methods) == 0 {
		panic("methods are required")
	}

	if a.sp == nil {
		panic("signing package is required")
	}

	return a
}

func DisableUserEnumerationProtection() Option {
	return func(a *authorizers) {
		a.disableUserEnumerationProtection = true
	}
}

func WithLogger(log *slog.Logger) Option {
	return func(a *authorizers) {
		a.log = log
	}
}

func WithRedirectPath(path string) Option {
	return func(a *authorizers) {
		a.redirectPath = path
	}
}

func WithTokenManager(tm tokens.TokenManager) Option {
	return func(a *authorizers) {
		a.tm = tm
	}
}

func WithMethods(methods ...authorizer_domain.Handler) Option {
	return func(a *authorizers) {
		a.methods = append(a.methods, methods...)
	}
}

func WithSigningPackage(sp signing.SigningPackageInterface) Option {
	return func(a *authorizers) {
		a.sp = sp
	}
}

func WithMinimumResponseTime(d time.Duration) Option {
	return func(a *authorizers) {
		a.minimumResponseTime = d
	}
}

func WithEmailAsUsername() Option {
	return func(a *authorizers) {
		a.emailAsUsername = true
	}
}
