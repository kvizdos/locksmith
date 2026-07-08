package method_oidc

import (
	"errors"
	"fmt"
	"net/http"
)

func (pv *oidcValidationSession) LoadRequest(r *http.Request) error {
	// Method check unnecessary; it is part of detect flow.
	oidcType := detectFlow(r)

	pv.flow = oidcType

	switch oidcType {
	case flowCode:
		pkceCookie, err := r.Cookie("ls_oidc_pkce")
		if err != nil {
			return fmt.Errorf("missing pkce verifier cookie: %w", err)
		}
		pv.pkceVerifier = pkceCookie.Value
		pv.selectBy = "generic_oidc_btn"
		pv.untrustedParsedCode = r.URL.Query().Get("code")
		// The OAuth2 "state" param carries the "return to this page after
		// login" target set by Begin; the provider echoes it back verbatim.
		pv.redirectTarget = sanitizeRedirectPath(r.URL.Query().Get("state"))
		return nil
	case flowCredential:
		// Google Identity Services' HTML API (data-login_uri) delivers the
		// credential via a native browser form POST, encoded as
		// application/x-www-form-urlencoded — not JSON. It carries a
		// "credential" field (the ID token) and a "select_by" field
		// indicating how the credential was produced (e.g. "btn",
		// "btn_confirm" for a rendered button click vs "auto", "fedcm",
		// "fedcm_auto" for an automatic One Tap/FedCM sign-in).
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse credential form body: %w", err)
		}
		credential := r.PostFormValue("credential")
		if credential == "" {
			return errors.New("missing credential in request body")
		}
		pv.untrustedCredentialToken = credential
		pv.selectBy = r.PostFormValue("select_by")
		if pv.selectBy == "" {
			pv.selectBy = "flow_credential"
		}
		// Same "state" convention as the code flow: the caller tells us where
		// to redirect after login via a "state" form field.
		pv.redirectTarget = sanitizeRedirectPath(r.PostFormValue("state"))
		return nil
	}

	return errors.New("unsupported flow")

}
