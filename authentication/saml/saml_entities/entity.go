package saml_entities

type EntityDescriptor struct {
	EntityID string `xml:"entityID,attr"`
	SP       spSSO  `xml:"SPSSODescriptor"`
}

type spSSO struct {
	WantAssertionsSigned bool `xml:"WantAssertionsSigned,attr"`

	NameIDFormats []string `xml:"NameIDFormat"`

	KeyDescriptors []keyDescriptor `xml:"KeyDescriptor"`
	ACS            []ACSService    `xml:"AssertionConsumerService"`
}

type keyDescriptor struct {
	Use     string  `xml:"use,attr"`
	KeyInfo keyInfo `xml:"KeyInfo"`
}

type keyInfo struct {
	X509Data x509Data `xml:"X509Data"`
}

type x509Data struct {
	Certs []string `xml:"X509Certificate"`
}

type ACSService struct {
	Location  string `xml:"Location,attr"`
	Binding   string `xml:"Binding,attr"`
	IsDefault bool   `xml:"isDefault,attr"`
}
