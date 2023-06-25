package database

import (
	"testing"
)

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
