package authorizer_domain

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnhandleableRequest  = errors.New("unhandleable request")
	ErrPasswordlessRequired = errors.New("passwordless required")

	ErrIDNotPresent = errors.New("id not present")

	ErrMethodNotSupported = errors.New("method not supported")
	ErrFailedToParse      = errors.New("failed to parse request")
)

/*
 * User Not Found Error
 */
var ErrUserNotFound = &UserNotFoundError{}

type UserNotFoundError struct {
	PresentedUsername string
	RegistrationHint  *RegistrationHint
}

type RegistrationHint struct {
	jwt.RegisteredClaims
	ProviderName string
	Email        string
	DisplayName  string
	Rosterable   bool
}

func (e *UserNotFoundError) Error() string {
	return "user not found"
}

func (e *UserNotFoundError) Is(target error) bool {
	return target == ErrUserNotFound
}

/*
 * Invalid Password Error
 */
var ErrInvalidPassword = &InvalidPasswordError{}

type InvalidPasswordError struct {
	Username string
	UserID   string
}

func (e *InvalidPasswordError) Error() string {
	return "invalid password"
}

func (e *InvalidPasswordError) Is(target error) bool {
	return target == ErrInvalidPassword
}
