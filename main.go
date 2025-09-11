package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/authentication/login"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/authentication/oauth"
	oauth_google_oidc "github.com/kvizdos/locksmith/authentication/oauth/oidc"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/authentication/xsrf"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/jwts"
	"github.com/kvizdos/locksmith/launchpad"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/ratelimits"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/routes"
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
	role, _ := authUser.GetRole()
	sid, _ := r.Context().Value("sid").(string)
	w.Write([]byte(fmt.Sprintf(`<html>
		<head>
			<script defer type="text/javascript" src="/api/auth/oauth/keep-alive.js"></script>
		</head>
		<body>
			<p>%s -- %t - %s %s %d - Your SID is %s</p>
		</body>
		</html>
		`, r.URL.Path, authUser.IsMagic(), authUser.GetUsername(), role.Name, len(role.Permissions), sid)))
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
	googleOIDC, err := oauth_google_oidc.NewOIDCConnection(context.Background(), oauth_google_oidc.OIDCConnectionParams{
		Issuer:       "https://accounts.google.com",
		ClientID:     "ExampleClientID",
		ClientSecret: "ExampleClientSecret",
		BaseURL:      "https://example.com",
		ProviderName: "google", // for the UI & a few backend things; make sure its unique!
		DB:           db,
		LoginInfoCallback: func(method string, user map[string]any) {
			fmt.Printf("User logged in via Google: %+v", user)
		},
	})
	if err != nil {
		panic(err)
	}

	routes.InitializeLocksmithRoutes(mux, db, routes.LocksmithRoutesOptions{
		AppName:            "Demo App",
		UseEmailAsUsername: true,
		OnboardPath:        "/onboard",
		InviteUsedRedirect: "/app",
		LoginInfoCallback: func(method string, user map[string]any) {
			fmt.Printf("User logged in via username / password: %+v", user)
		},
		OAuthProviders: []oauth.OAuthProvider{
			googleOIDC,
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
		InactivityLockDuration: map[string]time.Duration{
			"default": 100 * 365 * 24 * time.Hour, // If not set, defaults to 100 years.
			"admin":   100 * 365 * 24 * time.Hour,
		},
		Styling: pages.LocksmithPageStyling{
			LogoURL: "/components/locksmith.svg",
			InjectHeader: template.HTML(
				`<script>
					console.log("Loaded page.")
				</script>`,
			),
		},
		ResetPasswordOptions: routes.ResetPasswordOptions{
			SendResetToken: printResetToken,
		},
		DisablePublicRegistration: false,
		MinimumPasswordLength:     8,
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
			Users: map[string]launchpad.LocksmithLaunchpadUserOptions{
				"lp-admin": {
					DisplayName: "Administrator",
					Description: "Add a description for extra context.",
					Email:       "admin@admin.com",
					Role:        "admin",
					Redirect:    "/locksmith",
					Custom: map[string]interface{}{
						"customObject": "hello world",
					},
				},
				"lp-user": {
					DisplayName: "General User",
					Email:       "user@user.com",
					Role:        "user",
					Redirect:    "/app",
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
					Role:        "user",
					Redirect:    "/app",
				},
			},
		},
	})

	roles.AddPermissionsToRole("user", []string{
		"can.see.user.view",
		"can.see.both.view",
	})

	// Everyone can see this page, including magic-only!
	// You only need to be "authenticated" to see this page.
	// Use permissionless-endpoints with caution.
	mux.Handle("/app", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		RateLimit: ratelimits.NewRateLimiter(15, 15).WithSecondsLimits(5, 5),
	}))

	// Only the Logged-In users can see this page
	// Magic-only users cannot see this.
	mux.Handle("/user", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{
			"can.see.user.view",
		},
	}))

	mux.Handle("/both", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{
			"can.see.both.view",
		},
	}))

	// Logged in w/ Magic users can see this
	mux.Handle("/magic-prioritized", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		// Priority to Magic, so Magic permissions will override
		// default role permissions.
		PrioritizeMagic: true,
		MinimalPermissions: []string{
			"can.see.magic.view",
		},
	}))

	// If the user is also logged-in, they cannot see this page!
	mux.Handle("/magic-only", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db, endpoints.EndpointSecurityOptions{
		// No priority set, however only the Magic tokens have this permission
		// Because priority is on the token, logged-in users will not see this page.
		MinimalPermissions: []string{
			"can.see.magic.view",
		},
	}))

	roles.CreatePermissionSet([]roles.Permission{
		{
			Permission: "issue.token.demo",
			IsElevated: false,
			JWTOnly:    false,
			AvailableRoles: []string{
				"admin",
			},
			DontExposeFrontend: false,
		},
	})
	jwts.RegisterJWT("demo", jwts.RegisteredJWT{
		RequiredPermission: "issue.token.demo",
		AttachPermissions: []string{
			"svc.notify.connect",
		},
		ForAudience: "edvizion-notify",
		Issuer:      "my-app",
		ExpiresIn:   60 * time.Hour,
		SigningKey: func(context.Context) (string, error) {
			return "your_private_key_here", nil
		},
	})

	/*
		Replace this with your OWN key
		pkg, err := signing.CreateSigningPackage()
		marshaledPK, err := pkg.MarshalPrivate() // use this output as the "DecodePrivateKey" variable
	*/
	sp, _ := signing.DecodePrivateKey("MHcCAQEEIOXFnC40e/HNM6nn6iO8u3oA/KMoSyLrzarpJ/UMdTrKoAoGCCqGSM49AwEHoUQDQgAE8ZtLIHX8NYqAe0VukxPGZNHmOv84WVjRDPHATJq/go/eubOIB/ddQ4JG2tEtPqCKa+pso5l/vC1kIzIbZIJIFA==")
	magic.MagicSigningPackage = &sp
	xsrf.XSRFSigningPackage.Anonymous = &sp
	xsrf.XSRFSigningPackage.Authenticated = &sp
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
