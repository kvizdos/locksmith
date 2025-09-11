package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kvizdos/locksmith/authentication/validation"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/observability"
	"github.com/kvizdos/locksmith/users"
)

func LoginUser(db database.DatabaseAccessor, user users.LocksmithUserInterface, provider string, redirectPage string, LoginInfoCallback func(method string, user map[string]any), w http.ResponseWriter, r *http.Request) {
	if tokencookie, err := r.Cookie("token"); err == nil && tokencookie.Value != "" {
		expiresAtString, expiresErr := r.Cookie("ls_expires_at")

		if expiresErr == nil {
			expiresAtUnix, conversionErr := strconv.ParseInt(expiresAtString.Value, 10, 64)

			if conversionErr == nil {
				t := time.Unix(expiresAtUnix, 0).UTC()
				isBefore := time.Now().UTC().Add(11 * time.Minute).Before(t)
				// No refresh necessary
				if isBefore {
					u, err := validation.ValidateHTTPUserToken(r, db, validation.MagicValidation{})
					if err == nil && u.GetID() == user.GetID() {
						http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
						return
					}
				}
			}
		}
	}

	session, err := user.GeneratePasswordSession()

	if err != nil {
		fmt.Println("Error generating session token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Expire at 11 AM UTC
	now := time.Now().UTC()
	today11 := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, time.UTC)

	var next11 time.Time
	if now.Before(today11) {
		next11 = today11
	} else {
		next11 = today11.AddDate(0, 0, 1)
	}

	session.ExpiresAt = next11.Unix()
	// session.ExpiresAt = time.Now().UTC().Add(11 * time.Minute).Unix()

	err = user.SavePasswordSession(session, db)

	if err != nil {
		fmt.Println("Error saving session token to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.LOGGER.Log(logger.LOGIN, user.GetUsername()+" (OIDC)", logger.GetIPFromRequest(*r))

	observability.LoginSuccess.Inc()

	cookieValue := user.(users.LocksmithUser).GenerateCookieValueFromSession(session)

	// Expire Login XSRF cookie
	cookieXSRF := http.Cookie{Name: "login_xsrf", Value: "", Expires: time.Unix(0, 0), HttpOnly: true, Secure: true, Path: "/api/login", SameSite: http.SameSiteStrictMode}

	// Attach Session Cookie
	cookie := http.Cookie{Name: "token", Value: cookieValue, Expires: time.Unix(session.ExpiresAt, 0), HttpOnly: true, Secure: true, Path: "/"}

	sessionExpiresAtCookie := http.Cookie{Name: "ls_expires_at", Value: fmt.Sprintf("%d", session.ExpiresAt), Expires: time.Now().UTC().AddDate(10, 0, 0), HttpOnly: false, Secure: true, Path: "/"}

	oauthprovidercookie := http.Cookie{Name: "ls_oauth_provider", Value: provider, Expires: time.Now().UTC().AddDate(10, 0, 0), HttpOnly: false, Secure: true, Path: "/"}

	oauthhint := http.Cookie{Name: "ls_oauth_hint", Value: user.GetEmail(), Expires: time.Now().UTC().AddDate(10, 0, 0), HttpOnly: true, Secure: true, Path: "/"}

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &cookieXSRF)
	http.SetCookie(w, &oauthhint)
	http.SetCookie(w, &sessionExpiresAtCookie)
	http.SetCookie(w, &oauthprovidercookie)
	redirect, err := url.QueryUnescape(redirectPage)
	if err != nil {
		http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
		return
	}

	if LoginInfoCallback != nil {
		LoginInfoCallback(fmt.Sprintf("oauth-%s", strings.ToLower(provider)), user.ToMap())
	}

	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}
