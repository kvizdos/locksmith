package register_domain

import "errors"

var (
	// Dispatch-level errors, mirroring authenticator_domain's handler
	// dispatch sentinels.
	ErrUnhandleableRequest = errors.New("unhandleable registration request")
	ErrMethodNotSupported  = errors.New("method not supported")
	ErrFailedToParse       = errors.New("failed to parse registration request")

	// Registration business-rule errors.
	ErrRegistrationRoleMissing       = errors.New("registration role name must be set")
	ErrRegistrationRoleInvalid       = errors.New("registration role name is invalid")
	ErrPublicRegistrationDisabled    = errors.New("public registration disabled")
	ErrRegistrationMissingFields     = errors.New("missing fields")
	ErrRegistrationPasswordTooShort  = errors.New("password too short")
	ErrRegistrationIllegalUsername   = errors.New("illegal username characters")
	ErrRegistrationInvalidEmail      = errors.New("invalid email")
	ErrRegistrationBadInviteCode     = errors.New("bad invite code")
	ErrRegistrationInvalidInviteCode = errors.New("invalid code")
	ErrRegistrationTaken             = errors.New("taken")
	ErrRegistrationEmailBlocked      = errors.New("email blocked")
	ErrRegistrationInvalidHint       = errors.New("invalid registration hint")
)

// ErrRegistrationConfirmEmail is the sentinel matched by ConfirmEmailError,
// mirroring the authenticator_domain.ErrUserNotFound / UserNotFoundError
// pattern.
var ErrRegistrationConfirmEmail = &ConfirmEmailError{}

// ConfirmEmailError indicates the presented email requires user confirmation
// (e.g. a suspected typo) before registration can proceed.
type ConfirmEmailError struct {
	DidYouMean string
}

func (e *ConfirmEmailError) Error() string { return "confirm email" }

func (e *ConfirmEmailError) Is(target error) bool {
	return target == ErrRegistrationConfirmEmail
}
