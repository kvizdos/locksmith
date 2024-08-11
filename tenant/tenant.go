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

	HasEntitlement(entitlementName string) bool
	GetEntitlementInfo(entitlementName string) (*entitlements.TenantEntitlement, error)
	GetAttachedEntitlementIDs() []string

	ConfirmUserEntitlements(foundEntitlementIDs []string) []string

	// ToPublic() is passed as a JWT
	// for frontend rendering.
	//
	// **ONLY PASS SAFE VALUES TO THE FRONTEND**
	ToPublic() map[string]interface{}
}

type BaseTenant struct {
	ID           uuid.UUID                       `json:"uuid"`
	Entitlements entitlements.TenantEntitlements `json:"entitlements"`
}

func (b BaseTenant) ConfirmUserEntitlements(foundEntitlementIDs []string) []string {
	out := []string{}

	for _, entitlementID := range foundEntitlementIDs {
		if b.HasEntitlement(entitlementID) {
			out = append(out, entitlementID)
		}
	}

	return out
}

func (b BaseTenant) GetAttachedEntitlementIDs() []string {
	return b.Entitlements.GetIDs()
}

func (b BaseTenant) HasEntitlement(entitlementID string) bool {
	return b.Entitlements.Has(entitlementID)
}

func (b BaseTenant) GetEntitlementInfo(entitlementID string) (*entitlements.TenantEntitlement, error) {
	for _, entitlement := range b.Entitlements {
		if entitlement.ID == entitlementID {
			if !entitlement.IsActive() {
				return nil, fmt.Errorf("not activated")
			}
			return &entitlement, nil
		}
	}

	return nil, fmt.Errorf("not redeemed")
}

func (b BaseTenant) GetTenantID() uuid.UUID {
	return b.ID
}

func (b BaseTenant) ToPublic() map[string]interface{} {
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
		"uuid":         b.ID.String(),
		"entitlements": b.Entitlements.ToMap(),
	}
}

func (b BaseTenant) FromMap(input map[string]interface{}) Tenant {
	return BaseTenant{
		ID:           uuid.MustParse(input["uuid"].(string)),
		Entitlements: entitlements.TenantEntitlementsFromMap(input["entitlements"].([]interface{})),
	}
}
