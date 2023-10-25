package reset

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type ResetRouterAPIHandler struct {
	Database       database.DatabaseAccessor
	SendResetToken func(token string, user users.LocksmithUserInterface)
}

/*
Flow:
- [ ] User clicks "Forgot password" on login page
- [x] User enters their email
- [x] Show a screen to the user that "if the account exists, we've sent a link to your email address."
  - [x] Sends POST request to create & dispatch the reset token
  - POST /api/reset-password?email=email

- [x] Create a MAC for the user to access the PUT /api/reset-password endpoint
- [x] Send them a Notification with the MAC (we should make the SendMessage a variable on ResetRouterAPIHandler)
  - URL Format: /reset-password/reset?magic=<MAC>

- [x] User will enter their new password
- [x] Password will be changed
  - [x] MAC is expired
  - PUT /api/reset-password { email: email }

- [x] Redirect to /login
*/
func (h ResetRouterAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		// This page is only accessible through the
		// Magic Access Code.
		// Updates the actual password
		endpoints.SecureEndpointHTTPMiddleware(ResetPasswordAPIHandler{}, h.Database, endpoints.EndpointSecurityOptions{
			MinimalPermissions: []string{"magic.reset.password"},
			PrioritizeMagic:    true,
		}).ServeHTTP(w, r)
		return
	case http.MethodPost:
		// Creates the MAC if the user exists.
		username := r.URL.Query().Get("username")

		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, found := h.Database.FindOne("users", map[string]interface{}{
			"username": username,
		})

		if !found {
			w.WriteHeader(http.StatusOK)
			return
		}

		var lsUser users.LocksmithUserInterface
		users.LocksmithUser{}.ReadFromMap(&lsUser, user.(map[string]interface{}))
		token, err := lsUser.CreateMagicAuthenticationCode(h.Database, magic.MagicAuthenticationVariables{
			ForUserID:          lsUser.GetID(),
			AllowedPermissions: []string{"magic.reset.password"},
			TTL:                30 * time.Minute,
		})

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.SendResetToken(token, lsUser)

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
