package endpoints

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/validation"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/ratelimits"
	"github.com/kvizdos/locksmith/tenant"
	"github.com/kvizdos/locksmith/users"
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

// Returns int (status code)
// Will only let the request continue if
// the status is "200"
type EndpointSecurityCustomMiddleware func(users.LocksmithUserInterface, database.DatabaseAccessor) int

type EndpointSecurityOptions struct {
	// Specify required permissions to hit the endpoint
	// Handlers can check permissions by themselves after this point
	// for any conditional requirements.
	MinimalPermissions []string
	// Define the Entitlement(s) that aligns to the permission
	RequiresEntitlement []string
	// Eventually, add:
	// AllowAPITokens bool
	// If enabled, the API Key Management system will validate the permissions of the token
	BasicAuth EndpointSecurityBasicAuth
	// If you'd like to unwrap the Locksmith
	// context user into some other LocksmithUserInterface,
	// type it ehre.
	CustomUser validation.JWTCustomUserOptions
	// Define a custom Tenant interface
	// if you plan on being multi-tenant.
	TenantInterface tenant.Tenant
	// After initial confirmation of a user is confirmed,
	// you can use this function to validate endpoint-specific
	// validations.
	SecondaryValidation EndpointSecurityCustomMiddleware
	// PrioritizeMagic determines the precedence of authentication methods within the SecureEndpointMiddleware.
	// By default, a logged-in user's session (via the `token` cookie) is prioritized over the Magic Authentication Code (`magic` cookie).
	// Setting this to `true` gives preference to the `magic` cookie, even if a valid `token` cookie exists.
	// Please note that this action will revoke any permissions not explicitly specified in the `magic` cookie.
	PrioritizeMagic bool
	// Optionally, rate limit the endpoint.
	RateLimit *ratelimits.RateLimiter
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

	if len(secureOptions.RequiresEntitlement) == 0 {
		panic("secureOptions.RequiresEntitlement CANNOT be blank.")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sidCookie, err := r.Cookie("sid")
		sidValue := ""
		if err != nil {
			sid, err := authentication.GenerateRandomString(16)

			if err != nil {
				fmt.Println("Error generating random SID", err)
				sid = "random"
			}

			sidValue = sid

			cookie := http.Cookie{Name: "sid", Value: sid, HttpOnly: false, Secure: true, Path: "/"}
			http.SetCookie(w, &cookie)
		} else {
			sidValue = sidCookie.Value
		}

		var userInterface users.LocksmithUserInterface
		var userClaims jwt.Claims

		if secureOptions.CustomUser.CustomUser != nil {
			userInterface = secureOptions.CustomUser.CustomUser
			userClaims = secureOptions.CustomUser.Claims
		} else {
			userInterface = users.LocksmithUser{}
			userClaims = &users.BaseValidationClaims{}
		}

		magicQuery := r.URL.Query().Get("magic")
		if magicQuery == "" {
			magicCookie, err := r.Cookie("magic")
			if err == nil {
				magicQuery = magicCookie.Value
			}
		}

		user, err := validation.ValidateHTTPUserToken(r, db, validation.MagicValidation{
			Token:      magicQuery,
			Prioritize: secureOptions.PrioritizeMagic,
		}, validation.HTTPValidationOptions{
			UserType:       userInterface,
			Claims:         userClaims,
			ValidationKeys: validation.GetSigningKeys(),
		})

		if err != nil {
			fmt.Println("Failed to validate", err)
			c := &http.Cookie{
				Name:     "token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}
			mc := &http.Cookie{
				Name:     "magic",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}

			http.SetCookie(w, c)
			http.SetCookie(w, mc)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.URL.Query().Get("magic") != "" {
			cookie := http.Cookie{Name: "magic", Value: r.URL.Query().Get("magic"), HttpOnly: true, Secure: true, Path: "/"}
			http.SetCookie(w, &cookie)
		}

		if secureOptions.TenantInterface != nil {
			tenantInfo, err := tenant.GetTenantFromID(user.GetTenantID(), db, secureOptions.TenantInterface)
			if err != nil {
				fmt.Println("Failed to get Tenant ID!", user.GetTenantID(), err)
			} else {
				user.SetTenant(tenantInfo)
			}
		}

		entitlements := user.GetAttachedEntitlements()

		for _, entitlement := range secureOptions.RequiresEntitlement {
			if !slices.Contains(entitlements, entitlement) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		if len(secureOptions.MinimalPermissions) > 0 {
			for _, permission := range secureOptions.MinimalPermissions {
				if !user.HasPermission(permission) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), "authUser", user))
		r = r.WithContext(context.WithValue(r.Context(), "database", db))
		r = r.WithContext(context.WithValue(r.Context(), "sid", sidValue))

		if secureOptions.SecondaryValidation != nil {
			statusCode := secureOptions.SecondaryValidation(user, db)
			if statusCode != 200 {
				w.WriteHeader(statusCode)
				return
			}
		}

		// Final Rate Limiter for the Exact User
		// This provides a concrete answer about whether
		// or not the user should be able to access a resource.
		if secureOptions.RateLimit != nil {
			if !secureOptions.RateLimit.CanHandle(user.GetID()) {
				fmt.Println("Rate limited here!")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		}

		// launchpad.LaunchpadRequestMiddleware(next).ServeHTTP(w, r)
		next.ServeHTTP(w, r)
	})
}
