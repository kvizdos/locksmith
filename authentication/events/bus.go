package events

import "context"

type Handler func(context.Context, Envelope) error

type Subscription interface {
	Unsubscribe()
}

type Bus interface {
	Publish(ctx context.Context, event Envelope) error
	Subscribe(name EventName, handler Handler) Subscription
}
