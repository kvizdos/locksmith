package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/administration"
	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/authentication/login"
	"github.com/kvizdos/locksmith/authentication/register"
	"github.com/kvizdos/locksmith/authentication/reset"
	"github.com/kvizdos/locksmith/authentication/validation"
	captchaproviders "github.com/kvizdos/locksmith/captcha-providers"
	"github.com/kvizdos/locksmith/components"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/httpHelpers"
	"github.com/kvizdos/locksmith/launchpad"
	"github.com/kvizdos/locksmith/pages"
	sharedmemory "github.com/kvizdos/locksmith/shared-memory"
	"github.com/kvizdos/locksmith/shared-memory/providers"
	"github.com/kvizdos/locksmith/tenant"
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
	InviteUsedRedirect        string
	CustomUserRegistration    register.RegisterCustomUserFunc
	LaunchpadSettings         launchpad.LocksmithLaunchpadOptions
	Styling                   pages.LocksmithPageStyling
	ResetPasswordOptions      ResetPasswordOptions
	HIBPIntegrationOptions    hibp.HIBPSettings
	// map[roleName]time.Duration
	// Use "default" as a catch-all
	InactivityLockDuration         time.Duration
	MinimumPasswordLength          int
	NewRegistrationEvent           func(user users.LocksmithUserInterface)
	SharedMemory                   sharedmemory.MemoryProvider
	LoginSettings                  *login.LoginOptions
	DefaultRegistrationEntitlement string
	// ECDSA Keys used to sign tokens
	TokenSigningPackage validation.ValidationSigningKeys
	UseTenant           tenant.Tenant
}

type ResetPasswordOptions struct {
	SendResetToken func(token string, user users.LocksmithUserInterface)
}

