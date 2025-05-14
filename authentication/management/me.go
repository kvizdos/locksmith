package management

import (
	"log"
	"net/http"

	"github.com/kvizdos/locksmith/api_helpers"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/users"
)

type meEndpointHTTP struct{}

func (m meEndpointHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("authUser").(users.LocksmithUserInterface)

	if !ok {
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "failed",
		}, http.StatusInternalServerError)
		return
	}

	role, err := user.GetRole()

	if err != nil {
		log.Println(err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "failed",
		}, http.StatusInternalServerError)
		return
	}

	pub, err := user.ToPublic()

	out := map[string]interface{}{
		"info": map[string]any{
			"id":       user.GetID(),
			"username": user.GetUsername(),
			"email":    user.GetEmail(),
			"role":     role.Name,
		},
		"permissions": roles.GetExposable(role.Permissions),
	}

	if err != nil {
		api_helpers.WriteResponse(w, out, http.StatusOK)
		return
	}

	u := pub.MakeUserSafe()

	api_helpers.WriteResponse(w, map[string]interface{}{
		"info":        u,
		"permissions": roles.GetExposable(role.Permissions),
	}, http.StatusOK)
}
