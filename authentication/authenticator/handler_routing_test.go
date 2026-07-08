package authenticator

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerRouting_XHandlerHint(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db) // has "password" handler

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Header.Set("X-Handler-Hint", "password")

	handler, err := a.getHandler(req)
	if err != nil {
		t.Fatalf("expected handler, got error: %v", err)
	}
	if handler.Name() != "password" {
		t.Errorf("expected 'password', got %q", handler.Name())
	}
}

func TestHandlerRouting_PathValue(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodGet, "/api/login/password", nil)
	// Simulate path value routing
	req.SetPathValue("provider", "password")

	handler, err := a.getHandler(req)
	if err != nil {
		t.Fatalf("expected handler, got error: %v", err)
	}
	if handler.Name() != "password" {
		t.Errorf("expected 'password', got %q", handler.Name())
	}
}

func TestHandlerRouting_CanHandle_Fallthrough(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	// No hint, no path — password handler's CanHandle checks for POST
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)

	handler, err := a.getHandler(req)
	if err != nil {
		t.Fatalf("expected handler via CanHandle, got error: %v", err)
	}
	if handler.Name() != "password" {
		t.Errorf("expected 'password', got %q", handler.Name())
	}
}

func TestHandlerRouting_NoMatch_Returns500(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	// GET with no matching handler (password requires POST)
	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, req)

	if rr.Code != http.StatusMethodNotAllowed && rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 405 or 500 for no matching handler, got %d", rr.Code)
	}
}

func TestHandlerRouting_UnknownHint_Returns500(t *testing.T) {
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Header.Set("X-Handler-Hint", "nonexistent-handler")
	rr := httptest.NewRecorder()
	a.ServeLoginAPI(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for unknown handler hint, got %d", rr.Code)
	}
}
