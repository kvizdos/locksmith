package sessions

import (
	"errors"
	"testing"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
)

func TestBaseSessionGetPresentedUser(t *testing.T) {
	t.Parallel()

	b := NewBaseSession(WithUserID("kenton"))

	if got := b.GetPresentedUser(); got != "kenton" {
		t.Fatalf("expected kenton, got %s", got)
	}
}

func TestBaseSessionFindUserAggregateNoID(t *testing.T) {
	t.Parallel()

	b := NewBaseSession()

	_, _, err := b.FindUserAggregate()
	if !errors.Is(err, authenticator_domain.ErrIDNotPresent) {
		t.Fatalf("expected ErrIDNotPresent, got %v", err)
	}
}

func TestBaseSessionFindUserAggregateDefault(t *testing.T) {
	t.Parallel()

	b := NewBaseSession(WithUserID("kenton"))

	table, pipeline, err := b.FindUserAggregate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if table != "users" {
		t.Fatalf("expected default table 'users', got %s", table)
	}
	if len(pipeline) != 1 {
		t.Fatalf("expected single stage pipeline, got %d", len(pipeline))
	}

	match, ok := pipeline[0]["$match"].(map[string]any)
	if !ok {
		t.Fatalf("expected $match stage, got %+v", pipeline[0])
	}
	if match["username"] != "kenton" {
		t.Fatalf("expected username match on 'kenton', got %v", match["username"])
	}
}

func TestBaseSessionFindUserAggregateCustom(t *testing.T) {
	t.Parallel()

	called := false
	aggregateMaker := func(userID string) []map[string]any {
		called = true
		return []map[string]any{
			{"$match": map[string]any{"custom": userID}},
		}
	}

	b := NewBaseSession(WithUserID("kenton"), WithAggregate("custom_table", aggregateMaker))

	table, pipeline, err := b.FindUserAggregate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected aggregate maker to be called")
	}
	if table != "custom_table" {
		t.Fatalf("expected custom_table, got %s", table)
	}
	if len(pipeline) != 1 {
		t.Fatalf("expected single stage pipeline, got %d", len(pipeline))
	}
	match, ok := pipeline[0]["$match"].(map[string]any)
	if !ok {
		t.Fatalf("expected $match stage, got %+v", pipeline[0])
	}
	if match["custom"] != "kenton" {
		t.Fatalf("expected custom match on 'kenton', got %v", match["custom"])
	}
}

func TestLinkedSessionEmailVerified(t *testing.T) {
	t.Parallel()

	unset := NewLinkedSession()
	if unset.EmailVerified() {
		t.Fatal("expected EmailVerified to be false when email is unset")
	}

	set := NewLinkedSession()
	set.SetEmail("kenton@example.com")
	if !set.EmailVerified() {
		t.Fatal("expected EmailVerified to be true when email is set")
	}
}

func TestLinkedSessionAccessors(t *testing.T) {
	t.Parallel()

	l := NewLinkedSession(WithProvider("google"), WithSubject("subject-1"))

	if l.GetProvider() != "google" {
		t.Fatalf("unexpected provider: %s", l.GetProvider())
	}
	if l.GetSubject() != "subject-1" {
		t.Fatalf("unexpected subject: %s", l.GetSubject())
	}

	l.SetIssuer("https://issuer.example.com")
	if l.GetIssuer() != "https://issuer.example.com" {
		t.Fatalf("unexpected issuer: %s", l.GetIssuer())
	}

	l.SetSubject("subject-2")
	if l.GetSubject() != "subject-2" {
		t.Fatalf("unexpected subject after SetSubject: %s", l.GetSubject())
	}
}

func TestLinkedSessionSetEmailUpdatesUserID(t *testing.T) {
	t.Parallel()

	l := NewLinkedSession()
	l.SetEmail("kenton@example.com")

	if l.GetEmail() != "kenton@example.com" {
		t.Fatalf("unexpected email: %s", l.GetEmail())
	}
	if l.GetPresentedUser() != "kenton@example.com" {
		t.Fatalf("expected GetPresentedUser to reflect email, got %s", l.GetPresentedUser())
	}
}

func TestLinkedSessionImplementsInterfaces(t *testing.T) {
	t.Parallel()

	var _ authenticator_domain.FederatedIdentity = LinkedSession{}
	var _ authenticator_domain.VerifiedContact = LinkedSession{}
}
