package sessions

import "github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"

type BaseSession struct {
	userID string

	baseTable string
	aggregate func(userID string) []map[string]any
}

func (b BaseSession) GetPresentedUser() string {
	return b.userID
}

type optsFunc func(*BaseSession) *BaseSession

func WithUserID(userID string) optsFunc {
	return func(b *BaseSession) *BaseSession {
		b.userID = userID

		return b
	}
}

func WithAggregate(baseTable string, aggregateMaker func(userID string) []map[string]any) optsFunc {
	return func(b *BaseSession) *BaseSession {
		b.baseTable = baseTable
		b.aggregate = aggregateMaker
		return b
	}
}

func NewBaseSession(opts ...optsFunc) BaseSession {
	b := BaseSession{}
	for _, opt := range opts {
		b = *opt(&b)
	}

	return b
}

func (b BaseSession) FindUserAggregate() (string, []map[string]any, error) {
	if b.userID == "" {
		return "", nil, authenticator_domain.ErrIDNotPresent
	}

	if b.aggregate == nil {
		return "users", []map[string]any{
			{
				"$match": map[string]any{
					"username": b.userID,
				},
			},
		}, nil
	}

	return b.baseTable, b.aggregate(b.userID), nil
}
