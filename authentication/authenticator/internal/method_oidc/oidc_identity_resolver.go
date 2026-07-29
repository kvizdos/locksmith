package method_oidc

import (
	"context"
	"crypto/subtle"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

func (r *oidcValidationSession) ResolveIdentity(ctx context.Context) error {
	verifyCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var rawIDToken string

	switch r.flow {
	case flowCode:
		token, err := r.options.OauthConfig.Exchange(verifyCtx, r.untrustedParsedCode, oauth2.SetAuthURLParam("code_verifier", r.pkceVerifier))
		if err != nil {
			return fmt.Errorf("failed to exchange oidc code: %w", err)
		}
		var ok bool
		rawIDToken, ok = token.Extra("id_token").(string)
		if !ok {
			return fmt.Errorf("missing id_token")
		}
	case flowCredential:
		rawIDToken = r.untrustedCredentialToken
	default:
		return fmt.Errorf("unsupported flow")
	}

	idToken, err := r.options.Verifier.Verify(verifyCtx, rawIDToken)
	if err != nil {
		return fmt.Errorf("failed to verify id_token: %w", err)
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Nonce         string `json:"nonce"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return fmt.Errorf("failed to parse id_token claims: %w", err)
	}

	// The FedCM/One Tap "credential" flow is a bare bearer JWT with no PKCE
	// equivalent, so it must be bound to the browser that requested it via
	// the nonce cookie set alongside google_fcm.js (see
	// oauth.GoogleFCMNonceCookie). Without this check, an attacker could
	// obtain a valid credential for their own account and replay it against
	// a victim's browser via a forged cross-site form POST, logging the
	// victim into the attacker's account (login CSRF).
	if r.flow == flowCredential {
		if r.expectedNonce == "" || claims.Nonce == "" ||
			subtle.ConstantTimeCompare([]byte(claims.Nonce), []byte(r.expectedNonce)) != 1 {
			return fmt.Errorf("id_token nonce does not match expected sign-in session")
		}
	}

	r.SetSubject(idToken.Subject)
	r.SetIssuer(idToken.Issuer)
	r.displayName = claims.Name
	if claims.EmailVerified {
		r.SetEmail(claims.Email)
	}
	r.authoritativeToken = idToken
	return nil
}
