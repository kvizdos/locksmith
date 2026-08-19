package register_domain

import (
	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/authentication/verificationcodes"
)

type Request struct {
	Username string
	Password string
	Email    string
	Code     string

	ValidationOK bool

	// as of now meant for email-verification-protocol token
	// optional; requires opt-in to use.
	AutoVerificationPayload verificationcodes.AutoVerificationPayload

	Hint *registrationhints.Hint
}

func (r Request) HasRequiredFields() bool {
	if r.Hint != nil {
		return r.Username != "" && r.Email != ""
	}
	return r.Username != "" && r.Password != ""
}
