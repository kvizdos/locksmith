package oauth

import (
	_ "embed"
	"net/http"
)

//go:embed keep_alive.js
var keepAliveJSLBytes []byte

//go:embed google_fcm.js
var googleFCMJSLBytes []byte

type KeepAliveJSRoute struct{}

func (k KeepAliveJSRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(keepAliveJSLBytes)
}

type GoogleFCMJSRoute struct{}

func (g GoogleFCMJSRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(googleFCMJSLBytes)
}
