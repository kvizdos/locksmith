// Package xsrf provides XSRF token generation and validation for use with
// the authenticator's session-based flows. Tokens are short-lived, signed
// JWTs bound to a session ID, issued via an injected SigningPackageInterface.
package xsrf

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication/signing"
)

const audience = "xsrf"

var (
	ErrTokenExpired    = errors.New("xsrf token expired")
	ErrTokenInvalid    = errors.New("xsrf token invalid")
	ErrSessionMismatch = errors.New("xsrf session ID mismatch")
)

// claims is the JWT payload for an XSRF token.
// SessionID binds the token to a specific browser session.
type claims struct {
	jwt.RegisteredClaims
	SessionID string `json:"sid"`
}

// Manager generates and validates XSRF tokens using a provided signing package.
type Manager struct {
	sp  signing.SigningPackageInterface
	ttl time.Duration
}

// New creates an XSRFManager with the given signing package and token TTL.
func New(sp signing.SigningPackageInterface, ttl time.Duration) Manager {
	return Manager{sp: sp, ttl: ttl}
}

// Generate creates a signed XSRF token bound to the given session ID.
func (m Manager) Generate(sessionID string) (string, error) {
	if sessionID == "" {
		return "", errors.New("sessionID must not be empty")
	}

	now := time.Now().UTC()
	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
		SessionID: sessionID,
	}

	token, err := m.sp.CreateJWT(c)
	if err != nil {
		return "", fmt.Errorf("xsrf: sign token: %w", err)
	}

	return token, nil
}

// Validate verifies the token's signature, expiry, audience, and that its
// session ID matches the one presented by the caller.
func (m Manager) Validate(token string, sessionID string) error {
	c := &claims{}
	_, err := m.sp.ParseJWT(
		token,
		c,
		jwt.WithAudience(audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrTokenExpired
		}
		return fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	}

	if c.SessionID != sessionID {
		return ErrSessionMismatch
	}

	return nil
}
