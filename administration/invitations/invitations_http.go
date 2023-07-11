package invitations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

type AdministrationInviteUserHandler struct{}

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

	body, err := ioutil.ReadAll(r.Body)
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

	inviteCode, err := InviteUser(db, invite.Email, invite.Role, user.ID)

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

	w.Write([]byte(inviteCode))
}