func InitializeLocksmithRoutes(mux *http.ServeMux, db database.DatabaseAccessor, options LocksmithRoutesOptions) {
	if !options.DisableComponents {
		mux.HandleFunc("/components/", components.ServeComponents)
	}

	useSharedMemory := options.SharedMemory
	if useSharedMemory == nil {
		useSharedMemory = providers.NewRamSharedMemoryProvider()
	}

	if options.UseTenant == nil {
		options.UseTenant = tenant.BaseTenant{}
	}

	useLoginSettings := options.LoginSettings
	if useLoginSettings == nil {
		useLoginSettings = &login.LoginOptions{
			LockoutPolicy: login.LockoutPolicy{
				CaptchaAfter: 3,
				LockoutAfter: 10,
				ResetAfter:   24 * time.Hour,
				OnLockout: func(username string) {
					fmt.Println(username, "locked due to too many incorrect passwords")
				},
			},
		}
	}
	useLoginSettings.CaptchaProvider = captchaproviders.UseProvider

	InitializeLaunchpad(mux, db, options)

	jwtAPIHandler := httpHelpers.InjectDatabaseIntoContext(validation.JWTEndpointHandler{
		CustomUserOptions: validation.JWTCustomUserOptions{},
	}, db)
	mux.Handle("/jwt", jwtAPIHandler)

	if !options.DisableAPI {
		var lockAccountsAfter time.Duration

		if options.InactivityLockDuration == 0 {
			// If no lock period specified,
			// keep accounts open for 100 years.
			lockAccountsAfter = 24 * 365 * 100 * time.Hour
		} else {
			lockAccountsAfter = options.InactivityLockDuration
		}

		registrationAPIHandler := httpHelpers.InjectDatabaseIntoContext(register.RegistrationHandler{
			DefaultRoleName:                "user",
			DisablePublicRegistration:      options.DisablePublicRegistration,
			ConfigureCustomUser:            options.CustomUserRegistration,
			EmailAsUsername:                options.UseEmailAsUsername,
			HIBP:                           options.HIBPIntegrationOptions,
			MinimumLengthRequirement:       options.MinimumPasswordLength,
			NewRegistrationEvent:           options.NewRegistrationEvent,
			DefaultRegistrationEntitlement: options.DefaultRegistrationEntitlement,
		}, db)
		mux.Handle("/api/register", registrationAPIHandler)

		loginAPIHandler := httpHelpers.InjectDatabaseIntoContext(login.LoginHandler{
			HIBP:                options.HIBPIntegrationOptions,
			LockInactivityAfter: lockAccountsAfter,
			Options:             *useLoginSettings,
			SharedMemory:        useSharedMemory,
			TokenSigningKeys:    options.TokenSigningPackage,
			TenantInterface:     options.UseTenant,
		}, db)
		mux.Handle("/api/login", loginAPIHandler)

		listUsersAdminAPIHandler := endpoints.SecureEndpointHTTPMiddleware(administration.AdministrationListUsersHandler{}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions:  []string{"AUTHENTICATION.users.list.all"},
			RequiresEntitlement: []string{"SKU_AUTHENTICATION"},
		})
		mux.Handle("/api/users/list", listUsersAdminAPIHandler)

		lockStatusAPI := endpoints.SecureEndpointHTTPMiddleware(administration.AdministrationLockStatusAPI{
			LockInactivityAfter: lockAccountsAfter,
		}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions:  []string{"AUTHENTICATION.lock"},
			RequiresEntitlement: []string{"SKU_AUTHENTICATION"},
		})
		mux.Handle("/api/users/lock-status", lockStatusAPI)

		deleteUserAdminAPIHandler := endpoints.SecureEndpointHTTPMiddleware(administration.AdministrationDeleteUsersHandler{}, db, endpoints.EndpointSecurityOptions{
			RequiresEntitlement: []string{"SKU_AUTHENTICATION"},
		})
		mux.Handle("/api/users/delete", deleteUserAdminAPIHandler)

		if !options.DisableInvites {
			inviteUserAPIHandler := endpoints.SecureEndpointHTTPMiddleware(invitations.AdministrationInviteUserHandler{}, db, endpoints.EndpointSecurityOptions{
				MinimalPermissions:  []string{"AUTHENTICATION.user.invite"},
				RequiresEntitlement: []string{"SKU_AUTHENTICATION"},
			})
			mux.Handle("/api/users/invite", inviteUserAPIHandler)
		}

		// This endpoint requires a bit of dynamic Secure Endpointness,
		// so all of that is handled within it.
		mux.Handle("/api/reset-password", reset.ResetRouterAPIHandler{
			Database:              db,
			SendResetToken:        options.ResetPasswordOptions.SendResetToken,
			HIBP:                  options.HIBPIntegrationOptions,
			MinimumPasswordLength: options.MinimumPasswordLength,
		})
	}

	if !options.DisableUI {
		mux.Handle("/login", login.LoginPageMiddleware{
			Next: login.LoginPageHandler{
				AppName:         options.AppName,
				Styling:         options.Styling,
				EmailAsUsername: options.UseEmailAsUsername,
				OnboardingPath:  options.OnboardPath,
				CaptchaProvider: captchaproviders.UseProvider,
			},
		})
		mux.Handle("/register", httpHelpers.InjectDatabaseIntoContext(register.RegistrationPageHandler{
			AppName:                   options.AppName,
			DisablePublicRegistration: options.DisablePublicRegistration,
			Styling:                   options.Styling,
			EmailAsUsername:           options.UseEmailAsUsername,
			HasOnboarding:             len(options.OnboardPath) > 0,
			InviteUsedRedirect:        options.InviteUsedRedirect,
			HIBPIntegrationOptions:    options.HIBPIntegrationOptions,
			MinimumLengthRequirement:  options.MinimumPasswordLength,
		}, db))
		mux.Handle("/reset-password", httpHelpers.InjectDatabaseIntoContext(reset.ResetPasswordPageHandler{
			AppName:               options.AppName,
			Styling:               options.Styling,
			EmailAsUsername:       options.UseEmailAsUsername,
			ShowResetStage:        false,
			HIBP:                  options.HIBPIntegrationOptions,
			MinimumPasswordLength: options.MinimumPasswordLength,
		}, db))

		mux.Handle("/reset-password/reset", endpoints.SecureEndpointHTTPMiddleware(reset.ResetPasswordPageHandler{
			AppName:               options.AppName,
			Styling:               options.Styling,
			EmailAsUsername:       options.UseEmailAsUsername,
			ShowResetStage:        true,
			HIBP:                  options.HIBPIntegrationOptions,
			MinimumPasswordLength: options.MinimumPasswordLength,
		}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions:  []string{"AUTHENTICATION.magic.reset.password"},
			PrioritizeMagic:     true,
			RequiresEntitlement: []string{options.DefaultRegistrationEntitlement},
		}))
	}

	if !options.DisableLocksmithPage {
		serveAdminPage := endpoints.SecureEndpointHTTPMiddleware(administration.ServeAdminPage{}, db, endpoints.EndpointSecurityOptions{
			MinimalPermissions:  []string{"AUTHENTICATION.view.ls-admin"},
			RequiresEntitlement: []string{"SKU_AUTHENTICATION"},
		})
		mux.Handle("/locksmith", serveAdminPage)
	}
}
