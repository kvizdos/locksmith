package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/authentication/login"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/authentication/validation"
	"github.com/kvizdos/locksmith/authentication/xsrf"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/entitlements"
	"github.com/kvizdos/locksmith/launchpad"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/ratelimits"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/routes"
	"github.com/kvizdos/locksmith/tenant"
	"github.com/kvizdos/locksmith/users"
)

// import

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type TestAppHandler struct{}

func (th TestAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authUser := r.Context().Value("authUser").(users.LocksmithUserInterface)
	sid, _ := r.Context().Value("sid").(string)
	roles := authUser.GetRoles()
	w.Write([]byte(fmt.Sprintf("%t - %s %s - Your SID is %s", authUser.IsMagic(), authUser.GetUsername(), roles, sid)))
}

func printResetToken(token string, user users.LocksmithUserInterface) {
	fmt.Println(user.GetID(), token)
}

func sendWelcomeEmailExample(u users.LocksmithUserInterface) {
	fmt.Printf("Sending welcome email to %s\n", u.GetEmail())
}

func main() {
	// testPassword, _ := authentication.CompileLocksmithPassword("pass")

	// db := database.TestDatabase{
	// 	Tables: map[string]map[string]interface{}{
	// 		"users": {
	// 			"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
	// 				"id":       "c8531661-22a7-493f-b228-028842e09a05",
	// 				"username": "kenton",
	// 				"password": testPassword,
	// 				"sessions": []interface{}{},
	// 			},
	// 		},
	// 	},
	// }
	ctx, timeout := context.WithTimeout(context.Background(), 20*time.Second)
	db := database.MongoDatabase{
		Ctx:    ctx,
		Cancel: timeout,
	}
	err := db.Initialize("mongodb://localhost:27017", os.Getenv("database"))

	if err != nil {
		fmt.Println(err)
		return
	}

	mux := http.NewServeMux()
	// captchaproviders.UseProvider = providers.TurnstileCaptcha{
	// 	SiteKey:   "xxx",
	// 	SecretKey: "yyy",
	// }

	sp, _ := signing.DecodePrivateKey("MHcCAQEEIOXFnC40e/HNM6nn6iO8u3oA/KMoSyLrzarpJ/UMdTrKoAoGCCqGSM49AwEHoUQDQgAE8ZtLIHX8NYqAe0VukxPGZNHmOv84WVjRDPHATJq/go/eubOIB/ddQ4JG2tEtPqCKa+pso5l/vC1kIzIbZIJIFA==")
	magic.MagicSigningPackage = &sp
	xsrf.XSRFSigningPackage.Anonymous = &sp
	xsrf.XSRFSigningPackage.Authenticated = &sp

	validation.SetSigningKeys(validation.ValidationSigningKeys{
		PublicKey:  sp.Public,
		PrivateKey: sp.Private,
	})

	userManagementTenantID := "be1428a8-a23b-4cce-b694-4f90fd5ebf07"
	generalTenantID := "d00afe37-a1b8-4865-96f0-cff5c3003cea"
	baseAppSKUID := "SKU_TEST_APP"
	managementAppSKUID := "SKU_AUTHENTICATION"
	routes.InitializeLocksmithRoutes(mux, db, routes.LocksmithRoutesOptions{
		AppName:                        "Locksmith Demo UI",
		UseEmailAsUsername:             true,
		OnboardPath:                    "/onboard",
		InviteUsedRedirect:             "/app",
		DefaultRegistrationEntitlement: managementAppSKUID, // This needs to be given to every user.
		TokenSigningPackage: validation.ValidationSigningKeys{
			PublicKey:  &sp.Private.PublicKey,
			PrivateKey: sp.Private,
		},
		LoginSettings: &login.LoginOptions{
			LockoutPolicy: login.LockoutPolicy{
				CaptchaAfter: 2,
				LockoutAfter: 10,
				ResetAfter:   time.Duration(24 * time.Hour),
				OnLockout: func(username string) {
					fmt.Println(username, "locked out")
				},
			},
		},
		InactivityLockDuration: 24 * 30 * time.Hour,
		Styling: pages.LocksmithPageStyling{
			LogoURL: "https://example.com/logo.webp",
			InjectHeader: template.HTML(
				`<script>
					console.log("Loaded page.")
				</script>`,
			),
		},
		ResetPasswordOptions: routes.ResetPasswordOptions{
			SendResetToken: printResetToken,
		},
		MinimumPasswordLength: 8,
		HIBPIntegrationOptions: hibp.HIBPSettings{
			Enabled:                  false, // true to enable
			AppName:                  "Locksmith Demo",
			Enforcement:              hibp.STRICT,
			HTTPClient:               &http.Client{},
			PasswordSecurityInfoLink: "https://github.com/kvizdos",
		},
		NewRegistrationEvent: sendWelcomeEmailExample,
		LaunchpadSettings: launchpad.LocksmithLaunchpadOptions{
			Enabled:                       true,
			IsEarlyDevelopmentEnvironment: false,
			Caption:                       "Locksmith Launchpad helps demo your service. It allow stakeholders to easily swap between users and feel the product from every POV- without worrying about passwords.",
			AccessToken:                   "changeme",
			BootstrapDatabase: func(da database.DatabaseAccessor) {
				db.Drop("users")
			},
			RefreshButtonText: "Refresh Environment",
			Tenants: []tenant.Tenant{
				tenant.BaseTenant{
					ID: uuid.MustParse(userManagementTenantID),
					Entitlements: []entitlements.TenantEntitlement{
						{
							ID:            managementAppSKUID,
							Quantity:      100,
							QuantitySpent: 1,
							StartDate:     time.Now().AddDate(0, 0, -1).UTC(),
							EndDate:       time.Now().AddDate(0, 0, 30).UTC(),
						},
						{
							ID:            baseAppSKUID,
							Quantity:      100,
							QuantitySpent: 1,
							StartDate:     time.Now().AddDate(0, 0, -1).UTC(),
							EndDate:       time.Now().AddDate(0, 0, 30).UTC(),
						},
					},
				},
				tenant.BaseTenant{
					ID: uuid.MustParse(generalTenantID),
					Entitlements: []entitlements.TenantEntitlement{
						{
							ID:            baseAppSKUID,
							Quantity:      100,
							QuantitySpent: 1,
							StartDate:     time.Now().AddDate(0, 0, -1).UTC(),
							EndDate:       time.Now().AddDate(0, 0, 30).UTC(),
						},
					},
				},
			},
			Users: map[string]launchpad.LocksmithLaunchpadUserOptions{
				"lp-admin": {
					DisplayName: "Administrator",
					Email:       "admin@admin.com",
					Roles: []string{
						"AUTHENTICATION.admin",
						"TEST_APP.user",
					},
					Entitlements: []string{"SKU_AUTHENTICATION", "SKU_TEST_APP"},
					TenantID:     userManagementTenantID,
					Redirect:     "/locksmith",
					Custom: map[string]interface{}{
						"customObject": "hello world",
					},
				},
				"lp-user": {
					DisplayName: "General User",
					Email:       "user@user.com",
					Roles: []string{
						"SKU_TEST_APP.user",
					},
					Entitlements: []string{"SKU_AUTHENTICATION", "SKU_TEST_APP"},
					TenantID:     generalTenantID,
					Redirect:     "/app",
					Custom: map[string]interface{}{
						// If you need a static ID for testing / interacting
						// with other places, it's useful
						// to set that here.
						"id": "41084e13-a40a-42e7-aac6-19cba36b1d68",
					},
				},
				"lp-user-2": {
					DisplayName: "Another General User",
					Email:       "user2@user.com",
					TenantID:    generalTenantID,
					Roles: []string{
						"SKU_TEST_APP.user",
					},
					Entitlements: []string{"SKU_AUTHENTICATION", "SKU_TEST_APP"},

					Redirect: "/app",
				},
			},
		},
	})

	roles.AddPermissionsToRole("AUTHENTICATION.admin", []string{
		"AUTHENTICATION.can.see.user.view",
		"AUTHENTICATION.can.see.both.view",
	})

	roles.AddFrontendPermissionsToRole("AUTHENTICATION.admin", []string{"AUTHENTICATION.user.create", "AUTHENTICATION.user.delete", "AUTHENTICATION.user.modify"})

	// Everyone can see this page, including magic-only!
	// You only need to be "authenticated" to see this page.
	// Use permissionless-endpoints with caution.
	mux.Handle("/app", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		RateLimit:           ratelimits.NewRateLimiter(15, 15).WithSecondsLimits(5, 5),
		RequiresEntitlement: []string{"SKU_TEST_APP"},
	}))

	// Only the Logged-In users can see this page
	// Magic-only users cannot see this.
	mux.Handle("/user", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{
			"can.see.user.view",
		},
		RequiresEntitlement: []string{"SKU_TEST_APP"},
	}))

	mux.Handle("/both", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{
			"can.see.both.view",
		},
		RequiresEntitlement: []string{"SKU_TEST_APP"},
	}))

	// Logged in w/ Magic users can see this
	mux.Handle("/magic-prioritized", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		// Priority to Magic, so Magic permissions will override
		// default role permissions.
		PrioritizeMagic: true,
		MinimalPermissions: []string{
			"can.see.magic.view",
		},
		RequiresEntitlement: []string{"SKU_TEST_APP"},
	}))

	// If the user is also logged-in, they cannot see this page!
	mux.Handle("/magic-only", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		// No priority set, however only the Magic tokens have this permission
		// Because priority is on the token, logged-in users will not see this page.
		MinimalPermissions: []string{
			"can.see.magic.view",
		},
		RequiresEntitlement: []string{"SKU_TEST_APP"},
	}))

	/*
		Replace this with your OWN key
		pkg, err := signing.CreateSigningPackage()
		marshaledPK, err := pkg.MarshalPrivate() // use this output as the "DecodePrivateKey" variable
	*/
	// _, err = users.LocksmithUser{
	// 	ID: "41084e13-a40a-42e7-aac6-19cba36b1d68",
	// }.CreateMagicAuthenticationCode(db, magic.MagicAuthenticationVariables{
	// 	ForUserID: "41084e13-a40a-42e7-aac6-19cba36b1d68",
	// 	AllowedPermissions: []string{
	// 		"can.see.magic.view",
	// 		"can.see.magic.view.prioritized",
	// 		"can.see.both.view",
	// 	},
	// 	TTL: 15 * time.Minute,
	// })

	// fmt.Println("User Magic Key:", macID, err)
	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}

}
