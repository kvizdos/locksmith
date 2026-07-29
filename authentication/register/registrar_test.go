package register

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/events"
	"github.com/kvizdos/locksmith/authentication/register/register_domain"
	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/authentication/verificationcodes"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func newRosterableHint(overrides ...func(*registrationhints.Hint)) registrationhints.Hint {
	hint := registrationhints.Hint{
		ProviderName: "google",
		Email:        "rostered@example.com",
		Issuer:       "https://issuer.example.com",
		Subject:      "sub-123",
		Rosterable:   true,
	}
	for _, override := range overrides {
		override(&hint)
	}
	return hint
}

func newTestHintService(t *testing.T) registrationhints.Service {
	t.Helper()
	signer, err := signing.CreateSigningPackage()
	if err != nil {
		t.Fatalf("create signer: %v", err)
	}
	return registrationhints.Service{Signer: signer}
}

func newRegistrarTestDB(tables map[string]map[string]any) database.TestDatabase {
	if tables == nil {
		tables = map[string]map[string]any{}
	}
	if _, ok := tables["users"]; !ok {
		tables["users"] = map[string]any{}
	}
	return database.TestDatabase{Tables: tables}
}

func TestRegistrarRegisterCreatesLowercaseUserWithHashedPassword(t *testing.T) {
	db := newRegistrarTestDB(nil)
	r := NewRegistrar(db, WithDefaultRoleName("admin"))

	result, err := r.register(context.Background(), register_domain.Request{
		Username: "Kenton",
		Password: "password123",
		Email:    "Email@Example.com",
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if result.User.GetUsername() != "kenton" {
		t.Fatalf("username = %q, want %q", result.User.GetUsername(), "kenton")
	}
	if result.User.GetEmail() != "email@example.com" {
		t.Fatalf("email = %q, want %q", result.User.GetEmail(), "email@example.com")
	}

	raw, found := db.FindOne("users", map[string]any{"username": "kenton"})
	if !found {
		t.Fatal("new user was not inserted")
	}
	var lsu users.LocksmithUserInterface = users.LocksmithUser{}
	lsu.ReadFromMap(&lsu, raw.(map[string]any))
	passwordInfo := lsu.GetPasswordInfo()
	if passwordInfo.Password == "password123" {
		t.Fatal("password was stored in plaintext")
	}
	if len(passwordInfo.Salt) != 32 {
		t.Fatalf("salt length = %d, want 32", len(passwordInfo.Salt))
	}
}

func TestRegistrarRegisterRejectsDuplicateUsernameOrEmail(t *testing.T) {
	db := newRegistrarTestDB(map[string]map[string]any{
		"users": {
			"u1": map[string]any{
				"id":       "u1",
				"username": "kenton",
				"email":    "email@example.com",
			},
		},
	})
	r := NewRegistrar(db, WithDefaultRoleName("admin"))

	_, err := r.register(context.Background(), register_domain.Request{
		Username: "Kenton",
		Password: "password123",
		Email:    "other@example.com",
	})
	if !errors.Is(err, register_domain.ErrRegistrationTaken) {
		t.Fatalf("duplicate username error = %v, want ErrRegistrationTaken", err)
	}

	_, err = r.register(context.Background(), register_domain.Request{
		Username: "other",
		Password: "password123",
		Email:    "Email@Example.com",
	})
	if !errors.Is(err, register_domain.ErrRegistrationTaken) {
		t.Fatalf("duplicate email error = %v, want ErrRegistrationTaken", err)
	}
}

func TestRegistrarRegisterUsesInviteRoleIDAndExpiresInvite(t *testing.T) {
	code := "jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M"
	hasher := sha256.New()
	hasher.Write([]byte(code))
	hashedCode := fmt.Sprintf("%x", hasher.Sum(nil))

	db := newRegistrarTestDB(map[string]map[string]any{
		"users": {},
		"invites": {
			"invite1": map[string]any{
				"code":    hashedCode,
				"email":   "bob@bob.com",
				"role":    "admin",
				"inviter": "inviter-id",
				"sentAt":  time.Now().Unix(),
				"userid":  "invited-user-id",
			},
		},
	})
	r := NewRegistrar(db, WithDefaultRoleName("admin"), WithDisablePublicRegistration(true))

	result, err := r.register(context.Background(), register_domain.Request{
		Username: "bob",
		Password: "password123",
		Email:    "bob@bob.com",
		Code:     code,
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if !result.InviteUsed {
		t.Fatal("result.InviteUsed = false, want true")
	}
	if result.User.GetID() != "invited-user-id" {
		t.Fatalf("user id = %q, want invited-user-id", result.User.GetID())
	}
	if _, found := db.FindOne("invites", map[string]any{"code": hashedCode}); found {
		t.Fatal("invite was not expired")
	}
}

func TestRegistrarPasswordRegistrationDefersPendingInviteUntilEmailVerified(t *testing.T) {
	sender := &recordingVerificationSender{}
	db := newRegistrarTestDB(map[string]map[string]any{
		"users": {},
		"invites": {
			"invite1": map[string]any{
				"code":    "unused-hashed-code",
				"email":   "bob@bob.com",
				"role":    "admin",
				"inviter": "inviter-id",
				"sentAt":  time.Now().Unix(),
				"userid":  "invited-user-id",
			},
		},
	})
	r := NewRegistrar(db, WithDefaultRoleName("admin"), WithAccountVerifier(verificationcodes.NewVerifier(db, sender)))

	result, err := r.register(context.Background(), register_domain.Request{
		Username: "bob",
		Password: "password123",
		Email:    "bob@bob.com",
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if result.InviteUsed {
		t.Fatal("result.InviteUsed = true, want false (must not be applied before verification)")
	}
	if result.User.GetID() == "invited-user-id" {
		t.Fatal("user was assigned the invite's pinned id before verifying their email")
	}
	if !result.RequiresEmailVerification {
		t.Fatal("RequiresEmailVerification = false, want true because the email matches a pending invite")
	}
	if sender.calls != 1 {
		t.Fatalf("verification send calls = %d, want 1", sender.calls)
	}
	if _, found := db.FindOne("invites", map[string]any{"email": "bob@bob.com"}); !found {
		t.Fatal("invite should still be pending, not yet expired")
	}
}

func TestRegistrarPasswordRegistrationIgnoresPendingInviteWithoutAccountVerifier(t *testing.T) {
	db := newRegistrarTestDB(map[string]map[string]any{
		"users": {},
		"invites": {
			"invite1": map[string]any{
				"code":    "unused-hashed-code",
				"email":   "bob@bob.com",
				"role":    "admin",
				"inviter": "inviter-id",
				"sentAt":  time.Now().Unix(),
				"userid":  "invited-user-id",
			},
		},
	})
	// No WithAccountVerifier: there's no way for this user to ever prove
	// control of their email, so the pending-invite check must be a no-op.
	r := NewRegistrar(db, WithDefaultRoleName("admin"))

	result, err := r.register(context.Background(), register_domain.Request{
		Username: "bob",
		Password: "password123",
		Email:    "bob@bob.com",
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if result.RequiresEmailVerification {
		t.Fatal("RequiresEmailVerification = true, want false without a configured account verifier")
	}
	if _, found := db.FindOne("invites", map[string]any{"email": "bob@bob.com"}); !found {
		t.Fatal("invite should still be pending, not expired")
	}
}

func TestRegistrarHintedRegistrationClaimsMatchingInviteByVerifiedEmail(t *testing.T) {
	db := newRegistrarTestDB(map[string]map[string]any{
		"users": {},
		"invites": {
			"invite1": map[string]any{
				"code":    "unused-hashed-code",
				"email":   "rostered@example.com",
				"role":    "admin",
				"inviter": "inviter-id",
				"sentAt":  time.Now().Unix(),
				"userid":  "invited-user-id",
			},
		},
	})
	r := NewRegistrar(db, WithDefaultRoleName("admin"), WithDisablePublicRegistration(true))
	hint := newRosterableHint(func(h *registrationhints.Hint) { h.EmailVerified = true })

	result, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if !result.InviteUsed {
		t.Fatal("result.InviteUsed = false, want true")
	}
	if result.User.GetID() != "invited-user-id" {
		t.Fatalf("user id = %q, want invited-user-id", result.User.GetID())
	}
	if _, found := db.FindOne("invites", map[string]any{"email": "rostered@example.com"}); found {
		t.Fatal("invite was not expired")
	}
}

func TestRegistrarHintedRegistrationIgnoresInviteWhenEmailUnverified(t *testing.T) {
	db := newRegistrarTestDB(map[string]map[string]any{
		"users": {},
		"invites": {
			"invite1": map[string]any{
				"code":    "unused-hashed-code",
				"email":   "rostered@example.com",
				"role":    "admin",
				"inviter": "inviter-id",
				"sentAt":  time.Now().Unix(),
				"userid":  "invited-user-id",
			},
		},
	})
	r := NewRegistrar(db, WithDefaultRoleName("admin"))
	// EmailVerified defaults to false: an unverified provider email must not
	// be able to claim someone else's invite.
	hint := newRosterableHint()

	result, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if result.InviteUsed {
		t.Fatal("result.InviteUsed = true, want false")
	}
	if result.User.GetID() == "invited-user-id" {
		t.Fatal("user was assigned the invite's pinned user id despite unverified email")
	}
	if _, found := db.FindOne("invites", map[string]any{"email": "rostered@example.com"}); !found {
		t.Fatal("invite should not have been expired")
	}
}

func TestRegistrarHintedRegistrationCreatesUserAndAuthLink(t *testing.T) {
	t.Parallel()

	db := newRegistrarTestDB(nil)
	r := NewRegistrar(db, WithDefaultRoleName("admin"), WithDisablePublicRegistration(true))
	hint := newRosterableHint()

	result, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}

	if result.Method != "hint" {
		t.Fatalf("Method = %q, want %q", result.Method, "hint")
	}
	if result.Provider != hint.ProviderName {
		t.Fatalf("Provider = %q, want %q", result.Provider, hint.ProviderName)
	}
	if !result.Background {
		t.Fatal("Background = false, want true")
	}
	if !result.CreatedAuthLink {
		t.Fatal("CreatedAuthLink = false, want true")
	}
	if result.User.GetUsername() != hint.Email {
		t.Fatalf("username = %q, want %q", result.User.GetUsername(), hint.Email)
	}
	if !result.User.GetPasswordInfo().Passwordless {
		t.Fatal("expected hinted user to be passwordless")
	}

	link, found := db.FindOne("auth_links", map[string]any{"subject": hint.Subject})
	if !found {
		t.Fatal("auth link was not created")
	}
	linkMap := link.(map[string]any)
	if linkMap["provider"] != hint.ProviderName {
		t.Fatalf("link provider = %v, want %q", linkMap["provider"], hint.ProviderName)
	}
	if linkMap["issuer"] != hint.Issuer {
		t.Fatalf("link issuer = %v, want %q", linkMap["issuer"], hint.Issuer)
	}
	if linkMap["user_id"] != result.User.GetID() {
		t.Fatalf("link user_id = %v, want %q", linkMap["user_id"], result.User.GetID())
	}
}

func TestRegistrarHintedRegistrationBypassesDisabledPublicRegistration(t *testing.T) {
	t.Parallel()

	db := newRegistrarTestDB(nil)
	r := NewRegistrar(db, WithDefaultRoleName("admin"), WithDisablePublicRegistration(true))

	// Password registration without an invite must still be blocked.
	_, err := r.register(context.Background(), register_domain.Request{
		Username: "kenton",
		Password: "password123",
		Email:    "kenton@example.com",
	})
	if !errors.Is(err, register_domain.ErrPublicRegistrationDisabled) {
		t.Fatalf("password registration error = %v, want ErrPublicRegistrationDisabled", err)
	}

	hint := newRosterableHint()
	result, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if err != nil {
		t.Fatalf("hinted registration returned error: %v", err)
	}
	if result.User == nil {
		t.Fatal("expected user to be created for valid hint")
	}
}

func TestRegistrarHintedRegistrationRejectsIncompleteHint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		hint registrationhints.Hint
	}{
		{
			name: "missing provider",
			hint: newRosterableHint(func(h *registrationhints.Hint) { h.ProviderName = "" }),
		},
		{
			name: "missing issuer",
			hint: newRosterableHint(func(h *registrationhints.Hint) { h.Issuer = "" }),
		},
		{
			name: "missing subject",
			hint: newRosterableHint(func(h *registrationhints.Hint) { h.Subject = "" }),
		},
		{
			name: "not rosterable",
			hint: newRosterableHint(func(h *registrationhints.Hint) { h.Rosterable = false }),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newRegistrarTestDB(nil)
			r := NewRegistrar(db, WithDefaultRoleName("admin"))

			hint := tt.hint
			_, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
			if !errors.Is(err, register_domain.ErrRegistrationInvalidHint) {
				t.Fatalf("error = %v, want ErrRegistrationInvalidHint", err)
			}
			if len(db.Tables["auth_links"]) != 0 {
				t.Fatal("auth link should not be created for invalid hint")
			}
		})
	}
}

func TestRegistrarHintedRegistrationReplayCannotRecreateExistingUser(t *testing.T) {
	t.Parallel()

	db := newRegistrarTestDB(nil)
	r := NewRegistrar(db, WithDefaultRoleName("admin"))
	hint := newRosterableHint()

	// First use of the hint creates the user and auth link.
	first, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if err != nil {
		t.Fatalf("first register returned error: %v", err)
	}
	if len(db.Tables["users"]) != 1 {
		t.Fatalf("users count after first use = %d, want 1", len(db.Tables["users"]))
	}
	if len(db.Tables["auth_links"]) != 1 {
		t.Fatalf("auth_links count after first use = %d, want 1", len(db.Tables["auth_links"]))
	}

	// Replaying the exact same signed hint (e.g. token reuse before
	// expiration) must not create a second user or a second auth link;
	// the existing-user duplicate check is the replay defense.
	replay := hint
	_, err = r.register(context.Background(), register_domain.Request{Hint: &replay})
	if !errors.Is(err, register_domain.ErrRegistrationTaken) {
		t.Fatalf("replay error = %v, want ErrRegistrationTaken", err)
	}

	if len(db.Tables["users"]) != 1 {
		t.Fatalf("users count after replay = %d, want 1 (no duplicate user)", len(db.Tables["users"]))
	}
	if len(db.Tables["auth_links"]) != 1 {
		t.Fatalf("auth_links count after replay = %d, want 1 (no duplicate auth link)", len(db.Tables["auth_links"]))
	}
	if _, found := db.FindOne("users", map[string]any{"email": first.User.GetEmail()}); !found {
		t.Fatal("original user should remain after replay attempt")
	}
}

func TestRegistrarHintedRegistrationRejectsDuplicateEmail(t *testing.T) {
	t.Parallel()

	hint := newRosterableHint()
	db := newRegistrarTestDB(map[string]map[string]any{
		"users": {
			"u1": map[string]any{
				"id":       "u1",
				"username": hint.Email,
				"email":    hint.Email,
			},
		},
	})
	r := NewRegistrar(db, WithDefaultRoleName("admin"))

	_, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if !errors.Is(err, register_domain.ErrRegistrationTaken) {
		t.Fatalf("error = %v, want ErrRegistrationTaken", err)
	}
	if len(db.Tables["auth_links"]) != 0 {
		t.Fatal("auth link should not be created when registration fails")
	}
}

func TestRegistrarHintedRegistrationMissingEmailFailsClosed(t *testing.T) {
	t.Parallel()

	db := newRegistrarTestDB(nil)
	r := NewRegistrar(db, WithDefaultRoleName("admin"))
	hint := newRosterableHint(func(h *registrationhints.Hint) { h.Email = "" })

	_, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if !errors.Is(err, register_domain.ErrRegistrationMissingFields) && !errors.Is(err, register_domain.ErrRegistrationIllegalUsername) {
		t.Fatalf("error = %v, want missing-fields-like error", err)
	}
	if len(db.Tables["auth_links"]) != 0 {
		t.Fatal("auth link should not be created for missing email")
	}
}

func TestRegistrarHintedRegistrationEventsUseHintMethodAndProvider(t *testing.T) {
	t.Parallel()

	bus := &recordingEventBus{}
	db := newRegistrarTestDB(nil)
	r := NewRegistrar(db, WithDefaultRoleName("admin"), WithEventBus(bus))
	hint := newRosterableHint()

	_, err := r.register(context.Background(), register_domain.Request{Hint: &hint})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}

	published := bus.singlePublish(t, events.EventRegistrationSucceeded)
	payload, ok := published.event.Payload.(events.RegistrationSucceededPayload)
	if !ok {
		t.Fatalf("payload type = %T, want events.RegistrationSucceededPayload", published.event.Payload)
	}
	if payload.Method != "hint" {
		t.Fatalf("Method = %q, want %q", payload.Method, "hint")
	}
	if payload.Provider != hint.ProviderName {
		t.Fatalf("Provider = %q, want %q", payload.Provider, hint.ProviderName)
	}
	if !payload.Background {
		t.Fatal("Background = false, want true")
	}
}
