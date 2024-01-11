package providers_test

import (
	"reflect"
	"testing"

	"github.com/kvizdos/locksmith/shared-memory/objects"
	"github.com/kvizdos/locksmith/shared-memory/providers"
)

func TestRamMemorySettingAndGetting(t *testing.T) {
	mem := providers.NewRamSharedMemoryProvider()

	attempt := objects.UserLoginAttempt{
		Attempts:    1,
		LastAttempt: 100,
	}

	_ = mem.SetMemory(objects.USER_LOGIN_ATTEMPTS, "user-id", attempt)

	alloc, found := mem.GetFromMemory(objects.USER_LOGIN_ATTEMPTS, "user-id")

	if !found {
		t.Error("Expected to find result.")
		return
	}

	if !reflect.DeepEqual(attempt, alloc) {
		t.Error("Expected values to be equal.")
	}
}

func TestRamMemoryReturnsNotFound(t *testing.T) {
	mem := providers.NewRamSharedMemoryProvider()

	attempt := objects.UserLoginAttempt{
		Attempts:    1,
		LastAttempt: 100,
	}

	_ = mem.SetMemory(objects.USER_LOGIN_ATTEMPTS, "user-id", attempt)

	_, found := mem.GetFromMemory(objects.USER_LOGIN_ATTEMPTS, "diff-id")

	if found {
		t.Error("Expected to NOT find result.")
		return
	}
}

func TestRamMemoryIncrement(t *testing.T) {
	mem := providers.NewRamSharedMemoryProvider()

	err := mem.Increment(objects.USER_LOGIN_ATTEMPTS, "user-id")

	if err != nil {
		t.Error(err)
		return
	}

	attempts, found := mem.GetFromMemory(objects.USER_LOGIN_ATTEMPTS, "user-id")

	if !found {
		t.Error("Expected to find result.")
		return
	}

	if attempts.(objects.UserLoginAttempt).Attempts != 1 {
		t.Errorf("Expected to get 1 attempt, got %d", attempts.(objects.UserLoginAttempt).Attempts)
	}
}

func TestRamMemoryIncrementAlreadyPulled(t *testing.T) {
	mem := providers.NewRamSharedMemoryProvider()

	attempt := objects.UserLoginAttempt{
		Attempts:    1,
		LastAttempt: 100,
	}

	_ = mem.SetMemory(objects.USER_LOGIN_ATTEMPTS, "user-id", attempt)

	pulledAttempts, _ := mem.GetFromMemory(objects.USER_LOGIN_ATTEMPTS, "user-id")

	// In regular code, if checks and things would go here.

	err := mem.Increment(objects.USER_LOGIN_ATTEMPTS, "user-id", pulledAttempts)

	if err != nil {
		t.Error(err)
		return
	}

	attempts, found := mem.GetFromMemory(objects.USER_LOGIN_ATTEMPTS, "user-id")

	if !found {
		t.Error("Expected to find result.")
		return
	}

	if attempts.(objects.UserLoginAttempt).Attempts != 2 {
		t.Errorf("Expected to get 1 attempt, got %d", attempts.(objects.UserLoginAttempt).Attempts)
	}
}
