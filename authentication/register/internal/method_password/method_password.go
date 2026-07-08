package method_password

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/register/register_domain"
	"github.com/kvizdos/locksmith/authentication/register/register_methods"
	"github.com/kvizdos/locksmith/database"
)

func NewPasswordRegistration(opts ...register_methods.PasswordOption) register_domain.Handler {
	options := register_methods.DefaultPasswordOptions()
	for _, opt := range opts {
		opt(&options)
	}
	return passwordRegistration{options: options}
}

type passwordRegistration struct {
	options register_methods.PasswordOptions
}

func (pr passwordRegistration) CanHandle(r *http.Request) bool {
	return r != nil && r.Method == http.MethodPost
}

func (pr passwordRegistration) Name() string {
	return "password"
}

func (pr passwordRegistration) Session(db database.DatabaseAccessor) register_domain.Session {
	return &passwordRegistrationSession{db: db, options: pr.options}
}

type passwordRegistrationSession struct {
	db      database.DatabaseAccessor
	options register_methods.PasswordOptions
	request register_domain.Request
}

type passwordRegistrationRequestDTO struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Code         string `json:"code"`
	ValidationOK bool   `json:"validationok,omitempty"`
}

func (prs *passwordRegistrationSession) LoadRequest(r *http.Request) error {
	var dto passwordRegistrationRequestDTO
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dto); err != nil {
		return err
	}

	prs.request = register_domain.Request{
		Username:     dto.Username,
		Password:     dto.Password,
		Email:        dto.Email,
		Code:         dto.Code,
		ValidationOK: dto.ValidationOK,
	}
	return nil
}

func (prs passwordRegistrationSession) RegistrationRequest() register_domain.Request {
	return prs.request
}
