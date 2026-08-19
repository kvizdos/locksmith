package method_password

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
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
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("unsupported content type: %w", authenticator_domain.ErrInvalidContentType)
	}

	const maxBodySize = 1 << 20 // 1 MB

	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize+1))
	if err != nil {
		return err
	}

	if len(body) > maxBodySize {
		return fmt.Errorf(
			"%w: max %d bytes",
			authenticator_domain.ErrRequestTooLarge,
			maxBodySize,
		)
	}

	// Restore the request body for anything downstream.
	r.Body = io.NopCloser(bytes.NewReader(body))

	var dto passwordRegistrationRequestDTO
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&dto); err != nil {
		return err
	}

	if prs.options.MinimumLength > 0 && len(dto.Password) < prs.options.MinimumLength {
		return fmt.Errorf(
			"password must be at least %d characters: %w",
			prs.options.MinimumLength,
			register_domain.ErrRegistrationPasswordTooShort,
		)
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
