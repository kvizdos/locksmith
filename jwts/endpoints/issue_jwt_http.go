package jwt_endpoints

import (
	"log"
	"net/http"

	"github.com/kvizdos/locksmith/api_helpers"
	"github.com/kvizdos/locksmith/jwts"
	"github.com/kvizdos/locksmith/users"
)

type issueJWTHTTP struct{}

func (m issueJWTHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	requestedJWT := r.URL.Query().Get("token")

	jwt, exists := jwts.GetJWT(requestedJWT)

	if !exists {
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "token id not enabled",
		}, http.StatusNotFound)
		return
	}

	if !role.HasPermission(jwt.RequiredPermission) {
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "permission denied",
		}, http.StatusForbidden)
		return
	}

	tokenStr, err := jwt.IssueJWT(r.Context(), jwts.IssueJWTOptions{
		Sub: user.GetID(),
		Req: r,
	})
	if err != nil {
		log.Println("issue jwt:", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{Reason: "failed to issue jwt"}, http.StatusInternalServerError)
		return
	}

	api_helpers.WriteResponse(w, map[string]string{"token": tokenStr}, http.StatusOK)
}
