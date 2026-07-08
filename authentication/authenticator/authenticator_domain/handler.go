package authenticator_domain

import (
	"context"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type Iconable interface {
	GetIconBytes() []byte
}

type Handler interface {
	CanHandle(r *http.Request) bool
	Name() string
	Passwordless() bool
	Session(db database.DatabaseAccessor) Session
}

type IdentityResolver interface {
	ResolveIdentity(ctx context.Context) error
}

// FederatedIdentity is implemented by sessions that authenticate via a
// third-party provider (OIDC, SAML, etc.). It carries the provider-level
// identity needed to look up or create an auth_link.
type FederatedIdentity interface {
	GetProvider() string
	GetSubject() string
	GetIssuer() string
}

// VerifiedContact is implemented by sessions that carry a contact point
// (e.g. email) that has been verified by the identity provider. It is used
// to auto-link an existing user account on first federated login.
type VerifiedContact interface {
	GetEmail() string
	EmailVerified() bool
}

type Session interface {
	LoadRequest(r *http.Request) error

	FindUserAggregate() (string, []map[string]any, error)
	GetPresentedUser() string

	IsAuthorized(user users.LocksmithUserInterface) error
}

type Beginnable interface {
	Begin(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

// Rosterable is implemented by sessions that can provide a registration hint
// when no existing user is found, triggering the auto-roster flow.
type Rosterable interface {
	RegistrationHint() *registrationhints.Hint
}

// FlowSource is implemented by sessions that can identify which UI surface
// produced the credential being presented (for example, Google Identity
// Services' "select_by" field distinguishes a rendered button click from an
// automatic One Tap/FedCM sign-in). It is optional; sessions that don't know
// how they were triggered simply don't implement it.
type FlowSource interface {
	GetSelectBy() string
}

// RedirectSource is implemented by sessions that can carry a caller-supplied
// "return to this page after login" target (the same concept as the legacy
// oauth/oidc package's "page"/"state" query parameter). It is optional; when
// implemented and non-empty, it overrides the TokenManager's default
// post-login redirect for that single login.
type RedirectSource interface {
	GetRedirectTarget() string
}
