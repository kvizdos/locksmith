//go:build !testing
// +build !testing

package webauth

import (
	"fmt"
	"os"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

var (
	wauth *webauthn.WebAuthn
	werr  error
)

// InitializeWebAuthn() MUST be called before anything else is run.
func InitializeWebAuthn() {
	if !validateRequiredEnvironmentVariables() {
		fmt.Println("MISSING REQUIRED ENVIRONMENT VARIABLES FOR AUTHENTICATION: AppName, DomainName, AuthOrigin")
	}
	wconfig := &webauthn.Config{
		RPDisplayName: os.Getenv("AppName"),              // Display Name for your site
		RPID:          os.Getenv("DomainName"),           // Generally the FQDN for your site
		RPOrigins:     []string{os.Getenv("AuthOrigin")}, // The origin URLs allowed for WebAuthn requests
	}

	if wauth, werr = webauthn.New(wconfig); werr != nil {
		fmt.Println(werr)
	}
}

func validateRequiredEnvironmentVariables() bool {
	// Uses Getenv instead of LookupEnv for testing purposes.
	rp := os.Getenv("AppName")
	dn := os.Getenv("DomainName")
	ao := os.Getenv("AuthOrigin")

	if len(rp) == 0 || len(dn) == 0 || len(ao) == 0 {
		return false
	}

	return true
}

func BeginRegisterWebAuthn(user users.LocksmithUserInterface, db database.DatabaseAccessor) (*protocol.CredentialCreation, error) {
	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.AuthenticatorAttachment("platform"),
		RequireResidentKey:      protocol.ResidentKeyRequired(),
		UserVerification:        protocol.VerificationRequired,
	}

	conveyancePref := protocol.PreferDirectAttestation

	credential, session, err := wauth.BeginRegistration(user, webauthn.WithAuthenticatorSelection(authSelect), webauthn.WithConveyancePreference(conveyancePref))

	if err != nil {
		return &protocol.CredentialCreation{}, fmt.Errorf("failed to create credential: %s", err.Error())
	}

	// SAVE SESSION TO DATABASE
	db.UpdateOne("users", map[string]interface{}{
		"username": user.GetUsername(),
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"websessions": session,
		},
	})

	return credential, nil
}
