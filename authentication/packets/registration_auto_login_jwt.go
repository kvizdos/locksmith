package packets

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication/signing"
)

type RegistrationAutoLoginJWTServiceInterface interface {
	Create(subject string) (string, error)
	Parse(tokenString string) (*RegistrationAutoLoginToken, error)
	Verify(tokenString string) error
}

type RegistrationAutoLoginToken struct {
	jwt.RegisteredClaims
	Source string `json:"source"` // "password", "google", "microsoft", etc.

}

type RegistrationAutoLoginJWTService struct {
	signer signing.SigningPackageInterface
	issuer string
	ttl    time.Duration
}

func NewRegistrationAutoLoginJWTService(
	signer signing.SigningPackageInterface,
	issuer string,
	ttl time.Duration,
) *RegistrationAutoLoginJWTService {
	return &RegistrationAutoLoginJWTService{
		signer: signer,
		issuer: issuer,
		ttl:    ttl,
	}
}

func (s *RegistrationAutoLoginJWTService) Create(subject string, source string) (string, error) {
	if subject == "" {
		return "", errors.New("subject is required")
	}

	if source == "" {
		return "", errors.New("source is required")
	}

	if s.signer == nil {
		return "", errors.New("signer is nil")
	}

	now := time.Now()

	claims := RegistrationAutoLoginToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings{"registration_auto_login"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
		Source: source,
	}

	return s.signer.CreateJWT(claims)
}

func (s *RegistrationAutoLoginJWTService) Parse(tokenString string) (*RegistrationAutoLoginToken, error) {
	if s.signer == nil {
		return nil, errors.New("signer is nil")
	}

	claims := &RegistrationAutoLoginToken{}

	token, err := s.signer.ParseJWT(
		tokenString,
		claims,
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience("registration_auto_login"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid registration auto-login token")
	}

	if claims.Source == "" {
		return nil, errors.New("source is missing")
	}

	if claims.Subject == "" {
		return nil, errors.New("subject is missing")
	}

	return claims, nil
}

func (s *RegistrationAutoLoginJWTService) Verify(tokenString string) error {
	_, err := s.Parse(tokenString)
	return err
}
