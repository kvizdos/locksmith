package method_oidc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
	"github.com/kvizdos/locksmith/database"
	"golang.org/x/oauth2"
)

func TestResolveIdentityCredentialFlowSetsVerifiedIdentity(t *testing.T) {
	t.Parallel()

	provider := newTestOIDCProvider(t)
	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		Verifier:     provider.verifier,
	})
	session.flow = flowCredential
	session.untrustedCredentialToken = provider.idToken(t, testOIDCClaims{
		Subject:       "subject-123",
		Email:         "kenton@example.com",
		EmailVerified: true,
		Name:          "Kenton V",
	})

	if err := session.ResolveIdentity(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.GetSubject() != "subject-123" {
		t.Fatalf("expected subject to be set, got %q", session.GetSubject())
	}
	if session.GetIssuer() != provider.issuer {
		t.Fatalf("expected issuer to be set, got %q", session.GetIssuer())
	}
	if session.GetEmail() != "kenton@example.com" {
		t.Fatalf("expected verified email to be set, got %q", session.GetEmail())
	}
	if session.displayName != "Kenton V" {
		t.Fatalf("expected display name to be set, got %q", session.displayName)
	}
	if session.authoritativeToken == nil {
		t.Fatal("expected authoritative token to be stored")
	}
}

func TestResolveIdentityCredentialFlowDoesNotSetUnverifiedEmail(t *testing.T) {
	t.Parallel()

	provider := newTestOIDCProvider(t)
	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		Verifier:     provider.verifier,
	})
	session.flow = flowCredential
	session.untrustedCredentialToken = provider.idToken(t, testOIDCClaims{
		Subject:       "subject-123",
		Email:         "unverified@example.com",
		EmailVerified: false,
		Name:          "Unverified User",
	})

	if err := session.ResolveIdentity(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.GetEmail() != "" {
		t.Fatalf("expected unverified email not to be set, got %q", session.GetEmail())
	}
	if session.displayName != "Unverified User" {
		t.Fatalf("expected display name to be set, got %q", session.displayName)
	}
}

func TestResolveIdentityCodeFlowExchangesCodeWithPKCEVerifier(t *testing.T) {
	t.Parallel()

	provider := newTestOIDCProvider(t)
	provider.tokenResponseIDToken = provider.idToken(t, testOIDCClaims{
		Subject:       "code-subject",
		Email:         "code@example.com",
		EmailVerified: true,
		Name:          "Code User",
	})
	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		Verifier:     provider.verifier,
		OauthConfig:  provider.oauthConfig,
	})
	session.flow = flowCode
	session.untrustedParsedCode = "auth-code"
	session.pkceVerifier = "pkce-verifier"

	if err := session.ResolveIdentity(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.seenCode != "auth-code" {
		t.Fatalf("expected code to be exchanged, got %q", provider.seenCode)
	}
	if provider.seenCodeVerifier != "pkce-verifier" {
		t.Fatalf("expected pkce verifier to be sent, got %q", provider.seenCodeVerifier)
	}
	if session.GetSubject() != "code-subject" {
		t.Fatalf("expected code flow subject to be set, got %q", session.GetSubject())
	}
	if session.GetEmail() != "code@example.com" {
		t.Fatalf("expected code flow email to be set, got %q", session.GetEmail())
	}
}

func TestResolveIdentityCodeFlowMissingIDToken(t *testing.T) {
	t.Parallel()

	provider := newTestOIDCProvider(t)
	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		Verifier:     provider.verifier,
		OauthConfig:  provider.oauthConfig,
	})
	session.flow = flowCode
	session.untrustedParsedCode = "auth-code"
	session.pkceVerifier = "pkce-verifier"

	err := session.ResolveIdentity(context.Background())
	if err == nil {
		t.Fatal("expected error for missing id_token")
	}
	if !strings.Contains(err.Error(), "missing id_token") {
		t.Fatalf("expected missing id_token error, got %v", err)
	}
}

func TestResolveIdentityReturnsVerificationError(t *testing.T) {
	t.Parallel()

	provider := newTestOIDCProvider(t)
	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{
		ProviderName: "google",
		Verifier:     provider.verifier,
	})
	session.flow = flowCredential
	session.untrustedCredentialToken = "not-a-jwt"

	err := session.ResolveIdentity(context.Background())
	if err == nil {
		t.Fatal("expected verification error")
	}
	if !strings.Contains(err.Error(), "failed to verify id_token") {
		t.Fatalf("expected verification error wrapper, got %v", err)
	}
}

func TestResolveIdentityUnsupportedFlow(t *testing.T) {
	t.Parallel()

	session := newOIDCValidationSession(database.TestDatabase{}, authenticator_methods.OIDCValidatorOptions{ProviderName: "google"})

	err := session.ResolveIdentity(context.Background())
	if err == nil {
		t.Fatal("expected unsupported flow error")
	}
	if !strings.Contains(err.Error(), "unsupported flow") {
		t.Fatalf("expected unsupported flow error, got %v", err)
	}
}

type testOIDCProvider struct {
	issuer               string
	server               *httptest.Server
	privateKey           *rsa.PrivateKey
	verifier             *oidc.IDTokenVerifier
	oauthConfig          *oauth2.Config
	tokenResponseIDToken string
	seenCode             string
	seenCodeVerifier     string
}

type testOIDCClaims struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
}

func newTestOIDCProvider(t *testing.T) *testOIDCProvider {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}

	provider := &testOIDCProvider{privateKey: privateKey}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	provider.server = server
	provider.issuer = server.URL

	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"issuer":                 provider.issuer,
			"authorization_endpoint": provider.issuer + "/auth",
			"token_endpoint":         provider.issuer + "/token",
			"jwks_uri":               provider.issuer + "/keys",
		})
	})
	mux.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"keys": []map[string]any{provider.jwk()},
		})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		provider.seenCode = r.Form.Get("code")
		provider.seenCodeVerifier = r.Form.Get("code_verifier")

		response := map[string]any{
			"access_token": "access-token",
			"token_type":   "Bearer",
		}
		if provider.tokenResponseIDToken != "" {
			response["id_token"] = provider.tokenResponseIDToken
		}
		writeJSON(t, w, response)
	})

	discoveredProvider, err := oidc.NewProvider(context.Background(), provider.issuer)
	if err != nil {
		t.Fatalf("create oidc provider: %v", err)
	}
	provider.verifier = discoveredProvider.Verifier(&oidc.Config{ClientID: "client-id"})
	provider.oauthConfig = &oauth2.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Endpoint:     discoveredProvider.Endpoint(),
		RedirectURL:  "https://example.com/api/login",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return provider
}

func (p *testOIDCProvider) idToken(t *testing.T, claims testOIDCClaims) string {
	t.Helper()

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":            p.issuer,
		"sub":            claims.Subject,
		"aud":            "client-id",
		"exp":            now.Add(time.Hour).Unix(),
		"iat":            now.Unix(),
		"email":          claims.Email,
		"email_verified": claims.EmailVerified,
		"name":           claims.Name,
	})
	token.Header["kid"] = "test-key"

	signed, err := token.SignedString(p.privateKey)
	if err != nil {
		t.Fatalf("sign id token: %v", err)
	}
	return signed
}

func (p *testOIDCProvider) jwk() map[string]any {
	publicKey := p.privateKey.Public().(*rsa.PublicKey)
	return map[string]any{
		"kty": "RSA",
		"use": "sig",
		"kid": "test-key",
		"alg": "RS256",
		"n":   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil && !errors.Is(err, http.ErrAbortHandler) {
		t.Fatalf("write json response: %v", err)
	}
}
