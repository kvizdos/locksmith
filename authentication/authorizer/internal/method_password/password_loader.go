package method_password

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
	"github.com/kvizdos/locksmith/authentication/authorizer/internal/sessions"
)

type loginRequestHTTP struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (pv *passwordValidatorSession) LoadRequest(r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.Join(
			fmt.Errorf("unsupported method: %s", r.Method),
			authorizer_domain.ErrMethodNotSupported,
		)
	}
	var loginRequest loginRequestHTTP
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		out := fmt.Errorf("failed to decode login request: %w", err)
		return errors.Join(out, authorizer_domain.ErrFailedToParse)
	}
	pv.presentedPassword = loginRequest.Password
	pv.BaseSession = sessions.NewBaseSession(
		sessions.WithUserID(loginRequest.Username),
	)
	return nil
}
