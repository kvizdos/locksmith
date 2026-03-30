package saml_handlers

import (
	"crypto"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/google/uuid"
	saml_auth "github.com/kvizdos/locksmith/authentication/saml/internal/auth"
	saml_idp "github.com/kvizdos/locksmith/authentication/saml/internal/idp"
	"github.com/kvizdos/locksmith/authentication/saml/saml_config"
	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
	dsig "github.com/russellhaering/goxmldsig"
)

func BuildAndSignSAMLResponse(
	cfg *saml_config.IdPConfig,
	sp *saml_entities.SAMLProvider,
	req *saml_auth.AuthnRequest,
	user *saml_idp.User,
) (string, error) {
	if cfg == nil || cfg.Signer == nil || cfg.Signer.Key == nil || cfg.Signer.Cert == nil || cfg.Signer.EntityID == "" {
		return "", errors.New("IdP signer not configured")
	}
	if sp == nil || req == nil || user == nil {
		return "", errors.New("nil input")
	}
	if sp.ACSURL == "" {
		return "", errors.New("missing SP ACSURL")
	}
	if sp.EntityID == "" {
		return "", errors.New("missing SP EntityID")
	}

	now := time.Now().UTC()
	issueInstant := now.Format(time.RFC3339)
	notBefore := now.Add(-60 * time.Second).Format(time.RFC3339)
	notOnOrAfter := now.Add(5 * time.Minute).Format(time.RFC3339)

	// IDs
	respID := "_" + uuid.NewString()
	assertID := "_" + uuid.NewString()
	sessionIndex := "_" + uuid.NewString()

	// Choose NameID value (be consistent!)
	// If you want opaque stable ID: use user.ID
	// If you want email: use user.Email
	nameIDValue := strings.TrimSpace(user.Email)
	if nameIDValue == "" {
		nameIDValue = strings.TrimSpace(user.ID)
	}
	if nameIDValue == "" {
		return "", errors.New("user has no stable identifier for NameID")
	}

	nameIDFormat := strings.TrimSpace(sp.NameIDFormat)
	if nameIDFormat == "" {
		nameIDFormat = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	}

	// Build UNSIGNED assertion first
	ass := samlAssertion{
		XmlnsSaml:    "urn:oasis:names:tc:SAML:2.0:assertion",
		ID:           assertID,
		Version:      "2.0",
		IssueInstant: issueInstant,
		Issuer:       samlIssuer{Value: cfg.EntityID},
		Subject: samlSubject{
			NameID: samlNameID{
				Format: nameIDFormat,
				Value:  nameIDValue,
			},
			SubjectConfirmation: samlSubjectConfirmation{
				Method: "urn:oasis:names:tc:SAML:2.0:cm:bearer",
				Data: samlSubjectConfirmationData{
					InResponseTo: req.ID,
					Recipient:    sp.ACSURL,
					NotOnOrAfter: notOnOrAfter,
				},
			},
		},
		Conditions: samlConditions{
			NotBefore:    notBefore,
			NotOnOrAfter: notOnOrAfter,
			AudienceRestriction: samlAudienceRestriction{
				Audience: sp.EntityID,
			},
		},
		AuthnStmt: samlAuthnStmt{
			AuthnInstant: issueInstant,
			SessionIndex: sessionIndex,
			AuthnContext: samlAuthnContext{
				// Common safe default
				ClassRef: "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport",
			},
		},
		// Optional attributes (SP might ignore, harmless)
		AttrStmtOpt: &samlAttrStmt{
			Attributes: []samlAttribute{
				{Name: "email", Values: []samlAttrValue{{Value: user.Email}}},
				{Name: "user_id", Values: []samlAttrValue{{Value: user.ID}}},
			},
		},
	}

	unsignedAssertionXML, err := xml.Marshal(ass)
	if err != nil {
		return "", err
	}

	// Sign the assertion (enveloped signature)
	ctx := dsig.NewDefaultSigningContext(cfg.Signer)
	ctx.Hash = crypto.SHA256
	ctx.Canonicalizer = dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(unsignedAssertionXML); err != nil {
		return "", err
	}

	assertionEl := doc.Root()
	if assertionEl == nil {
		return "", errors.New("missing assertion root element")
	}

	signedEl, err := ctx.SignEnveloped(assertionEl)
	if err != nil {
		return "", err
	}

	// serialize signed assertion element → []byte
	signedDoc := etree.NewDocument()
	signedDoc.SetRoot(signedEl)

	// optional but nice
	signedDoc.WriteSettings = etree.WriteSettings{
		CanonicalEndTags: true,
	}

	signedAssertionXML, err := signedDoc.WriteToBytes()
	if err != nil {
		return "", err
	}

	// Build outer Response with signed assertion injected
	resp := samlpResponse{
		XmlnsSamlp:   "urn:oasis:names:tc:SAML:2.0:protocol",
		XmlnsSaml:    "urn:oasis:names:tc:SAML:2.0:assertion",
		ID:           respID,
		Version:      "2.0",
		IssueInstant: issueInstant,
		Destination:  sp.ACSURL,
		InResponseTo: req.ID,
		Issuer:       samlIssuer{Value: cfg.EntityID},
		Status: samlpStatus{
			StatusCode: samlpStatusCode{
				Value: "urn:oasis:names:tc:SAML:2.0:status:Success",
			},
		},
		Assertion: signedAssertionXML,
	}
	respXML, err := xml.Marshal(resp)
	if err != nil {
		return "", err
	}

	respDoc := etree.NewDocument()
	if err := respDoc.ReadFromBytes(respXML); err != nil {
		return "", err
	}

	respEl := respDoc.Root()
	if respEl == nil {
		return "", errors.New("missing Response root element")
	}

	respCtx := dsig.NewDefaultSigningContext(cfg.Signer)
	respCtx.Hash = crypto.SHA256
	respCtx.IdAttribute = "ID"
	respCtx.Canonicalizer = dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")

	signedRespEl, err := respCtx.SignEnveloped(respEl)
	if err != nil {
		return "", err
	}

	finalDoc := etree.NewDocument()
	finalDoc.SetRoot(signedRespEl)

	finalBytes, err := finalDoc.WriteToBytes()
	if err != nil {
		return "", err
	}

	final := append([]byte(xml.Header), finalBytes...)

	// Return base64 for HTML form field
	return base64.StdEncoding.EncodeToString(final), nil
}

