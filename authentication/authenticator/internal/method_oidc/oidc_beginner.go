package method_oidc

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"golang.org/x/oauth2"
)

func (pv oidcHandler) Begin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		return errors.Join(
			fmt.Errorf("unsupported method: %s", r.Method),
			authenticator_domain.ErrMethodNotSupported,
		)
	}

	if pv.options.OauthConfig == nil {
		return errors.New("oauth config is nil")
	}

	log := ctx.Value("log").(*slog.Logger)

	log.DebugContext(ctx, "starting oidc authorization",
		"provider", pv.options.ProviderName,
		"method", r.Method,
	)

	state, err := randomURLSafe(32)
	if err != nil {
		return fmt.Errorf("generate state: %w", err)
	}

	verifier, err := randomURLSafe(64)
	if err != nil {
		return fmt.Errorf("generate pkce verifier: %w", err)
	}

	challenge := pkceChallenge(verifier)

	log.DebugContext(ctx, "generated oidc state and pkce challenge")

	// TODO: ideally store these server-side in your token/session store.
	// Cookie is okay if signed/encrypted elsewhere; otherwise prefer server-side.
	http.SetCookie(w, &http.Cookie{
		Name:     "ls_oidc_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((10 * time.Minute).Seconds()),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "ls_oidc_pkce",
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((10 * time.Minute).Seconds()),
	})

	log.DebugContext(ctx, "stored oidc cookies",
		"secure", r.TLS != nil,
		"same_site", http.SameSiteLaxMode,
	)

	authURL := pv.options.OauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	log.DebugContext(ctx, "redirecting to oidc provider",
		"provider", pv.options.ProviderName,
		"redirect_uri", pv.options.OauthConfig.RedirectURL,
	)

	http.Redirect(w, r, authURL, http.StatusFound)
	return nil
}

func randomURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
