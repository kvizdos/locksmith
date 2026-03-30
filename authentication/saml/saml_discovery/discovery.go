package saml_discovery

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
)

type DiscoveryFunc func(r *http.Request, sp *saml_entities.SAMLProvider) (*IdPDiscovery, error)

type IdPDiscovery struct {
	userID string
	email  string
}

func (i IdPDiscovery) GetEmail() string {
	return i.email
}

func (i IdPDiscovery) GetUserID() string {
	return i.userID
}

func (i *IdPDiscovery) SetEmail(email string) *IdPDiscovery {
	i.email = email
	return i
}
func (i *IdPDiscovery) SetUserID(id string) *IdPDiscovery {
	i.userID = id
	return i
}

func NewIdPDiscovery() *IdPDiscovery {
	return &IdPDiscovery{}
}
