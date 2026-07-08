package register_domain

import "github.com/kvizdos/locksmith/authentication/registrationhints"

type Request struct {
	Username string
	Password string
	Email    string
	Code     string

	ValidationOK bool

	Hint *registrationhints.Hint
}

func (r Request) HasRequiredFields() bool {
	if r.Hint != nil {
		return r.Username != "" && r.Email != ""
	}
	return r.Username != "" && r.Password != ""
}
