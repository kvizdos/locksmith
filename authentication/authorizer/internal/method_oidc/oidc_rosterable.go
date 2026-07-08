package method_oidc

import "github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"

var _ authorizer_domain.Rosterable = (*oidcValidationSession)(nil)

func (r *oidcValidationSession) RegistrationHint() *authorizer_domain.RegistrationHint {
	if !r.options.Rosterable {
		return nil
	}
	return &authorizer_domain.RegistrationHint{
		ProviderName: r.options.ProviderName,
		Email:        r.GetEmail(),
		DisplayName:  r.displayName,
		Rosterable:   true,
	}
}
