package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/launchpad"
	"github.com/kvizdos/locksmith/pages"
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
	w.Write([]byte(fmt.Sprintf("%s %s %d", authUser.GetUsername(), role.Name, len(role.Permissions))))
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
	routes.InitializeLocksmithRoutes(mux, db, routes.LocksmithRoutesOptions{
		AppName:            "Locksmith Demo UI",
		UseEmailAsUsername: false,
		Styling: pages.LocksmithPageStyling{
			LogoURL: "https://example.com/logo.webp",
		},
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
				},
				"lp-user-2": {
					DisplayName: "Another General User",
					Email:       "user@user.com",
					Role:        "user",
					Redirect:    "/app",
				},
			},
		},
	})

	mux.Handle("/app", endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db))

	sp, _ := signing.DecodePrivateKey("MHcCAQEEIOXFnC40e/HNM6nn6iO8u3oA/KMoSyLrzarpJ/UMdTrKoAoGCCqGSM49AwEHoUQDQgAE8ZtLIHX8NYqAe0VukxPGZNHmOv84WVjRDPHATJq/go/eubOIB/ddQ4JG2tEtPqCKa+pso5l/vC1kIzIbZIJIFA==")
	magic.MagicSigningPackage = &sp
	macID, err := users.LocksmithUser{
		ID: "41084e13-a40a-42e7-aac6-19cba36b1d68",
	}.CreateMagicAuthenticationCode(db, magic.MagicAuthenticationVariables{
		ForUserID: "41084e13-a40a-42e7-aac6-19cba36b1d68",
		AllowedPermissions: []string{
			"test.1",
			"test.2",
		},
		TTL: time.Minute,
	})

	fmt.Println(macID, err)
	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}

}
