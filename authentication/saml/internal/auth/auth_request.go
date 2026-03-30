package saml_auth

import (
	"time"

	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
)

type AuthnRequest struct {
	ID           string
	Issuer       string
	ACSURL       string
	IssueInstant time.Time
	IsSigned     bool

	// raw XML if you want to verify signature later
	RawXML []byte
}

type ValidatedRequest struct {
	SP        *saml_entities.SAMLProvider
	RequestID string
	ACSURL    string
}
