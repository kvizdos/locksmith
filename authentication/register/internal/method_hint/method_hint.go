package method_hint

import (
	"errors"
	"net/http"
	"strings"

	"github.com/kvizdos/locksmith/authentication/register/register_domain"
	"github.com/kvizdos/locksmith/authentication/register/register_methods"
	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/database"
)

const authorizationScheme = "RegistrationHint"

func NewHintRegistration(opts ...register_methods.HintOption) register_domain.Handler {
	options := register_methods.DefaultHintOptions()
	for _, opt := range opts {
		opt(&options)
	}
	return hintRegistration{options: options}
}

type hintRegistration struct {
	options register_methods.HintOptions
}

func (hr hintRegistration) CanHandle(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.Method != http.MethodGet {
		return false
	}
	if !r.URL.Query().Has("hinted") {
		return false
	}
	if _, err := r.Cookie(registrationhints.CookieName); err == nil {
		return true
	}
	return strings.HasPrefix(r.Header.Get("Authorization"), authorizationScheme+" ")
}

func (hr hintRegistration) Name() string {
	return "hint"
}

func (hr hintRegistration) Session(db database.DatabaseAccessor) register_domain.Session {
	return &hintRegistrationSession{db: db, options: hr.options}
}

type hintRegistrationSession struct {
	db      database.DatabaseAccessor
	options register_methods.HintOptions
	request register_domain.Request
}

func (hrs *hintRegistrationSession) LoadRequest(r *http.Request) error {
	hint, err := hintFromRequest(hrs.options.Hints, r)
	if err != nil {
		return err
	}

	hrs.request = register_domain.Request{
		Username: hint.Email,
		Email:    hint.Email,
		Hint:     hint,
	}
	return nil
}

func (hrs hintRegistrationSession) RegistrationRequest() register_domain.Request {
	return hrs.request
}

func hintFromRequest(service registrationhints.Service, r *http.Request) (*registrationhints.Hint, error) {
	if r == nil {
		return nil, registrationhints.ErrMissingToken
	}

	if token, ok := authorizationToken(r.Header.Get("Authorization")); ok {
		return service.Parse(token)
	}

	hint, err := service.FromRequest(r)
	if errors.Is(err, registrationhints.ErrMissingToken) {
		return nil, registrationhints.ErrMissingToken
	}
	return hint, err
}

func authorizationToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || scheme != authorizationScheme || token == "" {
		return "", false
	}
	return token, true
}
