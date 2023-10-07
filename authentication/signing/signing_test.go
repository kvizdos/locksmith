package signing

import (
	"testing"
)

func TestSigningPackageFlow(t *testing.T) {
	// Test Creation
	pkg, err := CreateSigningPackage()

	if err != nil {
		t.Errorf("got a weird error: %s", err)
		return
	}

	// Test Marshalling
	privKeyBase, err := pkg.MarshalPrivate()

	if err != nil {
		t.Errorf("got a weird error while marshalling: %s", err)
		return
	}

	// Test Load Private Key
	pkgLoaded, err := DecodePrivateKey(privKeyBase)

	if err != nil {
		t.Errorf("got a weird error while marshalling: %s", err)
		return
	}

	testInput := "Hello World"

	signedFirst, err := pkg.Sign(testInput)

	if err != nil {
		t.Errorf("got a weird error while signing (1): %s", err)
		return
	}

	signedAfterLoaded, err := pkgLoaded.Sign(testInput)

	if err != nil {
		t.Errorf("got a weird error while signing (2): %s", err)
		return
	}

	firstValidatedWithLoadedPubkey := pkgLoaded.Validate(testInput, signedFirst)

	if !firstValidatedWithLoadedPubkey {
		t.Error("first signature should have validated with loaded public key.")
	}

	secondValidatedWithFirstPubkey := pkg.Validate(testInput, signedAfterLoaded)

	if !secondValidatedWithFirstPubkey {
		t.Error("second signature should have validated with original public key.")
	}
}
