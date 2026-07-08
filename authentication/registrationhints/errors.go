package registrationhints

import "errors"

var (
	ErrMissingSigner   = errors.New("registration hints: missing signer")
	ErrInvalidToken    = errors.New("registration hints: invalid token")
	ErrMissingToken    = errors.New("registration hints: missing token")
	ErrNotRosterable   = errors.New("registration hints: hint is not rosterable")
	ErrMissingEmail    = errors.New("registration hints: missing email")
	ErrMissingIdentity = errors.New("registration hints: missing issuer or subject")
)
