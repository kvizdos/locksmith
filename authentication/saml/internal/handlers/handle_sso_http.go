package saml_handlers

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"

	saml_idp "github.com/kvizdos/locksmith/authentication/saml/internal/idp"
	"github.com/kvizdos/locksmith/authentication/saml/saml_config"
	"github.com/kvizdos/locksmith/authentication/saml/saml_discovery"
)

func HandleSSO(cfg *saml_config.IdPConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Pull SAML context from middleware
		ctxVal := r.Context().Value(SAMLCtxKey{})
		if ctxVal == nil {
			http.Error(w, "missing SAML context", http.StatusInternalServerError)
			return
		}

		samlCtx := ctxVal.(*SAMLContext)

		discoveryInfo, err := cfg.Discover(r, samlCtx.Validated.SP)

		if err != nil {
			var discoveryErr saml_discovery.DiscoveryError
			if errors.As(err, &discoveryErr) {
				http.Redirect(w, r, fmt.Sprintf("/err?code=%s", discoveryErr.ErrorCode), http.StatusTemporaryRedirect)
				return
			} else {
				fmt.Printf("Failed to discover: %s\n", err.Error())
				http.Error(w, "failed to discover", http.StatusInternalServerError)
			}
			return
		}

		// 3. Build + sign SAML Response
		samlResp, err := BuildAndSignSAMLResponse(
			cfg,
			samlCtx.Validated.SP,
			samlCtx.AuthnRequest,
			&saml_idp.User{
				Email: discoveryInfo.GetEmail(),
				ID:    discoveryInfo.GetUserID(),
			},
		)
		if err != nil {
			http.Error(w, "failed to build SAML response", http.StatusInternalServerError)
			return
		}

		// 4. POST back to SP ACS
		renderAutoPOST(
			w,
			samlCtx.Validated.ACSURL,
			samlResp,
			samlCtx.RelayState,
		)
	}
}

func renderAutoPOST(
	w http.ResponseWriter,
	acsURL string,
	samlResponse string,
	relayState string,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	io.WriteString(w, `<!DOCTYPE html>
<html>
  <body onload="document.forms[0].submit()">
    <form method="POST" action="`+html.EscapeString(acsURL)+`">
      <input type="hidden" name="SAMLResponse" value="`+html.EscapeString(samlResponse)+`"/>`)

	if relayState != "" {
		io.WriteString(w, `
      <input type="hidden" name="RelayState" value="`+html.EscapeString(relayState)+`"/>`)
	}

	io.WriteString(w, `
      <noscript>
        <button type="submit">Continue</button>
      </noscript>
    </form>
  </body>
</html>`)
}
