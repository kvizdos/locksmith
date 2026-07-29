package authenticator

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/authentication/tokens"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/users"
)

// ── TestMain ──────────────────────────────────────────────────────────────────

func init() {
	roles.AVAILABLE_ROLES = map[string][]string{
		"user":  {"view.self"},
		"admin": {"view.self", "view.admin"},
	}
}

// ── mockTokenManager ─────────────────────────────────────────────────────────

type mockTokenManager struct {
	redirectPath string
}

func (m *mockTokenManager) Read(r *http.Request) (*authenticator_domain.Token, error) {
	return nil, nil
}

func (m *mockTokenManager) CreateAuthToken(user users.LocksmithUserInterface) (*authenticator_domain.Token, error) {
	lu := user.(users.LocksmithUser)
	return &authenticator_domain.Token{
		AuthToken: "test-token",
		ExpiresAt: time.Now().Add(time.Hour),
		User:      &lu,
	}, nil
}

func (m *mockTokenManager) PassToClient(w http.ResponseWriter, r *http.Request, token *authenticator_domain.Token) error {
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token.AuthToken,
	})
	http.Redirect(w, r, m.redirectPath, http.StatusSeeOther)
	return nil
}

var _ tokens.TokenManager = (*mockTokenManager)(nil)

// ── helpers ───────────────────────────────────────────────────────────────────

func newTestDB(tables map[string]map[string]any) database.TestDatabase {
	if tables == nil {
		tables = map[string]map[string]any{}
	}
	return database.TestDatabase{Tables: tables}
}

func newTestAuthorizer(db database.DatabaseAccessor, opts ...Option) *authorizers {
	sp, _ := signing.CreateSigningPackage()
	base := []Option{
		WithTokenManager(&mockTokenManager{redirectPath: "/app"}),
		WithMethods(AllowMethodPassword()),
		WithSigningPackage(&sp),
	}
	return NewAuthorizer(db, append(base, opts...)...)
}

func compiledPassword(plain string) authentication.PasswordInfo {
	p, _ := authentication.CompileLocksmithPassword(plain)
	return p
}

func makeUser(id, username, email string, password authentication.PasswordInfo, role string) map[string]any {
	return map[string]any{
		"id":                     id,
		"username":               username,
		"email":                  email,
		"password":               password,
		"sessions":               []any{},
		"role":                   role,
		"emailVerified":          false,
		"needsEmailVerification": false,
		"oauth":                  "",
	}
}

func postLoginReq(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Body = http.NoBody
	if body != "" {
		req = httptest.NewRequest(http.MethodPost, "/api/login", stringReader(body))
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func stringReader(s string) *singleReader {
	return &singleReader{data: s, pos: 0}
}

type singleReader struct {
	data string
	pos  int
}

func (r *singleReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, nil
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *singleReader) Close() error { return nil }

// ── mockFederatedHandler / mockFederatedSession ───────────────────────────────

// mockFederatedSession is a configurable Session that also implements
// FederatedIdentity, VerifiedContact, and optionally Rosterable.
type mockFederatedSession struct {
	provider      string
	subject       string
	issuer        string
	email         string
	emailVerified bool
	rosterable    bool
	displayName   string
}

func (s *mockFederatedSession) LoadRequest(r *http.Request) error                 { return nil }
func (s *mockFederatedSession) GetPresentedUser() string                          { return s.email }
func (s *mockFederatedSession) IsAuthorized(_ users.LocksmithUserInterface) error { return nil }
func (s *mockFederatedSession) FindUserAggregate() (string, []map[string]any, error) {
	return "", nil, nil
}

// FederatedIdentity
func (s *mockFederatedSession) GetProvider() string { return s.provider }
func (s *mockFederatedSession) GetSubject() string  { return s.subject }
func (s *mockFederatedSession) GetIssuer() string   { return s.issuer }

// VerifiedContact
func (s *mockFederatedSession) GetEmail() string    { return s.email }
func (s *mockFederatedSession) EmailVerified() bool { return s.emailVerified }

// Rosterable
func (s *mockFederatedSession) RegistrationHint() *registrationhints.Hint {
	if !s.rosterable {
		return nil
	}
	return &registrationhints.Hint{
		ProviderName: s.provider,
		Email:        s.email,
		DisplayName:  s.displayName,
		Issuer:       s.issuer,
		Subject:      s.subject,
		Rosterable:   true,
	}
}

// mockFederatedHandler wraps a mockFederatedSession as a Handler.
type mockFederatedHandler struct {
	session *mockFederatedSession
}

func (h mockFederatedHandler) CanHandle(r *http.Request) bool { return true }
func (h mockFederatedHandler) Name() string                   { return "mock-federated" }
func (h mockFederatedHandler) Passwordless() bool             { return true }
func (h mockFederatedHandler) Session(_ interface {
	FindOne(string, map[string]any) (any, bool)
}) authenticator_domain.Session {
	return h.session
}

// We need it to satisfy authenticator_domain.Handler which takes database.DatabaseAccessor
type mockFederatedHandlerFull struct {
	session *mockFederatedSession
}

func (h mockFederatedHandlerFull) CanHandle(r *http.Request) bool { return true }
func (h mockFederatedHandlerFull) Name() string                   { return "mock-federated" }
func (h mockFederatedHandlerFull) Passwordless() bool             { return true }
func (h mockFederatedHandlerFull) Session(db database.DatabaseAccessor) authenticator_domain.Session {
	return h.session
}
