package verificationcodes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kvizdos/locksmith/api_helpers"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type VerifierExchangeHTTP struct {
	Verifier Verifier
}

type bod struct {
	Code string `json:"code"`
}

func (rr VerifierExchangeHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authUser, _ := r.Context().Value("authUser").(users.LocksmithUser)

	if !authUser.RequiresEmailVerification() {
		api_helpers.WriteResponse(w, map[string]any{"success": true}, http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var codeBody bod
	err = json.Unmarshal(body, &codeBody)
	if err != nil {
		// handle the error
		fmt.Println("Error unmarshalling request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if codeBody.Code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified, err := rr.Verifier.CheckCode(r.Context(), VerifierMethod_EMAIL, authUser.GetID(), codeBody.Code)
	if err != nil {
		// handle the error
		fmt.Println("Error verifying code:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !verified {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	_, err = db.UpdateOne("users", map[string]any{
		"id": authUser.GetID(),
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: {
			"emailVerified":           true,
			"emailVerificationMethod": "email",
			"needsEmailVerification":  false,
		},
	})
	if err != nil {
		// handle the error
		fmt.Println("Error updating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rr.Verifier.DeleteCode(r.Context(), VerifierMethod_EMAIL, authUser.GetID())
	if err != nil {
		// handle the error
		fmt.Println("Error deleting code:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	api_helpers.WriteResponse(w, map[string]any{"success": true}, http.StatusOK)
}
