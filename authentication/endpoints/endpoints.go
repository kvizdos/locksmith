package endpoints

import (
	"context"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication/validation"
	"github.com/kvizdos/locksmith/database"
)

type EndpointSecurityOptions struct {
	// Specify required permissions to hit the endpoint
	// Handlers can check permissions by themselves after this point
	// for any conditional requirements.
	MinimalPermissions []string
	// Eventually, add:
	// AllowAPITokens bool
	// If enabled, the API Key Management system will validate the permissions of the token
}

func SecureEndpointHTTPMiddleware(next http.Handler, db database.DatabaseAccessor, opts ...EndpointSecurityOptions) http.Handler {
	var secureOptions EndpointSecurityOptions
	if len(opts) == 0 {
		secureOptions = EndpointSecurityOptions{}
	} else {
		secureOptions = opts[0]
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject the database into the request
		r = r.WithContext(context.WithValue(r.Context(), "database", db))

		user, err := validation.ValidateHTTPUserToken(r, db)

		if err != nil {
			c := &http.Cookie{
				Name:     "token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}

			http.SetCookie(w, c)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if len(secureOptions.MinimalPermissions) > 0 {
			userRole, err := user.GetRole()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, permission := range secureOptions.MinimalPermissions {
				if !userRole.HasPermission(permission) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), "authUser", user))

		next.ServeHTTP(w, r)
	})
}
