package method_password

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/authenticator/internal/sessions"
)

type loginRequestHTTP struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (pv *passwordValidatorSession) LoadRequest(r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.Join(
			fmt.Errorf("unsupported method: %s", r.Method),
			authenticator_domain.ErrMethodNotSupported,
		)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("unsupported content type: %w", authenticator_domain.ErrInvalidContentType)
	}

	var loginRequest loginRequestHTTP

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&loginRequest); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return fmt.Errorf("%w: max %d bytes", authenticator_domain.ErrRequestTooLarge, maxBytesErr.Limit)
		}

		out := fmt.Errorf("failed to decode login request: %w", err)
		return errors.Join(out, authenticator_domain.ErrFailedToParse)
	}

	if pv.options.MinPasswordLength > 0 && len(loginRequest.Password) < pv.options.MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters: %w", pv.options.MinPasswordLength, authenticator_domain.ErrPasswordTooShort)
	}

	pv.presentedPassword = loginRequest.Password
	pv.BaseSession = sessions.NewBaseSession(
		sessions.WithUserID(loginRequest.Username),
	)
	return nil
}
