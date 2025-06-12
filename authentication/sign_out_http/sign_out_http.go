package sign_out

import (
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/launchpad"
)

type SignOutHTTP struct{}

func (m SignOutHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Attach Session Cookie
	cookie := http.Cookie{Name: "token", Value: "", Expires: time.Unix(0, 0), HttpOnly: true, Secure: true, Path: "/"}
	sessionExpiresAtCookie := http.Cookie{Name: "ls_expires_at", Value: "", Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"}
	oauthprovidercookie := http.Cookie{Name: "ls_oauth_provider", Value: "", Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"}
	launchpadcookie := http.Cookie{Name: "LaunchpadUser", Value: "", Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"}

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &sessionExpiresAtCookie)
	http.SetCookie(w, &oauthprovidercookie)

	if launchpad.IS_ENABLED {
		http.SetCookie(w, &launchpadcookie)
	}

	http.SetCookie(w, &http.Cookie{Name: "magic", Value: "", Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "sid", Value: "", Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"})

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
