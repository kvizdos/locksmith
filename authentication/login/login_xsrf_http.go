package login

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/xsrf"
)

// This confirms that:
// - A Session ID exists, and if not generates one
// - A Login_XSRF is created to work with LoginHandler{}
type LoginPageMiddleware struct {
	Next http.Handler
}

func (h LoginPageMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionIDCookie, err := r.Cookie("sid")
	sessionIDValue := ""

	if err != nil {
		sid, err := authentication.GenerateRandomString(16)

		if err != nil {
			fmt.Println("Error generating random SID", err)
			sid = "random"
		}

		sessionIDValue = sid

		cookie := http.Cookie{Name: "sid", Value: sid, HttpOnly: false, Secure: true, Path: "/"}
		http.SetCookie(w, &cookie)
	} else {
		sessionIDValue = sessionIDCookie.Value
	}

	currentXSRF, err := r.Cookie("login_xsrf")
	if err == nil {
		xsrfConfirmed := xsrf.Confirm(currentXSRF.Value, sessionIDValue)

		if xsrfConfirmed {
			r = r.WithContext(context.WithValue(r.Context(), "login_xsrf", currentXSRF.Value))
			h.Next.ServeHTTP(w, r)
			return
		}
	}

	// Generate Login XSRF for SID
	loginXSRF, _ := xsrf.GenerateXSRFForSession(sessionIDValue, 1*time.Hour)

	cookieAPI := http.Cookie{Name: "login_xsrf", Value: loginXSRF, SameSite: http.SameSiteStrictMode, HttpOnly: true, Secure: true, Path: "/api/login"}
	http.SetCookie(w, &cookieAPI)

	r = r.WithContext(context.WithValue(r.Context(), "login_xsrf", loginXSRF))

	h.Next.ServeHTTP(w, r)
}
