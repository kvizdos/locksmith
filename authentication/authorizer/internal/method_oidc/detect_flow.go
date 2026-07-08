package method_oidc

import "net/http"

type oidcFlow int

const (
	flowNone oidcFlow = iota
	flowCode
	flowCredential
)

func detectFlow(r *http.Request) oidcFlow {
	// Authorization Code Flow
	if r.Method == http.MethodGet &&
		r.URL.Query().Has("state") &&
		r.URL.Query().Has("code") {
		return flowCode
	}

	// FedCM / Credential Flow: browser POSTs the ID token directly
	if r.Method == http.MethodPost {
		return flowCredential
	}

	return flowNone
}
