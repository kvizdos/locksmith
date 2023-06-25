package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"kv.codes/locksmith/administration"
	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/authentication/login"
	"kv.codes/locksmith/authentication/register"
	"kv.codes/locksmith/authentication/validation"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/httpHelpers"
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
	fs := http.FileServer(http.Dir("./components"))
	http.Handle("/components/", http.StripPrefix("/components/", fs))

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
	err := db.Initialize("mongodb://localhost:27017", "locksmith")

	if err != nil {
		fmt.Println(err)
		return
	}

	registrationAPIHandler := httpHelpers.InjectDatabaseIntoContext(register.RegistrationHandler{}, db)
	loginAPIHandler := httpHelpers.InjectDatabaseIntoContext(login.LoginHandler{}, db)

	serveAppPage := validation.ValidateUserTokenMiddleware(TestAppHandler{}, db)

	http.Handle("/api/login", loginAPIHandler)
	http.Handle("/api/register", registrationAPIHandler)

	http.Handle("/app", serveAppPage)
	http.HandleFunc("/login", login.ServeLoginPage)
	http.HandleFunc("/register", register.ServeRegisterPage)

	http.HandleFunc("/locksmith", administration.ServeAdminPage)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
