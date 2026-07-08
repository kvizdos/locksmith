package register_domain

import "github.com/kvizdos/locksmith/users"

type Result struct {
	User users.LocksmithUserInterface

	Method   string
	Provider string

	// Background indicates the registration happened without direct user
	// interaction (e.g. rostered from a signed registration hint).
	Background bool

	CreatedAuthLink           bool
	InviteUsed                bool
	RequiresEmailVerification bool
}
