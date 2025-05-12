package service_keys

import (
	"github.com/kvizdos/locksmith/database"
)

func ListServiceKeys(db database.DatabaseAccessor, filterKey string, filterValue any) []ServiceKey {
	results := []ServiceKey{}

	var items []any
	var found bool

	if filterKey == "" {
		items, found = db.Find("service_keys", map[string]any{})
	} else {
		items, found = db.Find("service_keys", map[string]any{
			filterKey: filterValue,
		})
	}

	if !found {
		return []ServiceKey{}
	}

	for _, item := range items {
		results = append(results, ServiceKeyFromMap(item))
	}

	return results
}
