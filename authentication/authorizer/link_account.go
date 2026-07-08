package authorizer

import (
	"context"
	"time"

	"github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"
)

func (a authorizers) LinkAccount(
	ctx context.Context,
	userID string,
	provider string,
	issuer string,
	providerSubjectID string) error {
	_, err := a.db.InsertOne("auth_links", authorizer_domain.LinkedIdentity{
		Provider: provider,
		Issuer:   issuer,
		Subject:  providerSubjectID,
		UserID:   userID,
		LinkedAt: time.Now().UTC(),
	}.ToMap())
	return err
}
