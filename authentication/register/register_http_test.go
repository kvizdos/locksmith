package register

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/users"
)

func TestMain(m *testing.M) {
	roles.AVAILABLE_ROLES = map[string][]string{
		"admin": {
			"view.admin",
			"user.delete.self",
		},
	}

	m.Run()

	roles.AVAILABLE_ROLES = map[string][]string{}
}

func TestRegistrationHandlerMissingRole(t *testing.T) {
	handler := RegistrationHandler{}

	// Test Missing Username
	payload := `{"username": "kvizdos", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestRegistrationHandlerInvalidRole(t *testing.T) {
	handler := RegistrationHandler{
		DefaultRoleName: "not-set",
	}

	// Test Missing Username
	payload := `{"username": "kvizdos", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestRegistrationHandlerFailedUnmarshal(t *testing.T) {
	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	// Test Missing Username
	payload := `{"password" "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}

	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "could not unmarshal" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerMissingBodyParams(t *testing.T) {
	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	// Test Missing Username
	payload := `{"password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code (missing username): got %v, want %v", status, http.StatusBadRequest)
	}

	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "missing fields" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}

	// Test Missing Password
	payload = `{"username": "kenton"}`

	req, err = http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code (missing password): got %v, want %v", status, http.StatusBadRequest)
	}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "missing fields" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerUsernameTaken(t *testing.T) {
	testDb := database.TestDatabase{
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

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "kenton", "password": "password123", "email": "email@email.com"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusConflict)
	}
	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "taken" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerEmailTaken(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "kvizdos", "password": "password123", "email": "email@email.com"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusConflict)
	}
	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "taken" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerEmailInvalid(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton",
					"email":    "email@email.com",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "kvizdos", "password": "password123", "email": "email@ema"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
	}
	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "invalid email" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerSuccess(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "kenton", "password": "password123", "email": "email@email.com"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
		return
	}

	newUser, found := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	if !found {
		t.Errorf("new user did not get added to database.")
		return
	}

	user := newUser.(map[string]interface{})

	var lsu users.LocksmithUserInterface
	lsu = users.LocksmithUser{}
	lsu.ReadFromMap(&lsu, user)
	passwordInfo := lsu.GetPasswordInfo()
	if passwordInfo.Password == "password123" {
		t.Errorf("PASSWORD NOT ENCRYPTED!")
		return
	}

	if len(passwordInfo.Salt) != 32 {
		t.Errorf("Salt not a correct length, expected %d, got %d", 32, len(passwordInfo.Salt))
	}

	if user["email"] != "email@email.com" {
		t.Error("did not receive correct email on registration")
	}
}

func TestRegistrationHandlerDoesNotMeetLengthRequirement(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
	}

	payload := `{"username": "kenton", "password": "1234567", "email": "email@email.com"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
		return
	}

	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "password too short" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerDoesMeetLengthRequirementSuccess(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
	}

	payload := `{"username": "kenton", "password": "12345678", "email": "email@email.com"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("1: unexpected status code: got %v, want %v", status, http.StatusOK)
		return
	}

	// Check w/ extra long password
	handler = RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
	}

	payload = `{"username": "kenton3", "password": "ahisfignfdiofjfdifdfiidofgiohdfgiohdfgoidfhgiodfngonfiogfngoidf", "email": "email3@email.com"}`

	req, err = http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("2: unexpected status code: got %v, want %v - %s", status, http.StatusOK, rr.Body.String())
		return
	}
}

func TestRegistrationHandlerHIBPStrictPasswordExistsInBreach(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
		HIBP: hibp.HIBPSettings{
			Enabled:     true,
			AppName:     "Locksmith Integration Tests",
			Enforcement: hibp.STRICT,
		},
	}

	payload := `{"username": "kenton", "password": "password123", "email": "email@email.com"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusConflict)
		return
	}

	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "password pwned" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}

	if !regErr.PwnStatus {
		t.Error("should have TRUE pwn status")
	}
}

func TestRegistrationHandlerHIBPStrictPasswordExistsInBreachCantBeBypassedWithBody(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
		HIBP: hibp.HIBPSettings{
			Enabled:     true,
			AppName:     "Locksmith Integration Tests",
			Enforcement: hibp.STRICT,
		},
	}

	payload := `{"username": "kenton", "password": "password123", "email": "email@email.com", "pwnok": true}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusConflict)
		return
	}

	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "password pwned" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}

	if !regErr.PwnStatus {
		t.Error("should have TRUE pwn status")
	}
}

