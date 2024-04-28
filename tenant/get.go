package tenant

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/database"
)

func GetTenantFromID(tenantID uuid.UUID, db database.DatabaseAccessor, intoTenantInterface Tenant) (Tenant, error) {
	rawTenant, found := db.FindOne("tenants", map[string]interface{}{
		"uuid": tenantID.String(),
	})

	if !found {
		return nil, fmt.Errorf("tenant ID not found!")
	}

	tenant := intoTenantInterface.FromMap(rawTenant.(map[string]interface{}))

	return tenant, nil
}
