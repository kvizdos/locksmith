package events

type LoginSucceededPayload struct {
	UserID string

	Method       string
	Provider     string
	Passwordless bool

	// SelectBy carries provider-specific UI-surface metadata, when
	// available (see registrationhints.Hint.SelectBy).
	SelectBy string
}

type LoginFailedPayload struct {
	PresentedUsername string

	Method   string
	Provider string
	Reason   string
}

type RosterStartedPayload struct {
	Provider string

	// SelectBy carries provider-specific UI-surface metadata, when
	// available (see registrationhints.Hint.SelectBy).
	SelectBy string
}

type AccountLinkedPayload struct {
	UserID   string
	Provider string
	Issuer   string
	Subject  string
}

type AccountVerifiedPayload struct {
	UserID          string
	LoginOrRegister string
	Method          string
	Provider        string
	SelectBy        string
}

type SignOutPayload struct {
	UserID string
}
