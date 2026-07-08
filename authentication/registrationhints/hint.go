package registrationhints

import "github.com/golang-jwt/jwt/v5"

const (
	CookieName = "registration_hint"
	Audience   = "registration"
)

type Hint struct {
	jwt.RegisteredClaims

	ProviderName string `json:"provider_name"`
	Email        string `json:"email"`
	DisplayName  string `json:"display_name,omitempty"`

	// Required for creating an auth_links row during hinted registration.
	Issuer  string `json:"issuer,omitempty"`
	Subject string `json:"subject,omitempty"`

	Rosterable bool `json:"rosterable"`

	// SelectBy carries provider-specific UI-surface metadata, when available
	// (for example, Google Identity Services' "select_by" value, which
	// distinguishes a rendered button click from an automatic One
	// Tap/FedCM sign-in). It is informational only and is not used for any
	// security decision.
	SelectBy string `json:"select_by,omitempty"`
}
