package administration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/users"
)

func TestMain(m *testing.M) {
	roles.AVAILABLE_ROLES = map[string]roles.RoleInfo{
		"admin": {
			BackendPermissions: []string{"view.admin", "user.delete.self", "user.delete.other"},
		},
		"user": {
			BackendPermissions: []string{"view.admin", "user.delete.self"},
		},
		"norights": {},
	}

	m.Run()

	roles.AVAILABLE_ROLES = map[string]roles.RoleInfo{}
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
					"email":    "email@email.com",
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

func TestListUsersSpecificRole(t *testing.T) {
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
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role":     "user",
					"sessions": sessions,
				},

				"a8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "bob",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role":     "admin",
					"sessions": sessions,
				},
			},
		},
	}

	handler := AdministrationListUsersHandler{}

	req, err := http.NewRequest("GET", "/locksmith/api/list?role=user", nil)
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

	if usersList[0].Role != "user" {
		t.Errorf("error reading correct role, expected user got %s", usersList[0].Role)
	}
}

func TestListUsersMultipleRoles(t *testing.T) {
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
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role":     "user",
					"sessions": sessions,
				},

				"a8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "bob",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role":     "admin",
					"sessions": sessions,
				},

				"b8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "z8531661-22a7-493f-b228-028842e09a05",
					"username": "james",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"role":     "other",
					"sessions": sessions,
				},
			},
		},
	}

	handler := AdministrationListUsersHandler{}

	req, err := http.NewRequest("GET", "/locksmith/api/list?role=user&role=admin", nil)
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

	if len(usersList) != 2 {
		t.Errorf("received unexpected number of users, expected 1 got %d", len(usersList))
		return
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
					"email":    "email@email.com",
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
					"email":    "email@email.com",
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

	withUser := users.LocksmithUser{
		Username: "kenton",
		Role:     "admin",
	}
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusNotFound)
	}
}

func TestDeleteUserSelfHTTP(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
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

	withUser := users.LocksmithUser{
		Username: "kenton",
		Role:     "user",
	}
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusOK)
	}
}

func TestDeleteUserSelfUnauthorizedHTTP(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"sessions":     []interface{}{},
					"role":         "norights",
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

	withUser := users.LocksmithUser{
		Username: "kenton",
		Role:     "norights",
	}
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusUnauthorized)
	}
}

func TestDeleteUserHTTPUnauthorizedCantDeleteOtherUsers(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
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

	payload := `{"username": "bob"}`

	req, err := http.NewRequest("DELETE", "/locksmith/api/list", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	withUser := users.LocksmithUser{
		Username: "kenton",
		Role:     "user",
	}
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusUnauthorized)
	}
}

func TestDeleteUserHTTPAuthorizedCanDeleteOtherUsers(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"password": authentication.PasswordInfo{
						Password: "testpassword",
						Salt:     "testsalt",
					},
					"sessions":     []interface{}{},
					"role":         "admin",
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

	withUser := users.LocksmithUser{
		Username: "james",
		Role:     "admin",
	}
	req = req.WithContext(context.WithValue(req.Context(), "authUser", withUser))
	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code got %v, want %v", status, http.StatusOK)
	}
}
