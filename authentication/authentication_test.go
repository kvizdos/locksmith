package authentication

import (
	"testing"

	"github.com/kvizdos/locksmith/roles"
)

func TestMain(m *testing.M) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
		"user": {
			"view.admin",
			"user.delete.self",
		},
	}

	m.Run()

	roles.AVAILABLE_ROLES = map[string][]string{}
}

func TestValidateMalformedLocksmithPasswordNoPasswordOrSalt(t *testing.T) {
	password := PasswordInfo{}

	passed, err := ValidatePassword(password, "randompassword")

	if passed {
		t.Fatalf("Password validated with malformed locksmith password (no password or salt presented)")
		return
	}

	if err.Error() != "locksmith password not presented" {
		t.Fatalf("Didn't receive correct error message: %s", err.Error())
	}
}

func TestValidateMalformedLocksmithPasswordNoSalt(t *testing.T) {
	password := PasswordInfo{
		Password: "random",
	}

	passed, err := ValidatePassword(password, "randompassword")

	if passed {
		t.Fatalf("Password validated with malformed locksmith password (no salt presented)")
		return
	}

	if err.Error() != "locksmith salt not presented" {
		t.Fatalf("Didn't receive correct error message: %s", err.Error())
	}
}

func TestValidateInvalidInputPasswordLength(t *testing.T) {
	password := PasswordInfo{
		Password: "random",
		Salt:     "blah",
	}

	passed, err := ValidatePassword(password, "")

	if passed {
		t.Fatalf("Password validated with malformed input password (no input password sent)")
		return
	}

	if err.Error() != "no input password passed" {
		t.Fatalf("Didn't receive correct error message: %s", err.Error())
	}
}

func TestValidateInvalidPassword(t *testing.T) {
	password := PasswordInfo{
		Password: "37d80843ed0aff6f6d0fe7ea3bec9ef65a7abd20b9657b78a2e32bbb0ac1293d",
		Salt:     "b16d14c05b76596006e3d2edf5eabca0",
	}

	passed, err := ValidatePassword(password, "wrongpassword")

	if err != nil {
		t.Fatalf("Received unexpected error while testing invalid password: %s\n", err.Error())
		return
	}

	if passed {
		t.Fatal("Wrong password passed when it should've failed.")
		return
	}
}

func TestValidateValidPassword(t *testing.T) {
	password := PasswordInfo{
		Password: "37d80843ed0aff6f6d0fe7ea3bec9ef65a7abd20b9657b78a2e32bbb0ac1293d",
		Salt:     "b16d14c05b76596006e3d2edf5eabca0",
	}

	passed, err := ValidatePassword(password, "supersecurepassword123")

	if err != nil {
		t.Fatalf("Received unexpected error while testing invalid password: %s\n", err.Error())
		return
	}

	if !passed {
		t.Fatal("Wrong password passed when it should've failed.")
		return
	}
}
