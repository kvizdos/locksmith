package administration

import (
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

func ListUsers(db database.DatabaseAccessor, withStruct ...users.LocksmithUserStruct) ([]users.LocksmithUserStruct, error) {
	var useStruct users.LocksmithUserStruct

	if len(withStruct) == 0 {
		useStruct = users.LocksmithUser{}
	} else {
		useStruct = withStruct[0]
	}

	allUsers, found := db.Find("users", map[string]interface{}{})

	if !found {
		return []users.LocksmithUserStruct{}, nil
	}

	usersArray := make([]users.LocksmithUserStruct, len(allUsers))

	for i, user := range allUsers {
		useStruct.ReadFromMap(&useStruct, user.(map[string]interface{}))
		usersArray[i] = useStruct
	}

	return usersArray, nil
}
