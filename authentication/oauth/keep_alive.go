package oauth

import (
	_ "embed"
	"net/http"
)

//go:embed keep_alive.js
var keepAliveJSLBytes []byte

type KeepAliveJSRoute struct{}

func (k KeepAliveJSRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(keepAliveJSLBytes)
}
