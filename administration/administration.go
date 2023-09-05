package administration

import (
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

func DeleteUser(db database.DatabaseAccessor, username string) (bool, error) {
	deleted, err := db.DeleteOne("users", map[string]interface{}{
		"username": username,
	})

	if err != nil {
		return false, err
	}

	return deleted, nil
}

type ListUsersOptions struct {
	// Query by roles
	GetRoles []string
	// Allow for custom roles
	CustomInterface users.LocksmithUserInterface
}

func ListUsers(db database.DatabaseAccessor, opts ListUsersOptions) ([]users.PublicLocksmithUserInterface, error) {
	var useStruct users.LocksmithUserInterface

	if opts.CustomInterface == nil {
		useStruct = users.LocksmithUser{}
	} else {
		useStruct = opts.CustomInterface
	}

	query := map[string]interface{}{}

	if len(opts.GetRoles) > 1 {
		query["$or"] = []map[string]interface{}{}

		for _, role := range opts.GetRoles {
			query["$or"] = append(query["$or"].([]map[string]interface{}), map[string]interface{}{
				"role": role,
			})
		}
	} else if len(opts.GetRoles) == 1 {
		query = map[string]interface{}{
			"role": opts.GetRoles[0],
		}
	}

	allUsers, found := db.Find("users", query)

	if !found {
		return []users.PublicLocksmithUserInterface{}, nil
	}

	usersArray := make([]users.PublicLocksmithUserInterface, len(allUsers))

	for i, user := range allUsers {
		useStruct.ReadFromMap(&useStruct, user.(map[string]interface{}))
		public, err := useStruct.ToPublic()

		if err != nil {
			return []users.PublicLocksmithUserInterface{}, nil
		}

		usersArray[i] = public
	}

	return usersArray, nil
}
