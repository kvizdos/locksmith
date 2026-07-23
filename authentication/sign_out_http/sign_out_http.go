package sign_out

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication/events"
	"github.com/kvizdos/locksmith/users"
)

type SignOutHTTP struct {
	EventBus events.Bus
}

func (m SignOutHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Clear-Site-Data", `"cookies", "storage"`)

	if m.EventBus != nil {
		if authUser, ok := r.Context().Value("authUser").(users.LocksmithUserInterface); ok {
			envelope := events.EnrichEnvelope(r.Context(), events.Envelope{
				ID:         uuid.New().String(),
				Name:       events.EventSignOut,
				OccurredAt: time.Now(),
				Source:     "signout",
				Payload: events.SignOutPayload{
					UserID: authUser.GetID(),
				},
			})
			m.EventBus.Publish(r.Context(), envelope)
		}
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
