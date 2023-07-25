//go:build !enable_launchpad
// +build !enable_launchpad

package routes

import (
	"net/http"

	"github.com/kvizdos/locksmith/database"
)

// Purposefully empty.
// Do nothing with Launchpad if
// built on a Production system.
func InitializeLaunchpad(mux *http.ServeMux, db database.DatabaseAccessor, options LocksmithRoutesOptions) {
}
