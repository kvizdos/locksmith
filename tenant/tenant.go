package tenant

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/entitlements"
)

type Tenant interface {
	GetTenantID() uuid.UUID

	ToMap() map[string]interface{} // tenant ID must be saved as "uuid" in the map
	FromMap(input map[string]interface{}) Tenant

	HasEntitlement(entitlementName string) (entitlements.Entitlement, error)
	GetAttachedEntitlements() []string

	// ToPublic() is passed as a JWT
	// for frontend rendering.
	//
	// **ONLY PASS SAFE VALUES TO THE FRONTEND**
	//
	// No Public values are used for validation;
	// they can be different names than the backend
	// structure.
	ToPublic(userRole string) map[string]interface{}
}

type BaseTenant struct {
	ID           uuid.UUID `json:"uuid"`
	Entitlements []string  `json:"entitlements"`
}

func (b BaseTenant) GetAttachedEntitlements() []string {
	return b.Entitlements
}

func (b BaseTenant) HasEntitlement(entitlementName string) (entitlements.Entitlement, error) {
	for _, name := range b.Entitlements {
		if name == entitlementName {
			return entitlements.GetEntitlement(name), nil
		}
	}

	return entitlements.Entitlement{}, fmt.Errorf("entitlement not attached")
}

func (b BaseTenant) GetTenantID() uuid.UUID {
	return b.ID
}

func (b BaseTenant) ToPublic(userRole string) map[string]interface{} {
	// Optionally, return different tenant
	// information depending on the
	// user role.
	//
	// Regular users may not need all of the
	// same info as administrators.
	return map[string]interface{}{}
}

func (b BaseTenant) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"uuid": b.ID.String(),
	}
}

func (b BaseTenant) FromMap(input map[string]interface{}) Tenant {
	return BaseTenant{
		ID: uuid.MustParse(input["uuid"].(string)),
	}
}
