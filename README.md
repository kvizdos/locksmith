# Locksmith

Locksmith stands as a comprehensive authentication solution designed specifically for web applications. Beyond secure registration and strict username criteria, it boasts features such as Magic Tokens with Time-To-Live (TTL) for temporary access, a dynamic admin panel for user management, multifactor authentication options, and a user lockout system to combat brute force attacks. By securely storing authentication tokens as cookies, Locksmith hopes to provide users a seamless login experience without compromising on data security.

Key Features:

- Intuitive Web Interface: Easily integrate a pre-built web interface for registration and logins. Alternatively, utilize only the Locksmith Authentication endpoints with your custom UI.
- Token Storage: Authentication tokens are securely stored as cookies for enhanced security.
- Magic Tokens with Scoped Permissions: Grant users access to specific areas without requiring login using Magic Access Codes (MACs). Each token has a Time-To-Live (TTL) ensuring limited access duration. Ideal for password reset links or any notification-based URLs, allowing users to seamlessly interact with the app while maintaining limited permissions.
- Middleware for Security: Implement token and role validation middleware to safeguard protected endpoints effectively and with ease.
- Admin User Management: A comprehensive, pre-built admin panel that allows for viewing and removal of users with ease (or BYO and use our endpoints!).
- Role-based Access: Define user roles, create custom roles, and assign administrative privileges as required.
- User Invitation System: Simplify the onboarding process with a streamlined user invitation mechanism.
- Robust Password Policies: Set customizable password policies to ensure user data remains protected.
- Brute Force Mitigation: A user lockout system to deter and reduce potential brute force attacks.
- Multifactor Authentication: Incorporate multifactor authentication using TOTP, WebAuthn, and PassKeys for added security.
- Password Reset Capabilities: Integrated password reset flows and endpoints. Choose between our UI or integrate your own for a seamless experience.

Locksmith delivers a robust authentication system, balancing top-tier user account security with a seamless login experience.

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
1. Create a User structure that defines any new method and *inherits the* `LocksmithUserInterface`
2. Define your custom User structure
3. Override the default `ReadFromMap()` function to read in any new data that would be stored on the user

In code, this looks like:

```
type customUserInterface interface {
	users.LocksmithUserInterface
}

type customUser struct {
	users.LocksmithUser

	CustomObject string
}

func (c customUser) ReadFromMap(writeTo *users.LocksmithUserInterface, u map[string]interface{}) {
	// Load initial locksmith data
	var user users.LocksmithUserInterface
	c.LocksmithUser.ReadFromMap(&user, u)
	lsu := user.(users.LocksmithUser)

	converted := customUser{
		LocksmithUser: lsu,
	}

	// Begin setting your own values here.
	// Make sure to inherit all Locksmith values by using the code above.
	converted.CustomObject = u["customObject"].(string)

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

## Customizing ListUsers() API

By default, Locksmith's `ListUsers()` API returns a simplified version of the `LocksmithUser{}` struct with less sensitive information. When you use the code below, it will only return user info that is supplied by the `PublicLocksmithUserInterface`:

```
usersArr, err := ListUsers(testDb)
```

Would return an array of `PublicLocksmithUser{}`'s:

```
type PublicLocksmithUser struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	ActiveSessionCount int    `json:"sessions"`
	LastActive         int64  `json:"lastActive"`
}
```

Using the example above, we can add the "CustomObject" by:
1. Creating a structure (e.g. `publicCustomUser`) that inherits from `users.PublicLocksmithUser` (this will ensure no Locksmith features break)
2. Defining custom keys in the structure (make sure to add `json` tags)
3. Add a `FromRegular()` function to `publicCustomUser{}`
4. Add a `ToPublic()` function to `customUser{}`

Let's hit this step by step.

**1 + 2. Create the new structure and define keys:**
```
type publicCustomUser struct {
	users.PublicLocksmithUser

	CustomObject string `json:"customObject"`
}
```

**3. Add `FromRegular()` to `publicCustomUser{}`:**
```
func (u publicCustomUser) FromRegular(user users.LocksmithUserInterface) (users.PublicLocksmithUserInterface, error) {
	lsPub, err := u.PublicLocksmithUser.FromRegular(user)

	if err != nil {
		return publicCustomUser{}, nil
	}

	publicUser := publicCustomUser{
		PublicLocksmithUser: lsPub.(users.PublicLocksmithUser),
	}

	// Begin setting your own values here.
	// Make sure to inherit all Locksmith values by using the code above.
	publicUser.CustomObject = user.(customUser).CustomObject

	return publicUser, nil
}
```

**4. Add `ToPublic()` to `customUser{}`:**
```
func (u customUser) ToPublic() (users.PublicLocksmithUserInterface, error) {
	publicUser, err := publicCustomUser{}.FromRegular(u)

	return publicUser, err
}
```

Done! Now, when you want to list users with the custom data types, use the following code:

```
publicUser := customUser{}
usersArr, err := ListUsers(testDb, publicUser)
```

**Note that you're passing the default customUser{} and not the publicCustomUser{}**

Alongside `ListUsers()`, you can also pass the `customUser{}` interface into the `Administration*Handler{}` to return the values over the HTTP API:

```
listUsersAdminAPIHandler := validation.ValidateUserTokenMiddleware(administration.AdministrationListUsersHandler{
    UserInterface: customUser{},
}, db)
```

*Works with AdministrationListUsersHandler*

### Custom Users Security Notice

Even if you don't want custom values to be passed on the structure, **never just pass the LocksmithUser{} through a JSON serializer to get the output.** Always use the `.ToPublic()` output beforehand to ensure sensitive user data (like password hashes and sessions) are not leaked.

## Contribute

This project has a list of todo's at [TODO.md](./TODO.md). Feel free to submit a PR! Please reasonably test your code as you go along :)
