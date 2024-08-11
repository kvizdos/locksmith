//go:build enable_launchpad
// +build enable_launchpad

package routes

import (
	"fmt"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/login"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/launchpad"
)

func InitializeLaunchpad(mux *http.ServeMux, db database.DatabaseAccessor, options LocksmithRoutesOptions) {
	fmt.Println("WARNING !!! Configuring Launchpad.. DO NOT USE IN PRODUCTIN !!! WARNING")
	if options.LaunchpadSettings.Enabled {
		options.LaunchpadSettings.BootstrapDatabase(db)
		launchpad.BootstrapUsers(db, options.LaunchpadSettings.AccessToken, options.LaunchpadSettings.Users)
		launchpad.BootstrapTenants(db, options.LaunchpadSettings.AccessToken, options.LaunchpadSettings.Tenants)

		launchpadHandler := login.LoginPageMiddleware{
			Next: endpoints.SecureEndpointHTTPMiddleware(launchpad.LaunchpadHTTPHandler{
				AppName:                       options.AppName,
				Styling:                       options.Styling,
				AccessToken:                   options.LaunchpadSettings.AccessToken,
				AvailableUsers:                options.LaunchpadSettings.Users,
				Subtitle:                      options.LaunchpadSettings.Caption,
				RefreshButtonText:             options.LaunchpadSettings.RefreshButtonText,
				IsEarlyDevelopmentEnvironment: options.LaunchpadSettings.IsEarlyDevelopmentEnvironment,
			}, db, endpoints.EndpointSecurityOptions{
				BasicAuth: endpoints.EndpointSecurityBasicAuth{
					Enabled:  true,
					Username: "launchpad",
					Password: options.LaunchpadSettings.AccessToken,
				},
			}),
		}

		mux.Handle("/launchpad", launchpadHandler)

		launchpadRefresh := endpoints.SecureEndpointHTTPMiddleware(launchpad.LaunchpadRefreshHTTPHandler{
			AccessToken:       options.LaunchpadSettings.AccessToken,
			AvailableUsers:    options.LaunchpadSettings.Users,
			BootstrapDatabase: options.LaunchpadSettings.BootstrapDatabase,
		}, db, endpoints.EndpointSecurityOptions{
			BasicAuth: endpoints.EndpointSecurityBasicAuth{
				Enabled:  true,
				Username: "launchpad",
				Password: options.LaunchpadSettings.AccessToken,
			},
		})

		mux.Handle("/launchpad/refresh", launchpadRefresh)
	}
}
