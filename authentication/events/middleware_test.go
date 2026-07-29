package events

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithMiddlewareAppliesToEveryPublish(t *testing.T) {
	t.Parallel()

	bus := WithMiddleware(NewMemoryBus(), func(ctx context.Context, event Envelope) Envelope {
		if event.Metadata == nil {
			event.Metadata = map[string]string{}
		}
		event.Metadata["global"] = "always"
		return event
	})

	var got []Envelope
	bus.Subscribe(EventLoginSucceeded, func(ctx context.Context, event Envelope) error {
		got = append(got, event)
		return nil
	})
	bus.Subscribe(EventSignOut, func(ctx context.Context, event Envelope) error {
		got = append(got, event)
		return nil
	})

	if err := bus.Publish(context.Background(), Envelope{Name: EventLoginSucceeded}); err != nil {
		t.Fatalf("publish login: %v", err)
	}
	if err := bus.Publish(context.Background(), Envelope{Name: EventSignOut}); err != nil {
		t.Fatalf("publish sign out: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 events delivered, got %d", len(got))
	}
	for _, event := range got {
		if event.Metadata["global"] != "always" {
			t.Fatalf("event %q missing middleware metadata: %+v", event.Name, event)
		}
	}
}

func TestWithMiddlewareCanReadRequestFromContext(t *testing.T) {
	t.Parallel()

	bus := WithMiddleware(NewMemoryBus(), func(ctx context.Context, event Envelope) Envelope {
		if req, ok := RequestFromContext(ctx); ok {
			if event.Metadata == nil {
				event.Metadata = map[string]string{}
			}
			event.Metadata["tenant_id"] = req.Header.Get("X-Tenant-Id")
		}
		return event
	})

	var got Envelope
	bus.Subscribe(EventLoginSucceeded, func(ctx context.Context, event Envelope) error {
		got = event
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Header.Set("X-Tenant-Id", "acme")
	ctx := WithRequest(context.Background(), req)

	if err := bus.Publish(ctx, Envelope{Name: EventLoginSucceeded}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	if got.Metadata["tenant_id"] != "acme" {
		t.Fatalf("tenant_id = %q, want %q", got.Metadata["tenant_id"], "acme")
	}
}

func TestChainRunsMiddlewaresInOrder(t *testing.T) {
	t.Parallel()

	var order []string
	mw1 := func(ctx context.Context, e Envelope) Envelope {
		order = append(order, "mw1")
		return e
	}
	mw2 := func(ctx context.Context, e Envelope) Envelope {
		order = append(order, "mw2")
		return e
	}

	chained := Chain(mw1, nil, mw2)
	if chained == nil {
		t.Fatal("expected non-nil chain")
	}
	chained(context.Background(), Envelope{})

	if len(order) != 2 || order[0] != "mw1" || order[1] != "mw2" {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestChainReturnsNilWhenEmpty(t *testing.T) {
	t.Parallel()

	if Chain() != nil {
		t.Fatal("expected nil chain for no middlewares")
	}
	if Chain(nil, nil) != nil {
		t.Fatal("expected nil chain for all-nil middlewares")
	}
}

func TestWithMiddlewareNoMiddlewaresReturnsBusUnwrapped(t *testing.T) {
	t.Parallel()

	bus := NewMemoryBus()
	wrapped := WithMiddleware(bus)
	if wrapped != Bus(bus) {
		t.Fatal("expected WithMiddleware to return the original bus when no middlewares are given")
	}
}

func TestWithMiddlewareNilBusWrapsNoop(t *testing.T) {
	t.Parallel()

	called := false
	bus := WithMiddleware(nil, func(ctx context.Context, e Envelope) Envelope {
		called = true
		return e
	})

	if err := bus.Publish(context.Background(), Envelope{Name: EventLoginSucceeded}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !called {
		t.Fatal("expected middleware to run even when wrapping a nil bus")
	}
}
