package packets

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication/signing"
)

type RegistrationJWTServiceInterface interface {
	Create(oauthSourceName string, subject string) (string, error)
	Parse(tokenString string) (*RegistrationJWT, error)
	Verify(tokenString string) error

	ConfirmAutoLogin(tokenString string) (*RegistrationAutoLoginToken, error)
	CreateAutoLoginToken(oauthSourceName string, subject string) (string, error)
}

type RegistrationJWT struct {
	OAuthSourceName string `json:"oauth_source_name"`
	jwt.RegisteredClaims
}

type RegistrationJWTService struct {
	signer       signing.SigningPackageInterface
	issuer       string
	ttl          time.Duration
	autoLoginSvc *RegistrationAutoLoginJWTService
}

func NewRegistrationJWTService(
	signer signing.SigningPackageInterface,
	issuer string,
	ttl time.Duration,
	autoLoginSvc *RegistrationAutoLoginJWTService,
) *RegistrationJWTService {
	return &RegistrationJWTService{
		signer:       signer,
		issuer:       issuer,
		ttl:          ttl,
		autoLoginSvc: autoLoginSvc,
	}
}

func (s *RegistrationJWTService) ConfirmAutoLogin(tokenString string) (*RegistrationAutoLoginToken, error) {
	token, err := s.autoLoginSvc.Parse(tokenString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return token, nil
}

func (s *RegistrationJWTService) CreateAutoLoginToken(oauthSourceName string, subject string) (string, error) {
	return s.autoLoginSvc.Create(subject, oauthSourceName)
}

func (s *RegistrationJWTService) Create(oauthSourceName string, subject string) (string, error) {
	if oauthSourceName == "" {
		return "", errors.New("oauth source name is required")
	}

	if s.signer == nil {
		return "", errors.New("signer is nil")
	}

	now := time.Now()

	claims := RegistrationJWT{
		OAuthSourceName: oauthSourceName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			Audience:  jwt.ClaimStrings{"trusted-registration"},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	return s.signer.CreateJWT(claims)
}

func (s *RegistrationJWTService) Parse(tokenString string) (*RegistrationJWT, error) {
	if s.signer == nil {
		return nil, errors.New("signer is nil")
	}

	claims := &RegistrationJWT{}

	token, err := s.signer.ParseJWT(
		tokenString,
		claims,
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithAudience("trusted-registration"),
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid registration token")
	}

	if claims.OAuthSourceName == "" {
		return nil, errors.New("oauth source name is missing")
	}

	return claims, nil
}

func (s *RegistrationJWTService) Verify(tokenString string) error {
	_, err := s.Parse(tokenString)
	return err
}
