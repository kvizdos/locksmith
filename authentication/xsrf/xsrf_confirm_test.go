package xsrf

import (
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/signing"
)

func TestMain(m *testing.M) {
	pkg, _ := signing.CreateSigningPackage()
	XSRFSigningPackage.Anonymous = &pkg
	XSRFSigningPackage.Authenticated = &pkg
	m.Run()
}
func TestConfirmMismatchedIDs(t *testing.T) {
	xsrf, _ := GenerateXSRFForSession("test-id", 15*time.Minute)
	confirmed := Confirm(xsrf, "test-id-2")

	if confirmed {
		t.Errorf("should not have confirmed")
	}
}

func TestConfirmMatchingSignatures(t *testing.T) {
	xsrf, _ := GenerateXSRFForSession("test-id", 15*time.Minute)
	confirmed := Confirm(xsrf, "test-id")

	if !confirmed {
		t.Errorf("Error confirming")
	}
}
