package events

import "context"

type NoopBus struct{}

func (NoopBus) Publish(context.Context, Envelope) error { return nil }

func (NoopBus) Subscribe(EventName, Handler) Subscription { return noopSubscription{} }

type noopSubscription struct{}

func (noopSubscription) Unsubscribe() {}
