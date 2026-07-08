package method_oidc

import "github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"

var _ authenticator_domain.FlowSource = (*oidcValidationSession)(nil)

// GetSelectBy exposes Google Identity Services' "select_by" value (when the
// credential flow supplied one), letting consumers distinguish a rendered
// button click from an automatic One Tap/FedCM sign-in.
func (r *oidcValidationSession) GetSelectBy() string {
	return r.selectBy
}
