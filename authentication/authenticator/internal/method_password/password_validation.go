package method_password

import (
	"fmt"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"github.com/kvizdos/locksmith/authentication/authenticator/internal/sessions"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func newPasswordValidatorSession(db database.DatabaseAccessor, opts authenticator_methods.PasswordValidatorOptions) *passwordValidatorSession {
	return &passwordValidatorSession{
		db:          db,
		BaseSession: sessions.BaseSession{},
		options:     opts,
	}
}

type passwordValidatorSession struct {
	sessions.BaseSession

	// Common
	db      database.DatabaseAccessor
	options authenticator_methods.PasswordValidatorOptions

	// Password Specific Context
	presentedPassword string
}

func (pv passwordValidatorSession) IsAuthorized(user users.LocksmithUserInterface) error {
	passwordValidated, err := user.ValidatePassword(pv.presentedPassword)

	if err != nil {
		return fmt.Errorf("failed to validate password: %w", err)
	}

	if !passwordValidated {
		return authenticator_domain.ErrInvalidPassword
	}

	return nil
}
