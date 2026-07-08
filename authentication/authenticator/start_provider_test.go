package authenticator

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/database"
)

// ── mockBeginnableHandler ─────────────────────────────────────────────────────

// mockBeginnableHandler is a Handler that also implements Beginnable.
type mockBeginnableHandler struct {
	name        string
	beginErr    error
	beginCalled bool
	beginWrite  func(w http.ResponseWriter)
}

func (h *mockBeginnableHandler) CanHandle(r *http.Request) bool { return true }
func (h *mockBeginnableHandler) Name() string                   { return h.name }
func (h *mockBeginnableHandler) Passwordless() bool             { return true }
func (h *mockBeginnableHandler) Session(db database.DatabaseAccessor) authenticator_domain.Session {
	return nil
}
func (h *mockBeginnableHandler) Begin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	h.beginCalled = true
	if h.beginWrite != nil {
		h.beginWrite(w)
	}
	return h.beginErr
}

var _ authenticator_domain.Handler = (*mockBeginnableHandler)(nil)
var _ authenticator_domain.Beginnable = (*mockBeginnableHandler)(nil)

// mockNonBeginnableHandler is a Handler that does NOT implement Beginnable.
type mockNonBeginnableHandler struct {
	name string
}

func (h *mockNonBeginnableHandler) CanHandle(r *http.Request) bool { return true }
func (h *mockNonBeginnableHandler) Name() string                   { return h.name }
func (h *mockNonBeginnableHandler) Passwordless() bool             { return false }
func (h *mockNonBeginnableHandler) Session(db database.DatabaseAccessor) authenticator_domain.Session {
	return nil
}

var _ authenticator_domain.Handler = (*mockNonBeginnableHandler)(nil)

// ── ServeProviderStartAPI tests ───────────────────────────────────────────────

func TestServeProviderStartAPI_NoHandlerMatch_Returns500(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	// Request for an unknown provider.
	req := httptest.NewRequest(http.MethodGet, "/api/start/unknown", nil)
	req.SetPathValue("provider", "unknown-provider")
	rr := httptest.NewRecorder()

	a.ServeProviderStartAPI(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 when no handler matches, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestServeProviderStartAPI_HandlerNotBeginnabe_Returns500(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	nonBeginnable := &mockNonBeginnableHandler{name: "non-beginnable"}
	a := newTestAuthorizer(db, WithMethods(nonBeginnable))

	req := httptest.NewRequest(http.MethodGet, "/api/start/non-beginnable", nil)
	req.SetPathValue("provider", "non-beginnable")
	rr := httptest.NewRecorder()

	a.ServeProviderStartAPI(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for non-Beginnable handler, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestServeProviderStartAPI_BeginReturnsErrMethodNotSupported_Returns405(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	handler := &mockBeginnableHandler{
		name:     "my-oidc",
		beginErr: errors.Join(errors.New("unsupported method: DELETE"), authenticator_domain.ErrMethodNotSupported),
	}
	a := newTestAuthorizer(db, WithMethods(handler))

	req := httptest.NewRequest(http.MethodDelete, "/api/start/my-oidc", nil)
	req.SetPathValue("provider", "my-oidc")
	rr := httptest.NewRecorder()

	a.ServeProviderStartAPI(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for ErrMethodNotSupported, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestServeProviderStartAPI_BeginReturnsGenericError_Returns500(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	handler := &mockBeginnableHandler{
		name:     "my-oidc",
		beginErr: errors.New("something went wrong internally"),
	}
	a := newTestAuthorizer(db, WithMethods(handler))

	req := httptest.NewRequest(http.MethodGet, "/api/start/my-oidc", nil)
	req.SetPathValue("provider", "my-oidc")
	rr := httptest.NewRecorder()

	a.ServeProviderStartAPI(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for generic Begin error, got %d — body: %s", rr.Code, rr.Body.String())
	}
}

func TestServeProviderStartAPI_BeginSucceeds_WritesResponse(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)
	handler := &mockBeginnableHandler{
		name:     "my-oidc",
		beginErr: nil,
		beginWrite: func(w http.ResponseWriter) {
			w.Header().Set("Location", "https://provider.example.com/auth")
			w.WriteHeader(http.StatusFound)
		},
	}
	a := newTestAuthorizer(db, WithMethods(handler))

	req := httptest.NewRequest(http.MethodGet, "/api/start/my-oidc", nil)
	req.SetPathValue("provider", "my-oidc")
	rr := httptest.NewRecorder()

	a.ServeProviderStartAPI(rr, req)

	if !handler.beginCalled {
		t.Error("expected Begin to be called")
	}
	if rr.Code != http.StatusFound {
		t.Errorf("expected 302 redirect from Begin, got %d", rr.Code)
	}
	if rr.Header().Get("Location") != "https://provider.example.com/auth" {
		t.Errorf("expected redirect location, got %q", rr.Header().Get("Location"))
	}
}

func TestServeProviderStartAPI_UsesPathValueToRouteHandler(t *testing.T) {
	t.Parallel()
	db := newTestDB(nil)

	called := false
	handler := &mockBeginnableHandler{
		name:     "specific-provider",
		beginErr: nil,
		beginWrite: func(w http.ResponseWriter) {
			called = true
			w.WriteHeader(http.StatusOK)
		},
	}
	a := newTestAuthorizer(db, WithMethods(handler))

	req := httptest.NewRequest(http.MethodGet, "/api/start/specific-provider", nil)
	req.SetPathValue("provider", "specific-provider")
	rr := httptest.NewRecorder()

	a.ServeProviderStartAPI(rr, req)

	if !called {
		t.Error("expected correct handler's Begin to be called via path value routing")
	}
}
