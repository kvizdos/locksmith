package register

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/events"
	"github.com/kvizdos/locksmith/authentication/register/register_domain"
	"github.com/kvizdos/locksmith/authentication/textvalidation"
	"github.com/kvizdos/locksmith/authentication/tokens"
	"github.com/kvizdos/locksmith/authentication/verificationcodes"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type RegisterCustomUserFunc func(users.LocksmithUser, database.DatabaseAccessor) users.LocksmithUserInterface

// registrar is the core registration orchestrator, mirroring the
// authenticator package's `authorizers` type: a private struct constructed
// via functional options and exposed to callers only through its methods
// (ServeRegisterAPI, etc.).
type registrar struct {
	db  database.DatabaseAccessor
	log *slog.Logger

	methods  []register_domain.Handler
	tm       tokens.TokenManager
	eventBus events.Bus

	defaultRoleName           string
	disablePublicRegistration bool
	configureCustomUser       RegisterCustomUserFunc
	requiresEmailVerification func(context.Context, database.DatabaseAccessor, users.LocksmithUserInterface, textvalidation.ValidationResultEvaluator) bool
	accountVerifier           verificationcodes.Verifier
	emailValidation           textvalidation.EmailValidator
	emailAsUsername           bool
	minimumLengthRequirement  int

	requestEventMetadata func(*http.Request) events.ContextMetadata
}

type Option func(*registrar)

func NewRegistrar(db database.DatabaseAccessor, opts ...Option) *registrar {
	r := &registrar{db: db, log: slog.Default(), eventBus: events.NoopBus{}}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithLogger(logger *slog.Logger) Option {
	return func(r *registrar) {
		if logger != nil {
			r.log = logger
		}
	}
}

func WithEventBus(bus events.Bus) Option {
	return func(r *registrar) {
		if bus != nil {
			r.eventBus = bus
		}
	}
}

func WithTokenManager(tm tokens.TokenManager) Option {
	return func(r *registrar) { r.tm = tm }
}

func WithMethods(methods ...register_domain.Handler) Option {
	return func(r *registrar) { r.methods = append(r.methods, methods...) }
}

func WithDefaultRoleName(role string) Option {
	return func(r *registrar) { r.defaultRoleName = role }
}

func WithDisablePublicRegistration(disabled bool) Option {
	return func(r *registrar) { r.disablePublicRegistration = disabled }
}

func WithConfigureCustomUser(fn RegisterCustomUserFunc) Option {
	return func(r *registrar) { r.configureCustomUser = fn }
}

func WithRequiresEmailVerification(fn func(context.Context, database.DatabaseAccessor, users.LocksmithUserInterface, textvalidation.ValidationResultEvaluator) bool) Option {
	return func(r *registrar) { r.requiresEmailVerification = fn }
}

func WithAccountVerifier(verifier verificationcodes.Verifier) Option {
	return func(r *registrar) { r.accountVerifier = verifier }
}

func WithEmailValidation(validator textvalidation.EmailValidator) Option {
	return func(r *registrar) { r.emailValidation = validator }
}

func WithEmailAsUsername(enabled bool) Option {
	return func(r *registrar) { r.emailAsUsername = enabled }
}

func WithMinimumLengthRequirement(length int) Option {
	return func(r *registrar) { r.minimumLengthRequirement = length }
}

func WithRequestEventMetadata(fn func(*http.Request) events.ContextMetadata) Option {
	return func(r *registrar) { r.requestEventMetadata = fn }
}
