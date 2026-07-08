package authenticator

import (
	"context"
	"time"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/events"
)

func (a authorizers) LinkAccount(
	ctx context.Context,
	userID string,
	provider string,
	issuer string,
	providerSubjectID string) error {
	_, err := a.db.InsertOne("auth_links", authenticator_domain.LinkedIdentity{
		Provider: provider,
		Issuer:   issuer,
		Subject:  providerSubjectID,
		UserID:   userID,
		LinkedAt: time.Now().UTC(),
	}.ToMap())
	if err != nil {
		return err
	}

	a.publishAuthEvent(ctx, events.EventAccountLinked, events.AccountLinkedPayload{
		UserID:   userID,
		Provider: provider,
		Issuer:   issuer,
		Subject:  providerSubjectID,
	})
	return nil
}
