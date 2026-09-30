package invitations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication/events"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type HTTPInvite struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (i HTTPInvite) IsValid() bool {
	if i.Email == "" || i.Role == "" {
		return false
	}
	return true
}

type AdministrationInviteUserHandler struct {
	Bus events.Bus
}

// Requires an authUser to be passed to HTTP Context
// Preferably through SecureEndpoint middleware.
func (i AdministrationInviteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user, ok := r.Context().Value("authUser").(users.LocksmithUser)

	if !ok {
		fmt.Println("Inviting users endpoint is required to be wrapped in SecureEndpointHTTPMiddleware()")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, ok := r.Context().Value("database").(database.DatabaseAccessor)

	if !ok {
		fmt.Println("Inviting users endpoint is required to have database context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var invite HTTPInvite
	err = json.Unmarshal(body, &invite)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !invite.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inviteCode, _, err := InviteUser(db, invite.Email, invite.Role, user.ID)

	if err != nil {
		switch err.Error() {
		case "email already invited":
			w.WriteHeader(http.StatusConflict)
			return
		case "email already registered":
			w.WriteHeader(http.StatusConflict)
			return
		default:
			fmt.Printf("Error while inviting user: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if i.Bus != nil {
		envelope := events.EnrichEnvelope(r.Context(), events.Envelope{
			ID:         uuid.New().String(),
			Name:       events.EventInviteUser,
			OccurredAt: time.Now(),
			Payload: events.InviteUserPayload{
				Email:      invite.Email,
				Role:       invite.Role,
				InviteCode: inviteCode,
			},
		})

		if err := i.Bus.Publish(r.Context(), envelope); err != nil {
			fmt.Println("Error publishing event:", err)
		}
	}

	w.Write([]byte("OK"))
}
