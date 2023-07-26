package endpoints

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication/validation"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/launchpad"
)

type EndpointSecurityBasicAuth struct {
	Enabled bool

	Username string
	Password string
}

func (e EndpointSecurityBasicAuth) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(e.Username))
			expectedPasswordHash := sha256.Sum256([]byte(e.Password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

type EndpointSecurityOptions struct {
	// Specify required permissions to hit the endpoint
	// Handlers can check permissions by themselves after this point
	// for any conditional requirements.
	MinimalPermissions []string
	// Eventually, add:
	// AllowAPITokens bool
	// If enabled, the API Key Management system will validate the permissions of the token
	BasicAuth EndpointSecurityBasicAuth
}

func SecureEndpointHTTPMiddleware(next http.Handler, db database.DatabaseAccessor, opts ...EndpointSecurityOptions) http.Handler {
	var secureOptions EndpointSecurityOptions
	if len(opts) == 0 {
		secureOptions = EndpointSecurityOptions{}
	} else {
		secureOptions = opts[0]
	}
	if secureOptions.BasicAuth.Enabled {
		return secureOptions.BasicAuth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "database", db))
			next.ServeHTTP(w, r)

		}))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject the database into the request
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
		r = r.WithContext(context.WithValue(r.Context(), "database", db))

		// next.ServeHTTP(w, r)

		launchpad.LaunchpadRequestMiddleware(next).ServeHTTP(w, r)
	})
}
