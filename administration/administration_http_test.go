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

var usersList []users.LocksmithUser

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

	user := usersList[0]

	if user.PasswordInfo.Password != "" {
		t.Errorf("SECURITY VULN: received a password value on listing")
		return
	}

	if len(user.PasswordSessions) != 0 {
		t.Errorf("SECURITY VULN: password sessions are getting sent, expected 0 received %d", len(user.PasswordSessions))
		return
	}
}
