package saml_http

import (
	"net/http"

	"github.com/kvizdos/locksmith/authentication/endpoints"
	saml_handlers "github.com/kvizdos/locksmith/authentication/saml/internal/handlers"
	"github.com/kvizdos/locksmith/authentication/saml/saml_config"
	"github.com/kvizdos/locksmith/database"
)

func ServeSAML(db database.DatabaseAccessor, cfg *saml_config.IdPConfig) *http.ServeMux {
	apiMux := http.NewServeMux()

	ssoHandler := saml_handlers.HandleSSO(cfg)

	apiMux.HandleFunc("GET /metadata.xml", saml_handlers.ServeMetadata(cfg))

	apiMux.Handle("/sso", endpoints.SecureEndpointHTTPMiddleware(handleSSORequest(cfg.EnabledProviders)(
		ssoHandler,
	), db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{},
		BaseURL:            "/api/auth/saml",
		CustomUser:         cfg.GetUserDecoder(),
	}),
	)
	return apiMux
}
