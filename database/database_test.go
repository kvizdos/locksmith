package database

import (
	"testing"
)

func TestDeleteOne(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
				},
				"ananc-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "ananc-22a7-493f-b228-028842e09a05",
					"username": "bob",
				},
			},
		},
	}

	deleted, err := testDb.DeleteOne("users", map[string]interface{}{
		"username": "kenton",
	})

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if !deleted {
		t.Error("failed to delete item")
		return
	}

	if _, ok := testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"]; ok {
		t.Error("did not delete item")
	}

	if _, ok := testDb.Tables["users"]["ananc-22a7-493f-b228-028842e09a05"]; !ok {
		t.Error("deleted too much")
	}
}

func TestDeleteOneNothingDeleted(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
				},
				"ananc-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "ananc-22a7-493f-b228-028842e09a05",
					"username": "bob",
				},
			},
		},
	}

	deleted, err := testDb.DeleteOne("users", map[string]interface{}{
		"username": "kenton2",
	})

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if deleted {
		t.Error("should not have deleted user")
		return
	}

	if _, ok := testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"]; !ok {
		t.Error("deleted too much 1")
	}

	if _, ok := testDb.Tables["users"]["ananc-22a7-493f-b228-028842e09a05"]; !ok {
		t.Error("deleted too much 2")
	}
}

func TestInsertDatabase(t *testing.T) {
	testDb := TestDatabase{
		Tables: make(map[string]map[string]interface{}),
	}

	_, err := testDb.InsertOne("users", map[string]interface{}{
		"username": "kenton",
	})

	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestGetInsertedValue(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
				},
			},
		},
	}

	value, found := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	if found == false {
		t.Errorf("couldnt find c")
	}

	usernameValue, ok := value.(map[string]interface{})["username"]

	if !ok {
		t.Errorf("could not find username")
		return
	}

	if usernameValue != "kenton" {
		t.Errorf("found invalid username: %s", usernameValue)
	}
}

func TestUpdateOneSET(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
				},
			},
		},
	}

	_, err := testDb.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[DatabaseUpdateActions]map[string]interface{}{
		SET: {
			"username": "kvizdos",
		},
	})

	if err != nil {
		t.Errorf("failed to SET array: %s", err.Error())
		return
	}

	if testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"] != "kvizdos" {
		t.Errorf("found invalid updated key: %s", testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"])
	}
}

func TestUpdateManySET(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
				},
				"z8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "bob",
				},
				"b8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "b8531661-22a7-493f-b228-028842e09a05",
					"username": "jane",
				},
			},
		},
	}

	_, err := testDb.UpdateMany("users", map[string]interface{}{
		"id": "c8531661-22a7-493f-b228-028842e09a05",
	}, map[DatabaseUpdateActions]map[string]interface{}{
		SET: {
			"username": "overwrite",
		},
	})

	if err != nil {
		t.Errorf("failed to SET array: %s", err.Error())
		return
	}

	if testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"] != "overwrite" {
		t.Errorf("c8531661-22a7-493f-b228-028842e09a05 found invalid updated key: %s", testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"])
	}

	if testDb.Tables["users"]["z8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"] != "overwrite" {
		t.Errorf("z8531661-22a7-493f-b228-028842e09a05 found invalid updated key: %s", testDb.Tables["users"]["z8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"])
	}

	if testDb.Tables["users"]["b8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"] != "jane" {
		t.Errorf("b8531661-22a7-493f-b228-028842e09a05 found invalid updated key: %s", testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"])
	}
}

func TestUpdateOneSETMultiple(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"password": "abc",
				},
			},
		},
	}

	_, err := testDb.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[DatabaseUpdateActions]map[string]interface{}{
		SET: {
			"username": "kvizdos",
			"password": "cba",
		},
	})

	if err != nil {
		t.Errorf("failed to SET array: %s", err.Error())
		return
	}

	if testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"] != "kvizdos" {
		t.Errorf("found invalid updated username key: %s", testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"])
	}

	if testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["password"] != "cba" {
		t.Errorf("found invalid updated password key: %s", testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["password"])
	}
}

