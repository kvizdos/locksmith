package login

import (
	"net/http"
)

type LoginPageMiddleware struct {
	Next http.Handler
}

func (h LoginPageMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Next.ServeHTTP(w, r)
}
