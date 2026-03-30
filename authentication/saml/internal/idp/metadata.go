package saml_idp

import (
	"encoding/xml"
	"strings"

	"github.com/kvizdos/locksmith/authentication/saml/saml_config"
)

func BuildIdPMetadata(cfg *saml_config.IdPConfig) ([]byte, error) {
	cert := stripPEM(cfg.SigningCertPEM)

	md := EntityDescriptor{
		EntityID: cfg.EntityID,
		IDPSSO: IDPSSODescriptor{
			ProtocolSupport: "urn:oasis:names:tc:SAML:2.0:protocol",
			KeyDescriptor: KeyDescriptor{
				Use: "signing",
				KeyInfo: KeyInfo{
					X509Data: X509Data{
						Cert: cert,
					},
				},
			},
			SSOService: SingleSignOnSvc{
				Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
				Location: cfg.SSOURL,
			},
			NameIDFormats: cfg.NameIDFormats,
		},
	}

	return xml.MarshalIndent(md, "", "  ")
}

func stripPEM(pem string) string {
	out := ""
	for _, line := range strings.Split(pem, "\n") {
		if strings.HasPrefix(line, "-----") {
			continue
		}
		out += strings.TrimSpace(line)
	}
	return out
}
