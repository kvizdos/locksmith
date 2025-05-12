package invitations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/database"
)

type AdministrationInviteListHandler struct{}

// Requires an authUser to be passed to HTTP Context
// Preferably through SecureEndpoint middleware.
func (i AdministrationInviteListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	db, ok := r.Context().Value("database").(database.DatabaseAccessor)

	if !ok {
		fmt.Println("Inviting users endpoint is required to have database context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invites := ListInvites(db)

	js, _ := json.Marshal(invites)
	w.Write(js)
}
