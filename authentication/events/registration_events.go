package events

type RegistrationRequestedPayload struct {
	Username string
	Email    string

	Method   string
	Provider string

	InviteProvided bool
	Background     bool

	// SelectBy carries provider-specific UI-surface metadata, when
	// available (see registrationhints.Hint.SelectBy).
	SelectBy string
}

type RegistrationSucceededPayload struct {
	UserID   string
	Username string
	Email    string

	Method   string
	Provider string

	InviteUsed                bool
	RequiresEmailVerification bool
	Background                bool
	AutoLoginIssued           bool

	// SelectBy carries provider-specific UI-surface metadata, when
	// available (see registrationhints.Hint.SelectBy).
	SelectBy string
}

type RegistrationFailedPayload struct {
	Username string
	Email    string

	Method   string
	Provider string

	InviteProvided bool
	Background     bool

	Reason string

	// SelectBy carries provider-specific UI-surface metadata, when
	// available (see registrationhints.Hint.SelectBy).
	SelectBy string
}

type EmailVerificationSentPayload struct {
	UserID   string
	Username string
	Email    string

	Method string
	Target string
}
