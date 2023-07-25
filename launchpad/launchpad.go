//go:build enable_launchpad
// +build enable_launchpad

package launchpad

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
)

func BootstrapUsers(db database.DatabaseAccessor, accessToken string, importUsers map[string]LocksmithLaunchpadUserOptions) {
	password, _ := authentication.CompileLocksmithPassword(accessToken)

	for username, opts := range importUsers {
		_, found := db.FindOne("users", map[string]interface{}{
			"username": username,
		})

		if found {
			fmt.Printf("Launchpad user %s already registered.", username)
			continue
		}

		_, err := db.InsertOne("users", map[string]interface{}{
			"id":          uuid.New().String(),
			"username":    username,
			"password":    password,
			"email":       opts.Email,
			"sessions":    []interface{}{},
			"websessions": []interface{}{},
			"role":        opts.Role,
		})

		if err != nil {
			fmt.Println(err)
		}
	}

}
