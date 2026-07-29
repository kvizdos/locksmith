package events

import "context"

// Middleware runs against every Envelope immediately before it reaches a
// Bus's subscribers, no matter which package published it (login,
// registration, sign-out, rostering, account linking, etc).
//
// Register middleware once, when the Bus is constructed (see
// WithMiddleware), instead of wiring a per-caller callback into each
// orchestrator individually. A single wrapped Bus, handed to authenticator.WithEventBus,
// register.WithEventBus, sign_out_http.SignOutHTTP{EventBus: ...}, and
// routes.LocksmithRoutesOptions.Bus, enriches every event the same way.
//
// Middleware commonly pulls request-scoped data out of ctx (see
// WithRequest/RequestFromContext) and copies it onto the Envelope's
// Metadata, RequestID, TraceID, or even Payload.
type Middleware func(ctx context.Context, event Envelope) Envelope

// Chain composes middlewares in order, feeding each one's output into the
// next. Nil middlewares are skipped. Chain returns nil if there is nothing
// to run.
func Chain(middlewares ...Middleware) Middleware {
	filtered := make([]Middleware, 0, len(middlewares))
	for _, mw := range middlewares {
		if mw != nil {
			filtered = append(filtered, mw)
		}
	}
	if len(filtered) == 0 {
		return nil
	}

	return func(ctx context.Context, event Envelope) Envelope {
		for _, mw := range filtered {
			event = mw(ctx, event)
		}
		return event
	}
}

// WithMiddleware wraps bus so every Publish call is passed through
// middlewares (in order) before it reaches bus's subscribers. Build the
// wrapped bus once at startup and hand that single instance to every
// consumer, so the middleware chain applies everywhere the bus is used:
//
//	bus := events.WithMiddleware(events.NewMemoryBus(),
//		func(ctx context.Context, e events.Envelope) events.Envelope {
//			if req, ok := events.RequestFromContext(ctx); ok {
//				if e.Metadata == nil {
//					e.Metadata = map[string]string{}
//				}
//				e.Metadata["tenant_id"] = req.Header.Get("X-Tenant-Id")
//			}
//			return e
//		},
//	)
//
//	authorizer := authenticator.NewAuthorizer(db, authenticator.WithEventBus(bus), ...)
//	registrar := register.NewRegistrar(db, register.WithEventBus(bus), ...)
//	routes.InitializeLocksmithRoutes(mux, db, routes.LocksmithRoutesOptions{Bus: bus, ...})
//
// If bus is nil, WithMiddleware wraps NoopBus{}. If no middlewares are
// given, bus is returned unwrapped.
func WithMiddleware(bus Bus, middlewares ...Middleware) Bus {
	if bus == nil {
		bus = NoopBus{}
	}

	chain := Chain(middlewares...)
	if chain == nil {
		return bus
	}

	return &middlewareBus{next: bus, mw: chain}
}

type middlewareBus struct {
	next Bus
	mw   Middleware
}

func (b *middlewareBus) Publish(ctx context.Context, event Envelope) error {
	event = b.mw(ctx, event)
	return b.next.Publish(ctx, event)
}

func (b *middlewareBus) Subscribe(name EventName, handler Handler) Subscription {
	return b.next.Subscribe(name, handler)
}
