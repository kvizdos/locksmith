# Locksmith

Locksmith is a robust authentication project designed to simplify user authentication in web applications. It provides secure registration and login features, enforcing strong username requirements. By saving authentication tokens as cookies, Locksmith offers a seamless login experience while maintaining data security.

Key Features:

- User-friendly, drop-in web interface for registration and logins (or, bring your own and just use the Locksmith Authentication endpoints)
- Securely stores authentication tokens as cookies
- Token & role validation middleware to restrict access to protected endpoints
- Admin panel for user management, including viewing and removal of users
- User roles, customizable roles, and administrative privileges
- User invitation system for streamlined onboarding
- Customizable password policies to protect your users
- User Lockout system to reduce brute force attacks
- Multifactor authentication for TOTP and WebAuthn (including PassKeys)
- Built-in password reset flows and endpoints (bring your own UI, or use ours!)

With Locksmith, you can easily implement a reliable authentication system that ensures user account security and provides a smooth login experience.

## Get Started

If you are just using Locksmith for basic authentication, you can get started with this boilerplate code:

**this will be changing**

```
package main

func main() {
	fs := http.FileServer(http.Dir("./components"))

	// Open the Locksmith UI components if you aren't bringing your own UI
	http.Handle("/components/", http.StripPrefix("/components/", fs))

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

	// Secure endpoints with just one line!
	serveAppPage := validation.ValidateUserToken(TestAppHandler{}, db)

	// This will open the API endpoints for registration and logging in. As of now, these are static and the paths should not be changed.
	http.Handle("/api/login", loginAPIHandler)
	http.Handle("/api/register", registrationAPIHandler)

	// Serve our pre-built login and registration UI
	http.HandleFunc("/login", login.ServeLoginPage)
	http.HandleFunc("/register", register.ServeRegisterPage)

	// Serve your app!
	http.Handle("/app", serveAppPage)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Contribute

This project has a list of todo's at [TODO.md](./TODO.md). Feel free to submit a PR!