func TestRegistrationHandlerHIBPStrictPasswordDoesNotExistsInBreach(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
		HIBP: hibp.HIBPSettings{
			Enabled:     true,
			AppName:     "Locksmith Integration Tests",
			Enforcement: hibp.STRICT,
		},
	}

	pass, _ := authentication.GenerateRandomString(128)
	payload := fmt.Sprintf(`{"username": "kenton", "password": "%s", "email": "email@email.com"}`, pass)

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
		return
	}

	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}

	if regErr.PwnStatus {
		t.Errorf("password should not be pwned.")
	}
}

func TestRegistrationHandlerHIBPLooseDoesExistsInBreachSoFailsWithoutBodyConfirmation(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
		HIBP: hibp.HIBPSettings{
			Enabled:     true,
			AppName:     "Locksmith Integration Tests",
			Enforcement: hibp.LOOSE,
		},
	}

	payload := fmt.Sprintf(`{"username": "kenton", "password": "%s", "email": "email@email.com"}`, "password123")

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusConflict)
		return
	}

	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "password pwned" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
	if !regErr.PwnStatus {
		t.Errorf("pwnStatus should be TRUE")
	}
}

func TestRegistrationHandlerHIBPLooseDoesExistsInBreachPassesWithBodyConfirmation(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:          "admin",
		MinimumLengthRequirement: 8,
		HIBP: hibp.HIBPSettings{
			Enabled:     true,
			AppName:     "Locksmith Integration Tests",
			Enforcement: hibp.LOOSE,
		},
	}

	payload := fmt.Sprintf(`{"username": "kenton", "password": "%s", "email": "email@email.com", "pwnok": true}`, "password123")

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
		return
	}
}

type customUser struct {
	users.LocksmithUser

	CustomObject string `json:"customObject"`
}

func (c customUser) ReadFromMap(writeTo *users.LocksmithUserInterface, u map[string]interface{}) {
	// Load initial locksmith data
	var user users.LocksmithUserInterface
	c.LocksmithUser.ReadFromMap(&user, u)
	lsu := user.(users.LocksmithUser)

	converted := customUser{
		LocksmithUser: lsu,
	}

	converted.CustomObject = u["customObject"].(string)

	*writeTo = converted
}

func (c customUser) ToMap() map[string]interface{} {
	lsu := c.LocksmithUser.ToMap()

	lsu["customObject"] = c.CustomObject

	return lsu
}

func TestRegistrationHandlerSuccessCustomUser(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
		ConfigureCustomUser: func(lui users.LocksmithUser, db database.DatabaseAccessor) users.LocksmithUserInterface {
			user := customUser{
				LocksmithUser: lui,
			}
			user.CustomObject = "Hello World"
			return user
		},
	}

	payload := `{"username": "kenton", "password": "password123", "email": "email@email.com"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
		return
	}

	newUser, found := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
	})

	if !found {
		t.Errorf("new user did not get added to database.")
		return
	}

	user := newUser.(map[string]interface{})

	var lsu users.LocksmithUserInterface
	lsu = customUser{}
	lsu.ReadFromMap(&lsu, user)

	if user["email"] != "email@email.com" {
		t.Error("did not receive correct email on registration")
	}

	if user["customObject"] != "Hello World" {
		t.Errorf("did not receive custom object on registration: %s", user["customObject"])
	}
}

// Tests to see if the username does not match [A-z0-9]
func TestRegistrationHandlerInvalidUsername(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName: "admin",
	}

	payload := `{"username": "<i want xss>!", "password": "password123"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
		return
	}
	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "illegal username characters" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

