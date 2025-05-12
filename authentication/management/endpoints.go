package management

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/database"
)

func RouteManagementAPI(mux *http.ServeMux, db database.DatabaseAccessor) {
	// Create a sub-mux (subrouter) for API routes
	apiMux := http.NewServeMux()

	apiMux.Handle("/me", endpoints.SecureEndpointHTTPMiddleware(meEndpointHTTP{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{
			"human",
		},
	}))

	mux.Handle("/api/management/", http.StripPrefix("/api/management", apiMux))
}
