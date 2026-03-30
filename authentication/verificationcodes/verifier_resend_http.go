package verificationcodes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/api_helpers"
	"github.com/kvizdos/locksmith/users"
)

type VerifierResendHTTP struct {
	Verifier Verifier
}

func (rr VerifierResendHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authUser, _ := r.Context().Value("authUser").(users.LocksmithUser)

	if !authUser.RequiresEmailVerification() {
		api_helpers.WriteResponse(w, map[string]any{"success": true}, http.StatusOK)
		return
	}

	if err := rr.Verifier.SendVerification(r.Context(), authUser, VerifierMethod_EMAIL, authUser.GetEmail()); err != nil {
		if errors.Is(err, ErrVerificationRateLimited) {
			api_helpers.WriteResponse(w, map[string]any{"success": false, "error": "rate limited"}, http.StatusTooManyRequests)
			return
		}

		fmt.Println("Failed to send verification code:", err.Error())
		api_helpers.WriteResponse(w, map[string]any{"success": false}, http.StatusInternalServerError)
		return
	}

	api_helpers.WriteResponse(w, map[string]any{"success": true}, http.StatusOK)
}
