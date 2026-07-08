package registrationhints

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication/signing"
)

const defaultTTL = 60 * time.Second

type Service struct {
	Signer signing.SigningPackageInterface
	TTL    time.Duration
}

func (s Service) Create(hint Hint) (string, error) {
	if s.Signer == nil {
		return "", ErrMissingSigner
	}
	if !hint.Rosterable {
		return "", ErrNotRosterable
	}
	if hint.ProviderName == "" || hint.Issuer == "" || hint.Subject == "" {
		return "", ErrMissingIdentity
	}

	now := time.Now().UTC()
	ttl := s.TTL
	if ttl <= 0 {
		ttl = defaultTTL
	}

	if hint.ID == "" {
		hint.ID = uuid.NewString()
	}
	hint.IssuedAt = jwt.NewNumericDate(now)
	hint.ExpiresAt = jwt.NewNumericDate(now.Add(ttl))
	hint.Audience = jwt.ClaimStrings{Audience}

	return s.Signer.CreateJWT(&hint)
}

func (s Service) Parse(tokenString string) (*Hint, error) {
	if s.Signer == nil {
		return nil, ErrMissingSigner
	}
	if tokenString == "" {
		return nil, ErrMissingToken
	}

	hint := &Hint{}
	token, err := s.Signer.ParseJWT(tokenString, hint, jwt.WithAudience(Audience), jwt.WithExpirationRequired())
	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if hint.ExpiresAt == nil || !hint.ExpiresAt.After(time.Now().UTC()) {
		return nil, ErrInvalidToken
	}
	if !hasAudience(hint.Audience, Audience) {
		return nil, ErrInvalidToken
	}
	if !hint.Rosterable {
		return nil, ErrNotRosterable
	}
	if hint.ProviderName == "" || hint.Issuer == "" || hint.Subject == "" {
		return nil, ErrMissingIdentity
	}

	return hint, nil
}

func (s Service) FromRequest(r *http.Request) (*Hint, error) {
	if r == nil {
		return nil, ErrMissingToken
	}

	if header := r.Header.Get("X-Registration-Hint"); header != "" {
		return s.Parse(header)
	}

	cookie, err := r.Cookie(CookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, ErrMissingToken
		}
		return nil, ErrInvalidToken
	}

	return s.Parse(cookie.Value)
}

func SetCookie(w http.ResponseWriter, token string) {
	SetCookieWithTTL(w, token, defaultTTL)
}

func SetCookieWithTTL(w http.ResponseWriter, token string, ttl time.Duration) {
	if ttl <= 0 {
		ttl = defaultTTL
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/register",
		MaxAge:   int(ttl.Seconds()),
		Expires:  time.Now().UTC().Add(ttl),
	})
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/register",
	})
}

func hasAudience(audiences jwt.ClaimStrings, expected string) bool {
	for _, audience := range audiences {
		if audience == expected {
			return true
		}
	}
	return false
}
