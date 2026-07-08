package authenticator_domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
)

func TestUserNotFoundErrorMessage(t *testing.T) {
	t.Parallel()
	err := &authenticator_domain.UserNotFoundError{PresentedUsername: "kenton"}
	if err.Error() != "user not found" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestUserNotFoundErrorIsSentinel(t *testing.T) {
	t.Parallel()
	err := &authenticator_domain.UserNotFoundError{PresentedUsername: "kenton"}

	if !errors.Is(err, authenticator_domain.ErrUserNotFound) {
		t.Fatal("expected errors.Is to match ErrUserNotFound")
	}
}

func TestUserNotFoundErrorAs(t *testing.T) {
	t.Parallel()
	var wrapped error = errors.Join(errors.New("wrapper"), &authenticator_domain.UserNotFoundError{PresentedUsername: "kenton"})

	var target *authenticator_domain.UserNotFoundError
	if !errors.As(wrapped, &target) {
		t.Fatal("expected errors.As to unwrap UserNotFoundError")
	}
	if target.PresentedUsername != "kenton" {
		t.Fatalf("unexpected presented username: %s", target.PresentedUsername)
	}
}

func TestInvalidPasswordErrorMessage(t *testing.T) {
	t.Parallel()
	err := &authenticator_domain.InvalidPasswordError{Username: "kenton", UserID: "id-1"}
	if err.Error() != "invalid password" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestInvalidPasswordErrorIsSentinel(t *testing.T) {
	t.Parallel()
	err := &authenticator_domain.InvalidPasswordError{Username: "kenton"}

	if !errors.Is(err, authenticator_domain.ErrInvalidPassword) {
		t.Fatal("expected errors.Is to match ErrInvalidPassword")
	}
}

func TestInvalidPasswordErrorAs(t *testing.T) {
	t.Parallel()
	var wrapped error = errors.Join(errors.New("wrapper"), &authenticator_domain.InvalidPasswordError{Username: "kenton", UserID: "id-1"})

	var target *authenticator_domain.InvalidPasswordError
	if !errors.As(wrapped, &target) {
		t.Fatal("expected errors.As to unwrap InvalidPasswordError")
	}
	if target.Username != "kenton" || target.UserID != "id-1" {
		t.Fatalf("unexpected target fields: %+v", target)
	}
}

func TestSentinelErrorsAreDistinct(t *testing.T) {
	t.Parallel()

	sentinels := []error{
		authenticator_domain.ErrUnhandleableRequest,
		authenticator_domain.ErrPasswordlessRequired,
		authenticator_domain.ErrIDNotPresent,
		authenticator_domain.ErrMethodNotSupported,
		authenticator_domain.ErrFailedToParse,
	}

	for i, a := range sentinels {
		for j, b := range sentinels {
			if i == j {
				continue
			}
			if errors.Is(a, b) {
				t.Fatalf("expected sentinel errors to be distinct: %v vs %v", a, b)
			}
		}
	}
}

func TestLinkedIdentityToMapFromMapRoundTrip(t *testing.T) {
	t.Parallel()

	original := authenticator_domain.LinkedIdentity{
		Provider: "google",
		Subject:  "subject-123",
		UserID:   "user-abc",
		Issuer:   "https://accounts.google.com",
		LinkedAt: time.Unix(1700000000, 0),
	}

	m := original.ToMap()

	if m["provider"] != "google" {
		t.Fatalf("unexpected provider in map: %v", m["provider"])
	}
	if m["subject"] != "subject-123" {
		t.Fatalf("unexpected subject in map: %v", m["subject"])
	}
	if m["user_id"] != "user-abc" {
		t.Fatalf("unexpected user_id in map: %v", m["user_id"])
	}
	if m["issuer"] != "https://accounts.google.com" {
		t.Fatalf("unexpected issuer in map: %v", m["issuer"])
	}
	if m["linked_at"] != int64(1700000000) {
		t.Fatalf("unexpected linked_at in map: %v", m["linked_at"])
	}

	roundTripped := authenticator_domain.LinkedIdentityFromMap(m)

	if roundTripped.Provider != original.Provider {
		t.Fatalf("provider mismatch after round-trip: %s", roundTripped.Provider)
	}
	if roundTripped.Subject != original.Subject {
		t.Fatalf("subject mismatch after round-trip: %s", roundTripped.Subject)
	}
	if roundTripped.UserID != original.UserID {
		t.Fatalf("user id mismatch after round-trip: %s", roundTripped.UserID)
	}
	if roundTripped.Issuer != original.Issuer {
		t.Fatalf("issuer mismatch after round-trip: %s", roundTripped.Issuer)
	}
	if roundTripped.LinkedAt.Unix() != original.LinkedAt.Unix() {
		t.Fatalf("linked at mismatch after round-trip: %v", roundTripped.LinkedAt)
	}
}
