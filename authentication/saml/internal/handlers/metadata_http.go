package saml_handlers

import (
	"net/http"

	saml_idp "github.com/kvizdos/locksmith/authentication/saml/internal/idp"
	"github.com/kvizdos/locksmith/authentication/saml/saml_config"
)

func ServeMetadata(cfg *saml_config.IdPConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xmlBytes, err := saml_idp.BuildIdPMetadata(cfg)
		if err != nil {
			http.Error(w, "metadata error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/samlmetadata+xml")
		w.Write(xmlBytes)
	}
}
