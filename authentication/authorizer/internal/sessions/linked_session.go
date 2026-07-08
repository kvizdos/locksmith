package sessions

import "github.com/kvizdos/locksmith/authentication/authorizer/authorizer_domain"

var _ authorizer_domain.FederatedIdentity = LinkedSession{}
var _ authorizer_domain.VerifiedContact = LinkedSession{}

type LinkedSession struct {
	BaseSession

	subject      string
	providerName string
	email        string
	issuer       string
}

func (b LinkedSession) EmailVerified() bool {
	return b.email != ""
}

func (b LinkedSession) GetProvider() string {
	return b.providerName
}

func (b LinkedSession) GetSubject() string {
	return b.subject
}

func (b LinkedSession) GetEmail() string {
	return b.email
}

func (b LinkedSession) GetIssuer() string {
	return b.issuer
}

func (b *LinkedSession) SetIssuer(issuer string) {
	b.issuer = issuer
}

func (b *LinkedSession) SetSubject(subject string) {
	b.subject = subject
}

func (b *LinkedSession) SetEmail(email string) {
	b.email = email
	b.BaseSession.userID = email
}

type linkedSessionOptsFunc func(*LinkedSession) *LinkedSession

func WithSubject(subject string) linkedSessionOptsFunc {
	return func(b *LinkedSession) *LinkedSession {
		b.subject = subject

		return b
	}
}

func NewLinkedSession(opts ...linkedSessionOptsFunc) LinkedSession {
	b := LinkedSession{
		BaseSession: BaseSession{},
	}
	for _, opt := range opts {
		b = *opt(&b)
	}

	return b
}

func WithProvider(providerName string) linkedSessionOptsFunc {
	return func(b *LinkedSession) *LinkedSession {
		b.providerName = providerName

		return b
	}
}
