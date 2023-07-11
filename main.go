package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kvizdos/locksmith/administration"
	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/endpoints"
	"github.com/kvizdos/locksmith/authentication/login"
	"github.com/kvizdos/locksmith/authentication/register"
	"github.com/kvizdos/locksmith/components"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/httpHelpers"
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
	fmt.Fprintf(w, "Hello, %s @ %s", parsed.Username, parsed.Token)
}

func main() {
	http.HandleFunc("/components/", components.ServeComponents)

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

	registrationAPIHandler := httpHelpers.InjectDatabaseIntoContext(register.RegistrationHandler{
		DefaultRoleName:           "user",
		DisablePublicRegistration: false,
	}, db)
	loginAPIHandler := httpHelpers.InjectDatabaseIntoContext(login.LoginHandler{}, db)

	listUsersAdminAPIHandler := endpoints.SecureEndpointHTTPMiddleware(administration.AdministrationListUsersHandler{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{"users.list.all"},
	})
	deleteUserAdminAPIHandler := endpoints.SecureEndpointHTTPMiddleware(administration.AdministrationDeleteUsersHandler{}, db)

	inviteUserAPIHandler := endpoints.SecureEndpointHTTPMiddleware(invitations.AdministrationInviteUserHandler{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{"user.invite"},
	})

	serveAdminPage := endpoints.SecureEndpointHTTPMiddleware(administration.ServeAdminPage{}, db, endpoints.EndpointSecurityOptions{
		MinimalPermissions: []string{"view.ls-admin"},
	})

	serveAppPage := endpoints.SecureEndpointHTTPMiddleware(TestAppHandler{}, db)

	http.Handle("/api/login", loginAPIHandler)
	http.Handle("/api/register", registrationAPIHandler)

	http.Handle("/api/users/list", listUsersAdminAPIHandler)
	http.Handle("/api/users/delete", deleteUserAdminAPIHandler)
	http.Handle("/api/users/invite", inviteUserAPIHandler)

	http.Handle("/app", serveAppPage)
	http.Handle("/login", login.LoginPageHandler{})
	http.Handle("/register", httpHelpers.InjectDatabaseIntoContext(register.RegistrationPageHandler{
		DisablePublicRegistration: false,
	}, db))

	http.Handle("/locksmith", serveAdminPage)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
