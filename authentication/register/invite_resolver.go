package register

import (
	"strings"

	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication/register/register_domain"
	"github.com/kvizdos/locksmith/database"
)

// defaultInviteResolver adapts administration/invitations onto
// register_domain.InviteResolver. It's the invite-resolution equivalent of
// the internal/method_password and internal/method_hint adapters: it keeps
// register()'s business rules decoupled from the concrete invites storage.
type defaultInviteResolver struct{}

func newDefaultInviteResolver() register_domain.InviteResolver {
	return defaultInviteResolver{}
}

func (defaultInviteResolver) ResolveByCode(db database.DatabaseAccessor, code string) (register_domain.Invite, error) {
	invite, err := invitations.GetInviteFromCode(db, code)
	if err != nil {
		return register_domain.Invite{}, err
	}
	return toDomainInvite(invite), nil
}

func (defaultInviteResolver) ResolveActiveByEmail(db database.DatabaseAccessor, email string) (register_domain.Invite, bool, error) {
	invite, found, err := invitations.GetActiveInviteByEmail(db, strings.ToLower(email))
	if err != nil || !found {
		return register_domain.Invite{}, false, err
	}
	return toDomainInvite(invite), true, nil
}

func (defaultInviteResolver) Expire(db database.DatabaseAccessor, invite register_domain.Invite) {
	invitations.ExpireByEmail(db, invite.Email)
}

func toDomainInvite(invite invitations.Invitation) register_domain.Invite {
	return register_domain.Invite{
		Email:        invite.Email,
		Role:         invite.Role,
		AttachUserID: invite.AttachUserID,
	}
}
