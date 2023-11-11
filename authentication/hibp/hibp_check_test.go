package hibp

import (
	"net/http"
	"testing"

	"github.com/kvizdos/locksmith/authentication"
)

func TestCheckPasswordHIBPIntegrationIsBreached(t *testing.T) {
	hasBeenBreached := make(chan bool)
	httpClient := &http.Client{}

	go CheckPassword("Locksmith Integration Tests", "password123", hasBeenBreached, httpClient)
	isBreached := <-hasBeenBreached

	if !isBreached {
		t.Error("expected to fail!")
	}
}

func TestCheckPasswordHIBPIntegrationIsNotBreached(t *testing.T) {
	hasBeenBreached := make(chan bool)
	pass, _ := authentication.GenerateRandomString(128)
	httpClient := &http.Client{}
	go CheckPassword("Locksmith Integration Tests", pass, hasBeenBreached, httpClient)
	isBreached := <-hasBeenBreached

	if isBreached {
		t.Error("expected to pass!")
	}
}
