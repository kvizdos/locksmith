package webauth

import (
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"kv.codes/locksmith/users"
)

func FinishRegistration(w http.ResponseWriter, r *http.Request, user users.LocksmithUserInterface) {
	response, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	if err != nil {
		// Handle Error and return.

		return
	}

	// Get the session data stored from the function above
	sessions := user.GetWebAuthnSessions()

	if len(sessions) == 0 {
		// RETURN FAILED STATE
		return
	}

	credential, err := wauth.CreateCredential(user, sessions[0], response)
	if err != nil {
		// Handle Error and return.

		return
	}

	// If creation was successful, store the credential object
	// JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps

	// Add the user credential.
	user.AddNewWebAuthnCredential(credential)
}
