package authenticator_methods

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OIDConfig struct {
	Issuer       string
	BaseURL      string
	ProviderName string
	ClientID     string
	ClientSecret string
	Rosterable   bool
	LogoBytes    []byte
}
type OIDCValidatorOptions struct {
	ProviderName string
	Rosterable   bool
	Verifier     *oidc.IDTokenVerifier
	OauthConfig  *oauth2.Config
	LogoBytes    []byte
}

type OIDCValidatorOption func(*OIDCValidatorOptions)

func WithOIDC(c OIDConfig) OIDCValidatorOption {
	provider, err := oidc.NewProvider(context.Background(), c.Issuer)
	if err != nil {
		panic(fmt.Errorf("Failed to register provider: %w", err))
	}

	callbackURL := fmt.Sprintf("%s/api/login", c.BaseURL)

	config := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  callbackURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: c.ClientID})

	return func(opts *OIDCValidatorOptions) {
		opts.ProviderName = c.ProviderName
		opts.Rosterable = c.Rosterable
		opts.Verifier = verifier
		opts.OauthConfig = &config
		opts.LogoBytes = c.LogoBytes
	}
}
