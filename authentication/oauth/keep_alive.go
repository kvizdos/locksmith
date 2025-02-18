package oauth

import (
	_ "embed"
	"net/http"
)

//go:embed keep_alive.html
var keepAliveHTMLBytes []byte

type KeepAliveRoute struct{}

func (k KeepAliveRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("ls_oauth_provider"); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Write(keepAliveHTMLBytes)
}
