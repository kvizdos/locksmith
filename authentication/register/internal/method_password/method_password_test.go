package method_password

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPasswordRegistrationSessionLoadRequest(t *testing.T) {
	t.Parallel()

	handler := NewPasswordRegistration()
	session := handler.Session(nil)
	req := httptest.NewRequest("POST", "/register", strings.NewReader(`{"username":"alice","password":"correct horse","email":"alice@example.com","code":"invite","validationok":true}`))
	req.Header.Set("Content-Type", "application/json")
	if err := session.LoadRequest(req); err != nil {
		t.Fatalf("load request: %v", err)
	}

	got := session.RegistrationRequest()
	if got.Username != "alice" || got.Password != "correct horse" || got.Email != "alice@example.com" || got.Code != "invite" || !got.ValidationOK {
		t.Fatalf("unexpected registration request: %+v", got)
	}
	if got.Hint != nil {
		t.Fatalf("password registration should not set hint: %+v", got.Hint)
	}
}

func TestPasswordRegistrationSessionRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	handler := NewPasswordRegistration()
	session := handler.Session(nil)
	req := httptest.NewRequest("POST", "/register", strings.NewReader(`{"username":"alice","password":"correct horse","role":"admin"}`))

	if err := session.LoadRequest(req); err == nil {
		t.Fatal("expected unknown role field to be rejected")
	}
}
