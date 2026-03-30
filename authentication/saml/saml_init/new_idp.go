package saml_init

import (
	"github.com/kvizdos/locksmith/authentication/saml/saml_config"
	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
)

func NewIdPConfig(entityID, ssoURL, signingCertPEM string, nameIDFormats []string, supportEmail, organization string, enabledProviders ...*saml_entities.SAMLProvider) *saml_config.IdPConfig {
	return &saml_config.IdPConfig{
		EntityID:         entityID,
		SSOURL:           ssoURL,
		SigningCertPEM:   signingCertPEM,
		NameIDFormats:    nameIDFormats,
		SupportEmail:     supportEmail,
		Organization:     organization,
		EnabledProviders: enabledProviders,
	}
}