type samlpResponse struct {
	XMLName      xml.Name `xml:"samlp:Response"`
	XmlnsSamlp   string   `xml:"xmlns:samlp,attr"`
	XmlnsSaml    string   `xml:"xmlns:saml,attr"`
	ID           string   `xml:"ID,attr"`
	Version      string   `xml:"Version,attr"`
	IssueInstant string   `xml:"IssueInstant,attr"`
	Destination  string   `xml:"Destination,attr,omitempty"`
	InResponseTo string   `xml:"InResponseTo,attr,omitempty"`

	Issuer samlIssuer `xml:"saml:Issuer"`

	Status samlpStatus `xml:"samlp:Status"`

	Assertion []byte `xml:",innerxml"`
}

type samlpStatus struct {
	StatusCode samlpStatusCode `xml:"samlp:StatusCode"`
}

type samlpStatusCode struct {
	Value string `xml:"Value,attr"`
}

type samlIssuer struct {
	Value string `xml:",chardata"`
}

type samlAssertion struct {
	XMLName      xml.Name `xml:"saml:Assertion"`
	XmlnsSaml    string   `xml:"xmlns:saml,attr"`
	ID           string   `xml:"ID,attr"`
	Version      string   `xml:"Version,attr"`
	IssueInstant string   `xml:"IssueInstant,attr"`

	Issuer samlIssuer `xml:"saml:Issuer"`

	Subject     samlSubject    `xml:"saml:Subject"`
	Conditions  samlConditions `xml:"saml:Conditions"`
	AuthnStmt   samlAuthnStmt  `xml:"saml:AuthnStatement"`
	AttrStmtOpt *samlAttrStmt  `xml:"saml:AttributeStatement,omitempty"`
}

type samlSubject struct {
	NameID              samlNameID              `xml:"saml:NameID"`
	SubjectConfirmation samlSubjectConfirmation `xml:"saml:SubjectConfirmation"`
}

type samlNameID struct {
	Format string `xml:"Format,attr,omitempty"`
	Value  string `xml:",chardata"`
}

type samlSubjectConfirmation struct {
	Method string                      `xml:"Method,attr"`
	Data   samlSubjectConfirmationData `xml:"saml:SubjectConfirmationData"`
}

type samlSubjectConfirmationData struct {
	InResponseTo string `xml:"InResponseTo,attr,omitempty"`
	Recipient    string `xml:"Recipient,attr,omitempty"`
	NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
}

type samlConditions struct {
	NotBefore           string                  `xml:"NotBefore,attr"`
	NotOnOrAfter        string                  `xml:"NotOnOrAfter,attr"`
	AudienceRestriction samlAudienceRestriction `xml:"saml:AudienceRestriction"`
}

type samlAudienceRestriction struct {
	Audience string `xml:"saml:Audience"`
}

type samlAuthnStmt struct {
	AuthnInstant string           `xml:"AuthnInstant,attr"`
	SessionIndex string           `xml:"SessionIndex,attr,omitempty"`
	AuthnContext samlAuthnContext `xml:"saml:AuthnContext"`
}

type samlAuthnContext struct {
	ClassRef string `xml:"saml:AuthnContextClassRef"`
}

type samlAttrStmt struct {
	Attributes []samlAttribute `xml:"saml:Attribute"`
}

type samlAttribute struct {
	Name   string          `xml:"Name,attr"`
	Values []samlAttrValue `xml:"saml:AttributeValue"`
}

type samlAttrValue struct {
	Value string `xml:",chardata"`
}
