package saml_entities

type SAMLProvider struct {
	Nickname string
	// Identity
	EntityID string // md.EntityID

	// Where we are allowed to send assertions
	ACSURL  string // ONE chosen AssertionConsumerService.Location
	Binding string // should be HTTP-POST

	// Verification of inbound requests (optional)
	SigningCertPEM *string // SP signing cert, if present

	// Assertion expectations
	NameIDFormat         string // e.g. emailAddress or unspecified
	WantAssertionsSigned bool   // md.SPSSODescriptor.WantAssertionsSigned

	// Operational flags
	Enabled bool
}
