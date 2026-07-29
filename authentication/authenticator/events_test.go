package authenticator

import (
	"context"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/kvizdos/locksmith/authentication/events"
)

type recordingEventBus struct {
	mu        sync.Mutex
	published []events.Envelope
}

func (b *recordingEventBus) Publish(ctx context.Context, event events.Envelope) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.published = append(b.published, event)
	return nil
}

func (b *recordingEventBus) Subscribe(events.EventName, events.Handler) events.Subscription {
	return recordingSubscription{}
}

type recordingSubscription struct{}

func (recordingSubscription) Unsubscribe() {}

func (b *recordingEventBus) eventsByName(name events.EventName) []events.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []events.Envelope
	for _, e := range b.published {
		if e.Name == name {
			out = append(out, e)
		}
	}
	return out
}

func (b *recordingEventBus) single(t *testing.T, name events.EventName) events.Envelope {
	t.Helper()
	got := b.eventsByName(name)
	if len(got) != 1 {
		t.Fatalf("events named %q = %d, want 1: %+v", name, len(got), got)
	}
	return got[0]
}

func TestServeLoginAPIPublishesLoginSucceeded(t *testing.T) {
	t.Parallel()

	user := makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user")
	db := newTestDB(map[string]map[string]any{"users": {"u1": user}})
	bus := &recordingEventBus{}
	a := newTestAuthorizer(db, WithEventBus(bus))

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"hunter2"}`))

	if rr.Code != 303 && rr.Code != 200 {
		t.Fatalf("unexpected status: %d body=%s", rr.Code, rr.Body.String())
	}

	envelope := bus.single(t, events.EventLoginSucceeded)
	payload, ok := envelope.Payload.(events.LoginSucceededPayload)
	if !ok {
		t.Fatalf("payload type = %T, want events.LoginSucceededPayload", envelope.Payload)
	}
	if payload.UserID != "u1" {
		t.Fatalf("UserID = %q, want %q", payload.UserID, "u1")
	}
	if len(bus.eventsByName(events.EventLoginFailed)) != 0 {
		t.Fatal("did not expect a login failed event on success")
	}
}

func TestServeLoginAPIPublishesLoginFailedOnInvalidPassword(t *testing.T) {
	t.Parallel()

	user := makeUser("u1", "kenton", "k@example.com", compiledPassword("hunter2"), "user")
	db := newTestDB(map[string]map[string]any{"users": {"u1": user}})
	bus := &recordingEventBus{}
	a := newTestAuthorizer(db, WithEventBus(bus))

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"kenton","password":"wrong-password"}`))

	envelope := bus.single(t, events.EventLoginFailed)
	payload, ok := envelope.Payload.(events.LoginFailedPayload)
	if !ok {
		t.Fatalf("payload type = %T, want events.LoginFailedPayload", envelope.Payload)
	}
	if payload.Reason != "invalid_password" {
		t.Fatalf("Reason = %q, want %q", payload.Reason, "invalid_password")
	}
	if strings.Contains(payload.PresentedUsername+payload.Reason, "wrong-password") {
		t.Fatal("failed login event must not leak the presented password")
	}
	if len(bus.eventsByName(events.EventLoginSucceeded)) != 0 {
		t.Fatal("did not expect a login succeeded event on failure")
	}
}

func TestServeLoginAPIPublishesLoginFailedOnUnknownUser(t *testing.T) {
	t.Parallel()

	db := newTestDB(map[string]map[string]any{"users": {}})
	bus := &recordingEventBus{}
	a := newTestAuthorizer(db, WithEventBus(bus))

	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, postLoginReq(`{"username":"ghost","password":"whatever"}`))

	envelope := bus.single(t, events.EventLoginFailed)
	payload := envelope.Payload.(events.LoginFailedPayload)
	if payload.Reason != "user_not_found" {
		t.Fatalf("Reason = %q, want %q", payload.Reason, "user_not_found")
	}
}

func TestLinkAccountPublishesAccountLinked(t *testing.T) {
	t.Parallel()

	db := newTestDB(map[string]map[string]any{"auth_links": {}})
	bus := &recordingEventBus{}
	a := newTestAuthorizer(db, WithEventBus(bus))

	if err := a.LinkAccount(context.Background(), "user-abc", "google", "https://accounts.google.com", "sub-xyz"); err != nil {
		t.Fatalf("LinkAccount returned error: %v", err)
	}

	envelope := bus.single(t, events.EventAccountLinked)
	payload, ok := envelope.Payload.(events.AccountLinkedPayload)
	if !ok {
		t.Fatalf("payload type = %T, want events.AccountLinkedPayload", envelope.Payload)
	}
	if payload.UserID != "user-abc" || payload.Provider != "google" || payload.Subject != "sub-xyz" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}
