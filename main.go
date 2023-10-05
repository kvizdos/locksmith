package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/launchpad"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/routes"
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
	c, _ := r.Cookie("token")
	parsed, _ := authentication.ParseToken(c.Value)
	// w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><script src=\"/components/persona-switcher.component.js\" type=\"module\"></script><p>Hello World: %s - %s</p></html>", parsed.Username, parsed.Token)
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
				fmt.Println("Nothing to bootstrap.")
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

	res, _ := db.FindPaginated("users", map[string]interface{}{}, 1, "")
	fmt.Println("RES", res)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}

}
