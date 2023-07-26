//go:build !enable_launchpad
// +build !enable_launchpad

package launchpad

import (
	"net/http"
)

func LaunchpadRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
