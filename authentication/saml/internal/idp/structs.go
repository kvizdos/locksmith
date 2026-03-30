package saml_idp

import "encoding/xml"

type EntityDescriptor struct {
	XMLName  xml.Name         `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntityDescriptor"`
	EntityID string           `xml:"entityID,attr"`
	IDPSSO   IDPSSODescriptor `xml:"IDPSSODescriptor"`
}

type IDPSSODescriptor struct {
	ProtocolSupport string          `xml:"protocolSupportEnumeration,attr"`
	KeyDescriptor   KeyDescriptor   `xml:"KeyDescriptor"`
	SSOService      SingleSignOnSvc `xml:"SingleSignOnService"`
	NameIDFormats   []string        `xml:"NameIDFormat"`
}

type KeyDescriptor struct {
	Use     string  `xml:"use,attr"`
	KeyInfo KeyInfo `xml:"KeyInfo"`
}

type KeyInfo struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	X509Data X509Data `xml:"X509Data"`
}

type X509Data struct {
	Cert string `xml:"X509Certificate"`
}

type SingleSignOnSvc struct {
	Binding  string `xml:"Binding,attr"`
	Location string `xml:"Location,attr"`
}
