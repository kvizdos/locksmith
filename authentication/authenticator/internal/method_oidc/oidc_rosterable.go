package method_oidc

import "github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"

var _ authenticator_domain.Rosterable = (*oidcValidationSession)(nil)

func (r *oidcValidationSession) RegistrationHint() *authenticator_domain.RegistrationHint {
	if !r.options.Rosterable {
		return nil
	}
	return &authenticator_domain.RegistrationHint{
		ProviderName: r.options.ProviderName,
		Email:        r.GetEmail(),
		DisplayName:  r.displayName,
		Rosterable:   true,
	}
}
