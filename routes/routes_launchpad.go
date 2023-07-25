//go:build enable_launchpad
// +build enable_launchpad

package routes

import (
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/launchpad"
)

func InitializeLaunchpad(mux *http.ServeMux, db database.DatabaseAccessor, options LocksmithRoutesOptions) {
	fmt.Println("configuring launchpad..")
	if options.LaunchpadSettings.Enabled {
		launchpad.BootstrapUsers(db, options.LaunchpadSettings.AccessToken, options.LaunchpadSettings.Users)
		options.LaunchpadSettings.BootstrapDatabase(db)

		launchpadHandler := endpoints.SecureEndpointHTTPMiddleware(launchpad.LaunchpadHTTPHandler{
			AppName:        options.AppName,
			Styling:        options.Styling,
			AccessToken:    options.LaunchpadSettings.AccessToken,
			AvailableUsers: options.LaunchpadSettings.Users,
			Subtitle: options.LaunchpadSettings.Caption,
		}, db, endpoints.EndpointSecurityOptions{
			BasicAuth: endpoints.EndpointSecurityBasicAuth{
				Enabled:  true,
				Username: "launchpad",
				Password: options.LaunchpadSettings.AccessToken,
			},
		})

		mux.Handle("/launchpad", launchpadHandler)
	}
}