// Test invite code stuff..
func TestRegistrationHandlerInvalidInviteCode(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:           "admin",
		DisablePublicRegistration: true,
	}

	payload := `{"username": "kenton", "password": "password123", "email": "email@email.com", "code": "asdadsasd"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
		return
	}
	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "bad invite code" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerInviteCodeNotFound(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"code":    "09c47c463871acf1da8c5c6dbcc73723d35d1aa262f7cfb29bad901d9640cd03",
					"email":   "bob@bob.com",
					"role":    "user",
					"inviter": "a-uuid",
					"sentAt":  time.Now().Unix(),
					"userid":  "bobbb",
				},
			},
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:           "admin",
		DisablePublicRegistration: true,
	}

	payload := `{"username": "kenton", "password": "password123", "email": "notbob@bob.com", "code": "zyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
		return
	}
	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "invalid code" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerIncorrectEmail(t *testing.T) {
	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"code":    "09c47c463871acf1da8c5c6dbcc73723d35d1aa262f7cfb29bad901d9640cd03",
					"email":   "bob@bob.com",
					"role":    "user",
					"inviter": "a-uuid",
					"sentAt":  time.Now().Unix(),
					"userid":  "bobbb",
				},
			},
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:           "admin",
		DisablePublicRegistration: true,
	}

	payload := `{"username": "kenton", "password": "password123", "email": "notbob@bob.com", "code": "jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusBadRequest)
		return
	}
	regErr := &registrationResponse{}
	regErr.Unmarshal(rr.Body.Bytes())

	if regErr.Error != "invalid email" {
		t.Errorf("got unexpected error message: %s", regErr.Error)
	}
}

func TestRegistrationHandlerWithInviteSuccess(t *testing.T) {
	hasher := sha256.New()
	hasher.Write([]byte("jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M"))
	hashedCode := hasher.Sum(nil)

	testDb := database.TestDatabase{
		Tables: map[string]map[string]interface{}{
			"invites": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"code":    fmt.Sprintf("%x", hashedCode),
					"email":   "bob@bob.com",
					"role":    "admin",
					"inviter": "a-uuid",
					"sentAt":  time.Now().Unix(),
					"userid":  "this-is-the-user-id",
				},
			},
			"users": {
				"c8531661-22a7-493f-b228-028842e09a05": map[string]interface{}{
					"id":       "c8531661-22a7-493f-b228-028842e09a05",
					"username": "kenton2",
					"email":    "email@email.com2",
					"sessions": []interface{}{"abc"},
				},
			},
		},
	}

	handler := RegistrationHandler{
		DefaultRoleName:           "admin",
		DisablePublicRegistration: true,
	}

	payload := `{"username": "kenton", "password": "password123", "email": "bob@bob.com", "code": "jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M"}`

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
		return
	}

	_, correct := testDb.FindOne("users", map[string]interface{}{
		"username": "kenton",
		"role":     "admin",
	})

	if !correct {
		t.Errorf("could not find user in database.")
	}

	_, shouldBeFalse := testDb.FindOne("invites", map[string]interface{}{
		"code": fmt.Sprintf("%x", hashedCode),
	})

	if shouldBeFalse {
		t.Errorf("invite was not expired after usage.")
	}
}

// Test Email-Only Support
func TestRegistrationWithEmail(t *testing.T) {
	testCases := []struct {
		DbSeed           map[string]interface{}
		Username         string
		Password         string
		ExpectStatusCode int
		ExpectErrorMsg   string
	}{
		{
			// Should pass with valid email that doesn't exist
			DbSeed:           map[string]interface{}{},
			Username:         "kvizdos@gmail.com",
			Password:         "password123",
			ExpectStatusCode: 200,
			ExpectErrorMsg:   "",
		},
		{
			// Should fail with a non-email passed as username
			DbSeed:           map[string]interface{}{},
			Username:         "kvizdos",
			Password:         "password123",
			ExpectStatusCode: 400,
			ExpectErrorMsg:   "illegal username characters",
		},
		{
			// Test with a email that already exists
			DbSeed: map[string]interface{}{
				"rand-id": map[string]interface{}{
					"username": "kvizdos@gmail.com",
					"email":    "kvizdos@gmail.com",
				},
			},
			Username:         "kvizdos@gmail.com",
			Password:         "helloworld",
			ExpectStatusCode: 409,
			ExpectErrorMsg:   "taken",
		},
	}

	for _, test := range testCases {
		testDb := database.TestDatabase{
			Tables: map[string]map[string]interface{}{
				"users": test.DbSeed,
			},
		}

		handler := RegistrationHandler{
			DefaultRoleName: "admin",
			EmailAsUsername: true,
		}

		payload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, test.Username, test.Password)

		req, err := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), "database", testDb))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.ExpectStatusCode {
			t.Errorf("unexpected status code: got %v, want %v", status, test.ExpectStatusCode)
			return
		}

		if test.ExpectErrorMsg != "" {
			regErr := &registrationResponse{}
			regErr.Unmarshal(rr.Body.Bytes())

			if regErr.Error != test.ExpectErrorMsg {
				t.Errorf("got unexpected error message: %s", regErr.Error)
			}
		}
	}
}
