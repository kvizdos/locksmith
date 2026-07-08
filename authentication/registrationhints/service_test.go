package registrationhints

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication/signing"
)

func TestServiceCreateAndParse(t *testing.T) {
	t.Parallel()

	signer, err := signing.CreateSigningPackage()
	if err != nil {
		t.Fatalf("create signer: %v", err)
	}

	tests := []struct {
		name          string
		hint          Hint
		mutate        func(*testing.T, Service, string) string
		wantCreateErr error
		wantParseErr  error
	}{
		{
			name: "valid hint parses successfully",
			hint: Hint{ProviderName: "google", Email: "user@example.com", Issuer: "https://issuer.example.com", Subject: "sub", Rosterable: true},
		},
		{
			name:          "non-rosterable hint fails closed",
			hint:          Hint{ProviderName: "google", Email: "user@example.com", Issuer: "https://issuer.example.com", Subject: "sub", Rosterable: false},
			wantCreateErr: ErrNotRosterable,
		},
		{
			name:          "missing issuer or subject fails closed",
			hint:          Hint{ProviderName: "google", Email: "user@example.com", Rosterable: true},
			wantCreateErr: ErrMissingIdentity,
		},
		{
			name: "wrong audience fails closed",
			hint: Hint{ProviderName: "google", Email: "user@example.com", Issuer: "https://issuer.example.com", Subject: "sub", Rosterable: true},
			mutate: func(t *testing.T, svc Service, token string) string {
				wrong := Hint{ProviderName: "google", Email: "user@example.com", Issuer: "https://issuer.example.com", Subject: "sub", Rosterable: true}
				wrong.Audience = []string{"wrong"}
				wrong.ExpiresAt = nil
				badToken, err := signer.CreateJWT(&wrong)
				if err != nil {
					t.Fatalf("create wrong-audience token: %v", err)
				}
				return badToken
			},
			wantParseErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := Service{Signer: signer, TTL: time.Minute}
			token, err := svc.Create(tt.hint)
			if tt.wantCreateErr != nil {
				if !errors.Is(err, tt.wantCreateErr) {
					t.Fatalf("expected %v, got %v", tt.wantCreateErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("create hint: %v", err)
			}
			if tt.mutate != nil {
				token = tt.mutate(t, svc, token)
			}

			parsed, err := svc.Parse(token)
			if tt.wantParseErr != nil {
				if !errors.Is(err, tt.wantParseErr) {
					t.Fatalf("expected %v, got %v", tt.wantParseErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("parse hint: %v", err)
			}
			if parsed.ID == "" || parsed.IssuedAt == nil || parsed.ExpiresAt == nil {
				t.Fatalf("expected id, issued-at, and expiration to be set: %+v", parsed)
			}
			if len(parsed.Audience) != 1 || parsed.Audience[0] != Audience {
				t.Fatalf("expected audience %q, got %v", Audience, parsed.Audience)
			}
		})
	}
}

func TestServiceParseExpiredHintFailsClosed(t *testing.T) {
	t.Parallel()

	signer, err := signing.CreateSigningPackage()
	if err != nil {
		t.Fatalf("create signer: %v", err)
	}
	svc := Service{Signer: signer, TTL: time.Minute}

	expired := Hint{ProviderName: "google", Email: "user@example.com", Issuer: "https://issuer.example.com", Subject: "sub", Rosterable: true}
	expired.Audience = []string{Audience}
	expired.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Second))
	token, err := signer.CreateJWT(&expired)
	if err != nil {
		t.Fatalf("create expired hint: %v", err)
	}

	if _, err := svc.Parse(token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected invalid token error, got %v", err)
	}
}

func TestCookieHelpers(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	SetCookieWithTTL(rr, "token-value", 2*time.Minute)

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != CookieName || cookie.Value != "token-value" {
		t.Fatalf("unexpected cookie identity: %+v", cookie)
	}
	if !cookie.HttpOnly || !cookie.Secure {
		t.Fatalf("expected secure httponly cookie: %+v", cookie)
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
	if cookie.Path != "/api/register" {
		t.Fatalf("expected /api/register path, got %q", cookie.Path)
	}
	if cookie.MaxAge != 120 || cookie.Expires.IsZero() {
		t.Fatalf("expected max-age/expires aligned to ttl: %+v", cookie)
	}
}

func TestFromRequest(t *testing.T) {
	t.Parallel()

	signer, err := signing.CreateSigningPackage()
	if err != nil {
		t.Fatalf("create signer: %v", err)
	}
	svc := Service{Signer: signer, TTL: time.Minute}
	token, err := svc.Create(Hint{ProviderName: "google", Email: "user@example.com", Issuer: "https://issuer.example.com", Subject: "sub", Rosterable: true})
	if err != nil {
		t.Fatalf("create hint: %v", err)
	}

	tests := []struct {
		name    string
		request func() *http.Request
		wantErr error
	}{
		{
			name: "loads from header",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/register", nil)
				r.Header.Set("X-Registration-Hint", token)
				return r
			},
		},
		{
			name: "loads from cookie",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/register", nil)
				r.AddCookie(&http.Cookie{Name: CookieName, Value: token})
				return r
			},
		},
		{
			name:    "missing token fails closed",
			request: func() *http.Request { return httptest.NewRequest(http.MethodGet, "/register", nil) },
			wantErr: ErrMissingToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := svc.FromRequest(tt.request())
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
