package xsrf_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/authenticator/addons/xsrf"
	"github.com/kvizdos/locksmith/authentication/signing"
)

func newManager(t *testing.T, ttl time.Duration) xsrf.Manager {
	t.Helper()
	pkg, err := signing.CreateSigningPackage()
	if err != nil {
		t.Fatalf("failed to create signing package: %v", err)
	}
	return xsrf.New(&pkg, ttl)
}

func TestGenerate_EmptySessionID(t *testing.T) {
	m := newManager(t, time.Minute)
	_, err := m.Generate("")
	if err == nil {
		t.Fatal("expected error for empty session ID")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		sessionID  string
		validateID string
		ttl        time.Duration
		wantErr    error
	}{
		{
			name:       "valid token and matching session",
			sessionID:  "session-abc",
			validateID: "session-abc",
			ttl:        time.Minute,
			wantErr:    nil,
		},
		{
			name:       "session ID mismatch",
			sessionID:  "session-abc",
			validateID: "session-xyz",
			ttl:        time.Minute,
			wantErr:    xsrf.ErrSessionMismatch,
		},
		{
			name:       "expired token",
			sessionID:  "session-abc",
			validateID: "session-abc",
			ttl:        -time.Second, // already expired
			wantErr:    xsrf.ErrTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newManager(t, tt.ttl)
			token, err := m.Generate(tt.sessionID)
			if err != nil {
				t.Fatalf("Generate() unexpected error: %v", err)
			}

			gotErr := m.Validate(token, tt.validateID)
			if tt.wantErr == nil && gotErr != nil {
				t.Errorf("Validate() unexpected error: %v", gotErr)
			}
			if tt.wantErr != nil && !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Validate() got error %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestValidate_TamperedToken(t *testing.T) {
	m := newManager(t, time.Minute)
	err := m.Validate("not.a.real.jwt", "session-abc")
	if !errors.Is(err, xsrf.ErrTokenInvalid) {
		t.Errorf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestValidate_WrongSigningKey(t *testing.T) {
	m1 := newManager(t, time.Minute)
	m2 := newManager(t, time.Minute) // different key

	token, err := m1.Generate("session-abc")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if err = m2.Validate(token, "session-abc"); !errors.Is(err, xsrf.ErrTokenInvalid) {
		t.Errorf("expected ErrTokenInvalid for wrong key, got %v", err)
	}
}
