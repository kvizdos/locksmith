package register_domain

import "github.com/kvizdos/locksmith/database"

// Invite is the registration package's normalized view of an invitation,
// decoupled from any specific invite storage/implementation.
type Invite struct {
	Email        string
	Role         string
	AttachUserID string
}

// InviteResolver resolves invitations during registration. It plays the
// same role as verificationcodes.Verifier / textvalidation.EmailValidator:
// a small, injectable seam so register()'s business rules don't couple
// directly to a concrete invites implementation.
type InviteResolver interface {
	// ResolveByCode looks up the invite for an explicitly supplied invite
	// code. Any error is treated as fatal to the registration attempt (the
	// caller handed us a code, so we owe them a concrete answer).
	ResolveByCode(db database.DatabaseAccessor, code string) (Invite, error)

	// ResolveActiveByEmail opportunistically looks for an active invite
	// matching an email when no code was supplied, e.g. for an OAuth/hinted
	// registration whose provider-verified email happens to match an
	// existing invite. found=false, err=nil means "no invite; proceed as a
	// normal, uninvited registration."
	ResolveActiveByEmail(db database.DatabaseAccessor, email string) (invite Invite, found bool, err error)

	// Expire marks the invite as consumed once registration succeeds.
	Expire(db database.DatabaseAccessor, invite Invite)
}
