package method_oidc

import (
	"encoding/json"
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
		pv.untrustedParsedCode = r.URL.Query().Get("code")
		return nil
	case flowCredential:
		var body struct {
			Credential string `json:"credential"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return fmt.Errorf("failed to decode credential body: %w", err)
		}
		if body.Credential == "" {
			return errors.New("missing credential in request body")
		}
		pv.untrustedCredentialToken = body.Credential
		return nil
	}

	return errors.New("unsupported flow")

}
