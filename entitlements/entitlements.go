package entitlements

import (
	"time"
)

type TenantEntitlements []TenantEntitlement

func (t TenantEntitlements) Has(entitlementID string) bool {
	for _, entitlement := range t {
		if entitlement.ID == entitlementID && entitlement.IsActive() {
			return true
		}
	}

	return false
}

func (t TenantEntitlements) ToMap() []map[string]interface{} {
	out := make([]map[string]interface{}, len(t))
	for i, e := range t {
		out[i] = e.ToMap()
	}
	return out
}

func TenantEntitlementsFromMap(input []interface{}) TenantEntitlements {
	out := make(TenantEntitlements, len(input))
	for i, e := range input {
		out[i] = TenantEntitlementFromMap(e)
	}
	return out
}

func (t TenantEntitlements) GetIDs() []string {
	out := make([]string, len(t))
	for i, entitlement := range t {
		if !entitlement.IsActive() {
			continue
		}
		out[i] = entitlement.ID
	}

	return out
}

type TenantEntitlement struct {
	ID            string `json:"id"`
	Quantity      int64  `json:"qty"`
	QuantitySpent int64  `json:"qtySpent"`

	StartDate time.Time `json:"start"`
	EndDate   time.Time `json:"end"`
}

func (t TenantEntitlement) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       t.ID,
		"qty":      t.Quantity,
		"qtySpent": t.QuantitySpent,
		"start":    t.StartDate.Unix(),
		"end":      t.EndDate.Unix(),
	}
}

func TenantEntitlementFromMap(input interface{}) TenantEntitlement {
	inp := input.(map[string]interface{})
	return TenantEntitlement{
		ID:            inp["id"].(string),
		Quantity:      inp["qty"].(int64),
		QuantitySpent: inp["qtySpent"].(int64),
		StartDate:     time.Unix(inp["start"].(int64), 0).UTC(),
		EndDate:       time.Unix(inp["end"].(int64), 0).UTC(),
	}
}

func (t TenantEntitlement) IsActive() bool {
	comp := time.Now().UTC()
	return t.StartDate.Before(comp) && t.EndDate.After(comp)
}

type Entitlement struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// map[roleName][]permissions
	AddedPermissions map[string][]string `json:"permissions"`
}

var entitlements map[string]map[string][]string

func AddEntitlement(e Entitlement) {
	if entitlements == nil {
		entitlements = map[string]map[string][]string{}
	}
	entitlements[e.Name] = e.AddedPermissions
}

func GetEntitlement(name string) Entitlement {
	if perms, ok := entitlements[name]; ok {
		return Entitlement{
			Name:             name,
			AddedPermissions: perms,
		}
	}

	// Panic, because these are hard-coded.
	panic("Entitlement name " + name + " not found!")
}
