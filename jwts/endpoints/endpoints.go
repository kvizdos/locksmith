package jwt_endpoints

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/database"
)

func RouteJWTEndpoints(mux *http.ServeMux, db database.DatabaseAccessor) {
	// Create a sub-mux (subrouter) for API routes
	apiMux := http.NewServeMux()

	apiMux.Handle("/issue", endpoints.SecureEndpointHTTPMiddleware(issueJWTHTTP{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{
			"jwt.issue",
		},
	}))

	mux.Handle("/api/jwt/", http.StripPrefix("/api/jwt", apiMux))
}
