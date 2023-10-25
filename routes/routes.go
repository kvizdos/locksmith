package routes

import (
	"net/http"

	"github.com/kvizdos/locksmith/administration"
	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/login"
	"github.com/kvizdos/locksmith/authentication/register"
	"github.com/kvizdos/locksmith/authentication/reset"
	"github.com/kvizdos/locksmith/components"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/httpHelpers"
	"github.com/kvizdos/locksmith/launchpad"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/users"
)

type LocksmithRoutesOptions struct {
	AppName                   string
	DisableUI                 bool
	DisableAPI                bool
	DisableComponents         bool
	DisableInvites            bool
	DisablePublicRegistration bool
	DisableLocksmithPage      bool
	UseEmailAsUsername        bool
	OnboardPath               string
	CustomUserRegistration    register.RegisterCustomUserFunc
	LaunchpadSettings         launchpad.LocksmithLaunchpadOptions
	Styling                   pages.LocksmithPageStyling
	ResetPasswordOptions      ResetPasswordOptions
}

type ResetPasswordOptions struct {
	SendResetToken func(token string, user users.LocksmithUserInterface)
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
			ConfigureCustomUser:       options.CustomUserRegistration,
			EmailAsUsername:           options.UseEmailAsUsername,
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

		// This endpoint requires a bit of dynamic Secure Endpointness,
		// so all of that is handled within it.
		mux.Handle("/api/reset-password", reset.ResetRouterAPIHandler{
			Database:       db,
			SendResetToken: options.ResetPasswordOptions.SendResetToken,
		})
	}

	if !options.DisableUI {
		mux.Handle("/login", login.LoginPageHandler{
			AppName:         options.AppName,
			Styling:         options.Styling,
			EmailAsUsername: options.UseEmailAsUsername,
			OnboardingPath:  options.OnboardPath,
		})
		mux.Handle("/register", httpHelpers.InjectDatabaseIntoContext(register.RegistrationPageHandler{
			AppName:                   options.AppName,
			DisablePublicRegistration: options.DisablePublicRegistration,
			Styling:                   options.Styling,
			EmailAsUsername:           options.UseEmailAsUsername,
			HasOnboarding:             len(options.OnboardPath) > 0,
		}, db))
		mux.Handle("/reset-password", httpHelpers.InjectDatabaseIntoContext(reset.ResetPasswordPageHandler{
			AppName:         options.AppName,
			Styling:         options.Styling,
			EmailAsUsername: options.UseEmailAsUsername,
			ShowResetStage:  false,
		}, db))

		mux.Handle("/reset-password/reset", endpoints.SecureEndpointHTTPMiddleware(reset.ResetPasswordPageHandler{
			AppName:         options.AppName,
			Styling:         options.Styling,
			EmailAsUsername: options.UseEmailAsUsername,
			ShowResetStage:  true,
		}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions: []string{"magic.reset.password"},
			PrioritizeMagic:    true,
		}))
	}

	if !options.DisableLocksmithPage {
		serveAdminPage := endpoints.SecureEndpointHTTPMiddleware(administration.ServeAdminPage{}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions: []string{"view.ls-admin"},
		})
		mux.Handle("/locksmith", serveAdminPage)
	}
}
