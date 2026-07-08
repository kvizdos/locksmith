package authenticator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/database"
)

// ── getHandler tests ──────────────────────────────────────────────────────────

func TestGetHandler_PathValue_Matches(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodGet, "/api/login/password", nil)
	req.SetPathValue("provider", "password")

	h, err := a.getHandler(req)
	if err != nil {
		t.Fatalf("expected handler, got error: %v", err)
	}
	if h.Name() != "password" {
		t.Errorf("expected 'password', got %q", h.Name())
	}
}

func TestGetHandler_XHandlerHint_Matches(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Header.Set("X-Handler-Hint", "password")

	h, err := a.getHandler(req)
	if err != nil {
		t.Fatalf("expected handler, got error: %v", err)
	}
	if h.Name() != "password" {
		t.Errorf("expected 'password', got %q", h.Name())
	}
}

func TestGetHandler_CanHandle_Fallback(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	// No path value, no header — password handler's CanHandle checks for POST.
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)

	h, err := a.getHandler(req)
	if err != nil {
		t.Fatalf("expected handler via CanHandle, got error: %v", err)
	}
	if h.Name() != "password" {
		t.Errorf("expected 'password', got %q", h.Name())
	}
}

func TestGetHandler_NoMatch_ReturnsErrUnhandleableRequest(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	// GET with no hint — password handler only accepts POST via CanHandle.
	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)

	_, err := a.getHandler(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != authenticator_domain.ErrUnhandleableRequest {
		t.Errorf("expected ErrUnhandleableRequest, got %v", err)
	}
}

func TestGetHandler_PathValue_TakesPriorityOverHint(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)

	// Register two handlers with distinct names using a custom mock.
	secondHandler := &namedHandler{name: "second", canHandleResult: true}
	a := newTestAuthorizer(db, WithMethods(secondHandler))

	// Path value says "password", header says "second".
	req := httptest.NewRequest(http.MethodPost, "/api/login/password", nil)
	req.SetPathValue("provider", "password")
	req.Header.Set("X-Handler-Hint", "second")

	h, err := a.getHandler(req)
	if err != nil {
		t.Fatalf("expected handler, got error: %v", err)
	}
	// Path value wins: should route to "password".
	if h.Name() != "password" {
		t.Errorf("expected 'password' (path value wins), got %q", h.Name())
	}
}

func TestGetHandler_UnknownPathValue_ReturnsError(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodGet, "/api/login/nonexistent", nil)
	req.SetPathValue("provider", "nonexistent")

	_, err := a.getHandler(req)
	if err == nil {
		t.Fatal("expected error for unknown provider path, got nil")
	}
	if err != authenticator_domain.ErrUnhandleableRequest {
		t.Errorf("expected ErrUnhandleableRequest, got %v", err)
	}
}

func TestGetHandler_UnknownHint_ReturnsError(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Header.Set("X-Handler-Hint", "does-not-exist")

	_, err := a.getHandler(req)
	if err == nil {
		t.Fatal("expected error for unknown handler hint, got nil")
	}
	if err != authenticator_domain.ErrUnhandleableRequest {
		t.Errorf("expected ErrUnhandleableRequest, got %v", err)
	}
}

// ── namedHandler is a minimal Handler for testing getHandler routing ──────────

type namedHandler struct {
	name            string
	canHandleResult bool
}

func (h *namedHandler) CanHandle(r *http.Request) bool { return h.canHandleResult }
func (h *namedHandler) Name() string                   { return h.name }
func (h *namedHandler) Passwordless() bool             { return false }
func (h *namedHandler) Session(db database.DatabaseAccessor) authenticator_domain.Session {
	return nil
}

var _ authenticator_domain.Handler = (*namedHandler)(nil)