func TestUpdateOnePUSH(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{},
				},
			},
		},
	}

	_, err := testDb.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[DatabaseUpdateActions]map[string]interface{}{
		PUSH: {
			"sessions": "abc",
		},
	})

	if err != nil {
		t.Errorf("failed to PUSH array: %s", err.Error())
		return
	}

	dbo := testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})

	if dbo["username"] != "kenton" {
		t.Errorf("username key was updated when it wasn't supposed to be: %s", testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"])
	}

	if len(dbo["sessions"].([]interface{})) == 0 {
		t.Error("push did not occur")
		return
	}

	if dbo["sessions"].([]interface{})[0].(string) != "abc" {
		t.Errorf("found incorrect value: %s", dbo["sessions"].([]interface{})[0].(string))
		return
	}
}

func TestUpdateOnePUSHAppend(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	_, err := testDb.UpdateOne("users", map[string]interface{}{
		"username": "kenton",
	}, map[DatabaseUpdateActions]map[string]interface{}{
		PUSH: {
			"sessions": "cba",
		},
	})

	if err != nil {
		t.Errorf("failed to PUSH array: %s", err.Error())
		return
	}

	dbo := testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})

	if dbo["username"] != "kenton" {
		t.Errorf("username key was updated when it wasn't supposed to be: %s", testDb.Tables["users"]["c8531661-22a7-493f-b228-028842e09a05"].(map[string]interface{})["username"])
	}

	if len(dbo["sessions"].([]interface{})) == 1 {
		t.Error("push did not occur")
		return
	}

	if dbo["sessions"].([]interface{})[1].(string) != "cba" {
		t.Errorf("found incorrect value: %s", dbo["sessions"].([]interface{})[0].(string))
		return
	}
}

func TestFindOne(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{"abc"},
				},

				"abcxyz": map[string]interface{}{
					"id":       "abcxyz",
					"username": "joe",
					"sessions": []interface{}{},
				},
			},
		},
	}

	item, found := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	if !found {
		t.Error("couldn't find object in test database")
		return
	}

	type User struct {
		id       string
		username string
		sessions []interface{}
	}

	if item.(map[string]interface{})["id"] != "c8531661-22a7-493f-b228-028842e09a05" {
		t.Error("found incorrect ID")
		return
	}
}

func TestFindOneDoesNotExist(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	_, found := testDb.FindOne("users", map[string]interface{}{
		"username": "kentons",
	})

	if found {
		t.Error("shouldn't have found the object")
		return
	}
}

func TestFindGetsAllWithEmptyQuery(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"sessions": []interface{}{"abc"},
				},
				"abcdze12-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "bob",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	users, found := testDb.Find("users", map[string]interface{}{})

	if !found {
		t.Error("should have found the objects")
		return
	}

	expecting := 2
	if len(users) != expecting {
		t.Errorf("received incorrect number of results, expected %d got %d", expecting, len(users))
	}
}

func TestFindWithOrOnlyOne(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"role":     "abc",
				},
				"abcdze12-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "a8531661-22a7-493f-b228-028842e09a05",
					"username": "bob",
					"role":     "bca",
				},
				"bbcdze12-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "b8531661-22a7-493f-b228-028842e09a05",
					"username": "james",
					"role":     "bca",
				},
			},
		},
	}

	users, found := testDb.Find("users", map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"role": "abc",
			},
		},
	})

	if !found {
		t.Error("should have found the objects")
		return
	}

	expecting := 1
	if len(users) != expecting {
		t.Errorf("received incorrect number of results, expected %d got %d", expecting, len(users))
	}
}

func TestFindWithOrMany(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"role":     "abc",
				},
				"abcdze12-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "a8531661-22a7-493f-b228-028842e09a05",
					"username": "bob",
					"role":     "bca",
				},
				"bbcdze12-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "b8531661-22a7-493f-b228-028842e09a05",
					"username": "james",
					"role":     "bca",
				},
			},
		},
	}

	users, found := testDb.Find("users", map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"role": "abc",
			},
			{
				"role": "bca",
			},
		},
	})

	if !found {
		t.Error("should have found the objects")
		return
	}

	expecting := 3
	if len(users) != expecting {
		t.Errorf("received incorrect number of results, expected %d got %d", expecting, len(users))
	}
}

func TestFindReturnsZeroForEmptyTable(t *testing.T) {
	testDb := TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	users, found := testDb.Find("users", map[string]interface{}{})

	expecting := 0
	if len(users) != expecting {
		t.Errorf("received incorrect number of results, expected %d got %d", expecting, len(users))
	}

	if found {
		t.Error("should have found the objects")
		return
	}
}
