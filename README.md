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

In this example, `db` can be the MongoDB structure this package provides (which gives the features necessary for this project to work and is slightly limited), or any Database structure that follows this interface:

```
type DatabaseUpdateActions string

const (
	SET  DatabaseUpdateActions = "set"
	PUSH DatabaseUpdateActions = "push"
)

type DatabaseAccessor interface {
	InsertOne(table string, body map[string]interface{}) (interface{}, error)
	UpdateOne(table string, query map[string]interface{}, body map[DatabaseUpdateActions]map[string]interface{}) (interface{}, error)
	FindOne(table string, query map[string]interface{}) (interface{}, bool)
	Find(table string, query map[string]interface{}) ([]interface{}, bool)
}
```

### Customizing Users

Sometimes, you may want to store custom data on users. This can be accomplished by:
1. Create a User structure that defines any new method and *inherits the LocksmithUserInterface*
2. Define your custom User structure
3. Override the default `ReadFromMap()` function to read in any new data that would be stored on the user

In code, this looks like:

```
type CustomUserInterface interface {
	users.LocksmithUserInterface
}

type CustomUser struct {
	users.LocksmithUser

	customObject string
}

func (c CustomUser) ReadFromMap(writeTo *users.LocksmithUserInterface, u map[string]interface{}) {
	// Load initial locksmith data
	var user users.LocksmithUserInterface
	c.LocksmithUser.ReadFromMap(&user, u)
	lsu := user.(users.LocksmithUser)

	converted := customUser{
		LocksmithUser: lsu,
	}

	converted.customObject = u["customObject"].(string)

	*writeTo = converted
}
```

Once you've read in your user from a MongoDB database, you can use it in the code by running:

```
// mongoUserData is as such:
// var mongoUserData map[string]interface{}

var user LocksmithUserInterface
customUser{}.ReadFromMap(&user, mongoUserData)
converted := user.(CustomUser)
```

Specific Locksmith features can also be used with custom interfaces, like `ListUsers()` by passing a base structure that inherits from `LocksmithUserInterface`:

```
usersArr, err := administration.ListUsers(testDb, CustomUser{})

// This will return an array of CustomUsers instead of LocksmithUsers
```

## Contribute

This project has a list of todo's at [TODO.md](./TODO.md). Feel free to submit a PR! Please reasonably test your code as you go along :)
