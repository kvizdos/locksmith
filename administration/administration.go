package administration

import (
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

func ListUsers(db database.DatabaseAccessor, withStruct ...users.LocksmithUserInterface) ([]users.LocksmithUserInterface, error) {
	var useStruct users.LocksmithUserInterface

	if len(withStruct) == 0 {
		useStruct = users.LocksmithUser{}
	} else {
		useStruct = withStruct[0]
	}

	allUsers, found := db.Find("users", map[string]interface{}{})

	if !found {
		return []users.LocksmithUserInterface{}, nil
	}

	usersArray := make([]users.LocksmithUserInterface, len(allUsers))

	for i, user := range allUsers {
		useStruct.ReadFromMap(&useStruct, user.(map[string]interface{}))
		usersArray[i] = useStruct
	}

	return usersArray, nil
}
