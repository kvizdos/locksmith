package method_oidc

import (
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/registrationhints"
)

var _ authenticator_domain.Rosterable = (*oidcValidationSession)(nil)

func (r *oidcValidationSession) RegistrationHint() *registrationhints.Hint {
	if !r.options.Rosterable {
		return nil
	}
	return &registrationhints.Hint{
		ProviderName: r.options.ProviderName,
		Email:        r.GetEmail(),
		DisplayName:  r.displayName,
		Issuer:       r.GetIssuer(),
		Subject:      r.GetSubject(),
		Rosterable:   true,
		SelectBy:     r.selectBy,
	}
}
