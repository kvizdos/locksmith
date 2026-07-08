package events

import (
	"context"
	"errors"
	"testing"
)

func TestMemoryBusPublish(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		subscribe func(*testing.T, *MemoryBus) Subscription
		wantCalls int
		wantErr   bool
	}{
		{
			name:      "no subscribers succeeds",
			subscribe: func(t *testing.T, b *MemoryBus) Subscription { return nil },
		},
		{
			name: "subscriber receives event",
			subscribe: func(t *testing.T, b *MemoryBus) Subscription {
				return b.Subscribe(EventLoginSucceeded, func(ctx context.Context, event Envelope) error {
					if event.Name != EventLoginSucceeded {
						t.Fatalf("expected event name %q, got %q", EventLoginSucceeded, event.Name)
					}
					return nil
				})
			},
			wantCalls: 1,
		},
		{
			name: "handler error is returned",
			subscribe: func(t *testing.T, b *MemoryBus) Subscription {
				return b.Subscribe(EventLoginSucceeded, func(context.Context, Envelope) error {
					return errors.New("handler failed")
				})
			},
			wantErr: true,
		},
		{
			name: "unsubscribe stops delivery",
			subscribe: func(t *testing.T, b *MemoryBus) Subscription {
				sub := b.Subscribe(EventLoginSucceeded, func(context.Context, Envelope) error { return nil })
				sub.Unsubscribe()
				return sub
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bus := NewMemoryBus()
			calls := 0
			bus.Subscribe(EventLoginSucceeded, func(context.Context, Envelope) error {
				calls++
				return nil
			}).Unsubscribe()

			tt.subscribe(t, bus)
			if tt.wantCalls > 0 {
				bus.Subscribe(EventLoginSucceeded, func(context.Context, Envelope) error {
					calls++
					return nil
				})
			}

			err := bus.Publish(context.Background(), Envelope{Name: EventLoginSucceeded})
			if tt.wantErr && err == nil {
				t.Fatal("expected handler error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if calls != tt.wantCalls {
				t.Fatalf("expected %d calls, got %d", tt.wantCalls, calls)
			}
		})
	}
}

func TestMemoryBusMultipleSubscribersReceiveEvent(t *testing.T) {
	t.Parallel()

	bus := NewMemoryBus()
	calls := 0
	for i := 0; i < 2; i++ {
		bus.Subscribe(EventRegistrationSucceeded, func(context.Context, Envelope) error {
			calls++
			return nil
		})
	}

	if err := bus.Publish(context.Background(), Envelope{Name: EventRegistrationSucceeded}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}
