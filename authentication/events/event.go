package events

import "time"

type EventName string

const (
	EventRegistrationRequested EventName = "auth.registration.requested"
	EventRegistrationSucceeded EventName = "auth.registration.succeeded"
	EventRegistrationFailed    EventName = "auth.registration.failed"

	EventLoginRequested EventName = "auth.login.requested"
	EventLoginSucceeded EventName = "auth.login.succeeded"
	EventLoginFailed    EventName = "auth.login.failed"

	EventRosterStarted   EventName = "auth.roster.started"
	EventRosterSucceeded EventName = "auth.roster.succeeded"
	EventRosterFailed    EventName = "auth.roster.failed"

	EventAccountLinked         EventName = "auth.account_linked"
	EventEmailVerificationSent EventName = "auth.email_verification.sent"
	EventEmailVerified         EventName = "auth.email_verification.verified"

	EventSignOut EventName = "auth.sign_out"
)

type Envelope struct {
	ID         string
	Name       EventName
	OccurredAt time.Time

	RequestID string
	TraceID   string
	Source    string
	Metadata  map[string]string

	Payload any
}
