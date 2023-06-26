package administration

import (
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
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

func ListUsers(db database.DatabaseAccessor, withStruct ...users.LocksmithUserInterface) ([]users.PublicLocksmithUserInterface, error) {
	var useStruct users.LocksmithUserInterface

	if len(withStruct) == 0 {
		useStruct = users.LocksmithUser{}
	} else {
		useStruct = withStruct[0]
	}

	allUsers, found := db.Find("users", map[string]interface{}{})

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
