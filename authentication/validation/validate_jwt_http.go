package validation

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type JWTEndpointHandler struct {
	CustomUserOptions JWTCustomUserOptions
}

type JWTCustomUserOptions struct {
	CustomUser users.LocksmithUserInterface
	// Custom Registered Claims are defined
	// in users.LocksmithUserInterface.FinalizeAccessJWTClaims.
	//
	// Each new JSON entry needs a struct representation with valid json tags.
	Claims jwt.Claims
}

func (endpoint JWTEndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// tokenCookie, err := r.Cookie("token")

	// db, ok := r.Context().Value("database").(database.DatabaseAccessor)

	// if !ok {
	// 	log.Println("failed to access database in validation.jwtendpointhandler")
	// 	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	// 	return
	// }

	// if err == nil && tokenCookie.Expires.Before(time.Now().UTC().Add(-1*time.Minute)) {
	// 	fmt.Println("in here")
	// 	w.WriteHeader(http.StatusOK) // don't issue new ones till its about to expire
	// 	// This endpoint is NOT for verifying the jwt; it exists, so securenedpointmiddleware can take care of that.
	// 	return
	// }

	var userInterface users.LocksmithUserInterface

	if endpoint.CustomUserOptions.CustomUser != nil {
		userInterface = endpoint.CustomUserOptions.CustomUser
	} else {
		userInterface = users.LocksmithUser{}
	}

	// Validate token
	profileToken, profileErr := r.Cookie("profile")
	refreshToken, refreshErr := r.Cookie("refresh")

	if profileErr != nil || refreshErr != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Do they WORK? Are they matching and non-expired?
	validated := userInterface.ValidateRefreshJWT(refreshToken.Value, profileToken.Value, GetSigningKeys().PublicKey)

	if !validated {
		log.Println("Failed to refresh token as validation failed")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": strings.ToLower(loginReq.Username),
	})

	if !usernameExists {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	var tmpUser users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

	tokens, err := user.GenerateJWTCookies("top", GetSigningKeys().PrivateKey, db)
	if err != nil {
		fmt.Println("Failed to generate JWT cookies:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Attach Session Cookie
	accessCookie := http.Cookie{Name: "token", Value: tokens.Access, Expires: tokens.RefreshExpiresAt, HttpOnly: true, Secure: true, Path: "/"}
	profileCookie := http.Cookie{Name: "profile", Value: tokens.Profile, Expires: tokens.RefreshExpiresAt, HttpOnly: false, Secure: true, Path: "/"}
	refreshCookie := http.Cookie{Name: "refresh", Value: tokens.Refresh, Expires: tokens.RefreshExpiresAt, HttpOnly: true, Secure: true, Path: "/"}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &profileCookie)
	http.SetCookie(w, &refreshCookie)
}
