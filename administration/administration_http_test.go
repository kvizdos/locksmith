package administration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/users"
)

func TestListUsersInvalidMethod(t *testing.T) {
	handler := AdministrationListUsersHandler{}

	req, err := http.NewRequest("POST", "/locksmith/api/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestListUsersReceivesValidJSON(t *testing.T) {
	var usersList []users.PublicLocksmithUser

	var sessions []interface{}
	tempSessions := []authentication.PasswordSession{
		{
			Token:     "abc",
			ExpiresAt: 123,
		},
	}

	for _, sess := range tempSessions {
		sessions = append(sessions, sess)
	}

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"sessions": sessions,
				},
			},
		},
	}

	handler := AdministrationListUsersHandler{}

	req, err := http.NewRequest("GET", "/locksmith/api/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusOK)
	}

	err = json.Unmarshal(rr.Body.Bytes(), &usersList)

	if err != nil {
		t.Errorf("failed to decode JSON response")
		return
	}

	if len(usersList) != 1 {
		t.Errorf("received unexpected number of users, expected 1 got %d", len(usersList))
		return
	}

	if usersList[0].ActiveSessionCount != 1 {
		t.Errorf("error reading correct number of sessions, expected 1 got %d", usersList[0].ActiveSessionCount)
	}
}

func TestListUsersReceivesValidJSONWithCustomStruct(t *testing.T) {
	var usersList []publicCustomUser

	var sessions []interface{}
	tempSessions := []authentication.PasswordSession{
		{
			Token:     "abc",
			ExpiresAt: 123,
		},
	}

	for _, sess := range tempSessions {
		sessions = append(sessions, sess)
	}

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"sessions":     sessions,
					"customObject": "hello",
				},
			},
		},
	}

	handler := AdministrationListUsersHandler{
		UserInterface: customUser{},
	}

	req, err := http.NewRequest("GET", "/locksmith/api/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusOK)
	}

	err = json.Unmarshal(rr.Body.Bytes(), &usersList)

	if err != nil {
		t.Errorf("failed to decode JSON response")
		return
	}

	if len(usersList) != 1 {
		t.Errorf("received unexpected number of users, expected 1 got %d", len(usersList))
		return
	}

	if usersList[0].ActiveSessionCount != 1 {
		t.Errorf("error reading correct number of sessions, expected 1 got %d", usersList[0].ActiveSessionCount)
	}

	if usersList[0].CustomObject != "hello" {
		t.Errorf("could not read custom object: %s", usersList[0].CustomObject)
	}
}