package authenticator

import (
	"log/slog"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/database"
)

func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s: expected panic, got none", name)
		}
	}()
	fn()
}

func TestNewAuthorizerPanicsWithoutTokenManager(t *testing.T) {
	t.Parallel()

	db := newTestDB(nil)
	sp, _ := signing.CreateSigningPackage()

	mustPanic(t, "missing token manager", func() {
		NewAuthorizer(db,
			WithMethods(AllowMethodPassword()),
			WithSigningPackage(&sp),
		)
	})
}

func TestNewAuthorizerPanicsWithoutMethods(t *testing.T) {
	t.Parallel()

	db := newTestDB(nil)
	sp, _ := signing.CreateSigningPackage()

	mustPanic(t, "missing methods", func() {
		NewAuthorizer(db,
			WithTokenManager(&mockTokenManager{}),
			WithSigningPackage(&sp),
		)
	})
}

func TestNewAuthorizerPanicsWithoutSigningPackage(t *testing.T) {
	t.Parallel()

	db := newTestDB(nil)

	mustPanic(t, "missing signing package", func() {
		NewAuthorizer(db,
			WithTokenManager(&mockTokenManager{}),
			WithMethods(AllowMethodPassword()),
		)
	})
}

func TestNewAuthorizerSucceedsWithRequiredDeps(t *testing.T) {
	t.Parallel()

	db := newTestDB(nil)
	sp, _ := signing.CreateSigningPackage()

	a := NewAuthorizer(db,
		WithTokenManager(&mockTokenManager{}),
		WithMethods(AllowMethodPassword()),
		WithSigningPackage(&sp),
	)

	if a == nil {
		t.Fatal("expected non-nil authorizer")
	}
}

func TestOptionWithLogger(t *testing.T) {
	t.Parallel()

	customLogger := slog.Default()
	a := newTestAuthorizer(newTestDB(nil), WithLogger(customLogger))

	if a.log != customLogger {
		t.Fatal("expected logger to be set via WithLogger")
	}
}

func TestOptionWithRedirectPath(t *testing.T) {
	t.Parallel()

	a := newTestAuthorizer(newTestDB(nil), WithRedirectPath("/custom-app"))

	if a.redirectPath != "/custom-app" {
		t.Fatalf("expected redirectPath '/custom-app', got %s", a.redirectPath)
	}
}

func TestOptionWithTokenManager(t *testing.T) {
	t.Parallel()

	tm := &mockTokenManager{redirectPath: "/somewhere"}
	db := newTestDB(nil)
	sp, _ := signing.CreateSigningPackage()

	a := NewAuthorizer(db,
		WithTokenManager(tm),
		WithMethods(AllowMethodPassword()),
		WithSigningPackage(&sp),
	)

	if a.tm != tm {
		t.Fatal("expected token manager to be set via WithTokenManager")
	}
}

func TestOptionWithMethods(t *testing.T) {
	t.Parallel()

	db := newTestDB(nil)
	sp, _ := signing.CreateSigningPackage()
	method1 := AllowMethodPassword()

	a := NewAuthorizer(db,
		WithTokenManager(&mockTokenManager{}),
		WithMethods(method1),
		WithSigningPackage(&sp),
	)

	if len(a.methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(a.methods))
	}
}

func TestOptionWithMethodsAppends(t *testing.T) {
	t.Parallel()

	db := newTestDB(nil)
	sp, _ := signing.CreateSigningPackage()

	a := NewAuthorizer(db,
		WithTokenManager(&mockTokenManager{}),
		WithMethods(AllowMethodPassword()),
		WithMethods(AllowMethodPassword()),
		WithSigningPackage(&sp),
	)

	if len(a.methods) != 2 {
		t.Fatalf("expected 2 methods after appending, got %d", len(a.methods))
	}
}

func TestOptionWithSigningPackage(t *testing.T) {
	t.Parallel()

	sp, _ := signing.CreateSigningPackage()
	a := newTestAuthorizer(newTestDB(nil), WithSigningPackage(&sp))

	if a.sp != &sp {
		t.Fatal("expected signing package to be set via WithSigningPackage")
	}
}

func TestOptionWithMinimumResponseTime(t *testing.T) {
	t.Parallel()

	a := newTestAuthorizer(newTestDB(nil), WithMinimumResponseTime(250*time.Millisecond))

	if a.minimumResponseTime != 250*time.Millisecond {
		t.Fatalf("expected minimumResponseTime 250ms, got %v", a.minimumResponseTime)
	}
}

func TestOptionWithEmailAsUsername(t *testing.T) {
	t.Parallel()

	a := newTestAuthorizer(newTestDB(nil), WithEmailAsUsername())

	if !a.emailAsUsername {
		t.Fatal("expected emailAsUsername to be true")
	}
}

func TestOptionDisableUserEnumerationProtection(t *testing.T) {
	t.Parallel()

	a := newTestAuthorizer(newTestDB(nil), DisableUserEnumerationProtection())

	if !a.disableUserEnumerationProtection {
		t.Fatal("expected disableUserEnumerationProtection to be true")
	}
}

func TestGetUsernameNoun(t *testing.T) {
	t.Parallel()

	withUsername := newTestAuthorizer(newTestDB(nil))
	if withUsername.getUsernameNoun() != "Username" {
		t.Fatalf("expected 'Username', got %s", withUsername.getUsernameNoun())
	}

	withEmail := newTestAuthorizer(newTestDB(nil), WithEmailAsUsername())
	if withEmail.getUsernameNoun() != "Email" {
		t.Fatalf("expected 'Email', got %s", withEmail.getUsernameNoun())
	}
}

var _ database.DatabaseAccessor = database.TestDatabase{}
var _ authenticator_domain.Handler = (*mockFederatedHandlerFull)(nil)
