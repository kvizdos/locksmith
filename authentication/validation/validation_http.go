package validation

import (
	"context"
	"net/http"
	"time"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
)

func ValidateUserTokenMiddleware(next http.Handler, db database.DatabaseAccessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject the database into the request
		r = r.WithContext(context.WithValue(r.Context(), "database", db))

		// Validate token
		token, err := r.Cookie("token")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		parsedToken, err := authentication.ParseToken(token.Value)

		if err != nil {
			c := &http.Cookie{
				Name:     "token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}

			http.SetCookie(w, c)
			// Add logging that a fake token was passed
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, validated, err := ValidateToken(parsedToken, db)

		if err != nil {
			c := &http.Cookie{
				Name:    "token",
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),

				HttpOnly: true,
			}

			http.SetCookie(w, c)
			// Add logging that a fake token was passed
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !validated {
			c := &http.Cookie{
				Name:    "token",
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),

				HttpOnly: true,
			}

			http.SetCookie(w, c)
			// Add logging that a FORGED token was passed
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "authUser", user))

		next.ServeHTTP(w, r)
	})
}
