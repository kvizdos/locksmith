package saml_init

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"strings"

	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
)

func LoadServiceProviderFromMetadata(nickname string, xmlBytes []byte) (*saml_entities.SAMLProvider, error) {
	var md saml_entities.EntityDescriptor
	if err := xml.Unmarshal(xmlBytes, &md); err != nil {
		return nil, err
	}

	if md.EntityID == "" {
		return nil, errors.New("missing SP entityID")
	}

	// --- pick ACS (prefer HTTP-POST, default=true) ---
	var acs *saml_entities.ACSService
	for _, a := range md.SP.ACS {
		if a.Binding == "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" && a.IsDefault {
			acs = &a
			break
		}
	}
	if acs == nil {
		for _, a := range md.SP.ACS {
			if a.Binding == "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" {
				acs = &a
				break
			}
		}
	}
	if acs == nil {
		return nil, errors.New("no HTTP-POST ACS found")
	}

	// --- NameIDFormat ---
	nameIDFormat := "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	if len(md.SP.NameIDFormats) > 0 {
		nameIDFormat = strings.TrimSpace(md.SP.NameIDFormats[0])
	}

	// --- SP signing cert (optional) ---
	var signingCert *string
	for _, kd := range md.SP.KeyDescriptors {
		if kd.Use == "signing" && len(kd.KeyInfo.X509Data.Certs) > 0 {
			cert := pemWrapCert(kd.KeyInfo.X509Data.Certs[0])
			signingCert = &cert
			break
		}
	}

	return &saml_entities.SAMLProvider{
		Nickname:             nickname,
		EntityID:             md.EntityID,
		ACSURL:               acs.Location,
		Binding:              acs.Binding,
		SigningCertPEM:       signingCert,
		NameIDFormat:         nameIDFormat,
		WantAssertionsSigned: md.SP.WantAssertionsSigned,
		Enabled:              true,
	}, nil
}

func pemWrapCert(b64 string) string {
	b64 = strings.ReplaceAll(b64, "\n", "")
	decoded, _ := base64.StdEncoding.DecodeString(b64)

	return "-----BEGIN CERTIFICATE-----\n" +
		base64.StdEncoding.EncodeToString(decoded) +
		"\n-----END CERTIFICATE-----"
}
