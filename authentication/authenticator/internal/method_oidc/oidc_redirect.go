package method_oidc

import (
	"net/url"
	"strings"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
)

var _ authenticator_domain.RedirectSource = (*oidcValidationSession)(nil)

// GetRedirectTarget exposes the caller-supplied "return to this page after
// login" target, when one was provided and passed validation.
func (r *oidcValidationSession) GetRedirectTarget() string {
	return r.redirectTarget
}

// sanitizeRedirectPath validates that raw is safe to redirect to: a
// same-site, relative path. This mirrors the validation the legacy
// authentication/oauth/oidc package applies to its "page"/"state"
// parameter, preventing open-redirect attacks via an absolute or
// protocol-relative URL. Returns "" if raw is empty or unsafe.
func sanitizeRedirectPath(raw string) string {
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.IsAbs() || parsed.Host != "" || !strings.HasPrefix(parsed.Path, "/") {
		return ""
	}
	return parsed.String()
}
