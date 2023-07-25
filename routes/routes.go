package routes

import (
	"net/http"

	"github.com/kvizdos/locksmith/administration"
	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/login"
	"github.com/kvizdos/locksmith/authentication/register"
	"github.com/kvizdos/locksmith/components"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/httpHelpers"
	"github.com/kvizdos/locksmith/launchpad"
	"github.com/kvizdos/locksmith/pages"
)

type LocksmithRoutesOptions struct {
	AppName                   string
	DisableUI                 bool
	DisableAPI                bool
	DisableComponents         bool
	DisableInvites            bool
	DisablePublicRegistration bool
	DisableLocksmithPage      bool
	LaunchpadSettings         launchpad.LocksmithLaunchpadOptions
	Styling                   pages.LocksmithPageStyling
}

func InitializeLocksmithRoutes(mux *http.ServeMux, db database.DatabaseAccessor, options LocksmithRoutesOptions) {
	if !options.DisableComponents {
		mux.HandleFunc("/components/", components.ServeComponents)
	}

	InitializeLaunchpad(mux, db, options)

	if !options.DisableAPI {
		registrationAPIHandler := httpHelpers.InjectDatabaseIntoContext(register.RegistrationHandler{
			DefaultRoleName:           "user",
			DisablePublicRegistration: options.DisablePublicRegistration,
		}, db)
		mux.Handle("/api/register", registrationAPIHandler)

		loginAPIHandler := httpHelpers.InjectDatabaseIntoContext(login.LoginHandler{}, db)
		mux.Handle("/api/login", loginAPIHandler)

		listUsersAdminAPIHandler := endpoints.SecureEndpointHTTPMiddleware(administration.AdministrationListUsersHandler{}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions: []string{"users.list.all"},
		})
		mux.Handle("/api/users/list", listUsersAdminAPIHandler)

		deleteUserAdminAPIHandler := endpoints.SecureEndpointHTTPMiddleware(administration.AdministrationDeleteUsersHandler{}, db)
		mux.Handle("/api/users/delete", deleteUserAdminAPIHandler)

		if !options.DisableInvites {
			inviteUserAPIHandler := endpoints.SecureEndpointHTTPMiddleware(invitations.AdministrationInviteUserHandler{}, db, endpoints.EndpointSecurityOptions{
				MinimalPermissions: []string{"user.invite"},
			})
			mux.Handle("/api/users/invite", inviteUserAPIHandler)
		}
	}

	if !options.DisableUI {
		mux.Handle("/login", login.LoginPageHandler{
			AppName: options.AppName,
			Styling: options.Styling,
		})
		mux.Handle("/register", httpHelpers.InjectDatabaseIntoContext(register.RegistrationPageHandler{
			AppName:                   options.AppName,
			DisablePublicRegistration: options.DisablePublicRegistration,
			Styling:                   options.Styling,
		}, db))
	}

	if !options.DisableLocksmithPage {
		serveAdminPage := endpoints.SecureEndpointHTTPMiddleware(administration.ServeAdminPage{}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions: []string{"view.ls-admin"},
		})
		mux.Handle("/locksmith", serveAdminPage)
	}
}
