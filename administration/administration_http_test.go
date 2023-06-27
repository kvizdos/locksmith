package administration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
	"kv.codes/locksmith/users"
)

func TestMain(m *testing.M) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
		"user": {
			"view.admin",
			"user.delete.self",
		},
	}

	m.Run()

	roles.AVAILABLE_ROLES = map[string][]string{}
}

func TestListUsersInvalidMethod(t *testing.T) {
	handler := AdministrationListUsersHandler{}

	req, err := http.NewRequest("POST", "/locksmith/api/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusMethodNotAllowed)
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
					"role":     "user",
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
					"role":         "user",
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

	if usersList[0].Role != "user" {
		t.Errorf("did not receive role")
	}
}

func TestDeleteUserHTTPInvalidMethod(t *testing.T) {
	handler := AdministrationDeleteUsersHandler{}

	req, err := http.NewRequest("POST", "/locksmith/api/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestDeleteUserHTTPHandlesNoPayload(t *testing.T) {
	handler := AdministrationDeleteUsersHandler{}

	req, err := http.NewRequest("DELETE", "/locksmith/api/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteUserHTTPNonexistentUser(t *testing.T) {
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
					"sessions":     []interface{}{},
					"role":         "user",
					"customObject": "hello",
				},
			},
		},
	}

	handler := AdministrationDeleteUsersHandler{}

	payload := `{"username": "random"}`

	req, err := http.NewRequest("DELETE", "/locksmith/api/list", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusNotFound)
	}
}

func TestDeleteUserHTTP(t *testing.T) {
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
					"sessions":     []interface{}{},
					"role":         "user",
					"customObject": "hello",
				},
			},
		},
	}

	handler := AdministrationDeleteUsersHandler{}

	payload := `{"username": "kenton"}`

	req, err := http.NewRequest("DELETE", "/locksmith/api/list", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusOK)
	}
}
