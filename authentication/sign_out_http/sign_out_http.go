package sign_out

import (
	"net/http"
)

type SignOutHTTP struct{}

func (m SignOutHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Clear-Site-Data", `"cookies", "storage"`)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
