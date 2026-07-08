package method_oidc

import (
	"context"
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
	}
	if err := idToken.Claims(&claims); err != nil {
		return fmt.Errorf("failed to parse id_token claims: %w", err)
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
