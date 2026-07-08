package register

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/events"
	"github.com/kvizdos/locksmith/authentication/register/register_methods"
	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/authentication/signing"
	"github.com/kvizdos/locksmith/authentication/textvalidation"
	"github.com/kvizdos/locksmith/authentication/verificationcodes"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/users"
)

type recordingTokenManager struct {
	createCalls int
	passCalls   int
	createdFor  users.LocksmithUserInterface
	passedToken *authenticator_domain.Token
	createErr   error
	passErr     error
}

func (m *recordingTokenManager) Read(r *http.Request) (*authenticator_domain.Token, error) {
	return nil, nil
}

func (m *recordingTokenManager) CreateAuthToken(user users.LocksmithUserInterface) (*authenticator_domain.Token, error) {
	m.createCalls++
	m.createdFor = user
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &authenticator_domain.Token{AuthToken: "session-token", ExpiresAt: time.Now().Add(time.Hour)}, nil
}

func (m *recordingTokenManager) PassToClient(w http.ResponseWriter, r *http.Request, token *authenticator_domain.Token) error {
	m.passCalls++
	m.passedToken = token
	if m.passErr != nil {
		return m.passErr
	}
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "client-session", HttpOnly: true, Secure: true, Path: "/"})
	return nil
}

func newTestHintServiceForHTTP(t *testing.T) registrationhints.Service {
	t.Helper()
	signer, err := signing.CreateSigningPackage()
	if err != nil {
		t.Fatalf("create signer: %v", err)
	}
	return registrationhints.Service{Signer: signer}
}

// newTestRegistrarWithHint builds a *registrar wired with both the hint and
// password methods, with hint checked first (it only matches requests
// carrying a hint cookie/header, whereas password matches any POST).
func newTestRegistrarWithHint(db database.DatabaseAccessor, hints registrationhints.Service, opts ...Option) *registrar {
	base := []Option{WithMethods(
		AllowMethodHint(register_methods.WithHintService(hints)),
		AllowMethodPassword(),
	)}
	return NewRegistrar(db, append(base, opts...)...)
}

type staticEmailValidator struct {
	result textvalidation.ValidationResult
	dym    *string
}

func (v staticEmailValidator) Validate(context.Context, string) (textvalidation.ValidationResultEvaluator, error) {
	return staticEmailValidationResult{result: v.result, dym: v.dym}, nil
}

type staticEmailValidationResult struct {
	result textvalidation.ValidationResult
	dym    *string
}

func (r staticEmailValidationResult) Result(bool) (*string, textvalidation.ValidationResult) {
	return r.dym, r.result
}

func (r staticEmailValidationResult) DebugPrint(string) {}

type recordingVerificationSender struct {
	calls    int
	forValue string
}

func (s *recordingVerificationSender) SendVerificationEmail(_ context.Context, userEmail string, _ string, _ string) error {
	s.calls++
	s.forValue = userEmail
	return nil
}

func newRegistrationTestDB(usersSeed map[string]interface{}) database.TestDatabase {
	if usersSeed == nil {
		usersSeed = map[string]interface{}{}
	}
	return database.TestDatabase{Tables: map[string]map[string]interface{}{"users": usersSeed}}
}

func performRegistrationRequest(t *testing.T, r *registrar, payload string) *httptest.ResponseRecorder {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeRegisterAPI(rr, req)
	return rr
}

type recordingEventBus struct {
	mu        sync.Mutex
	published []publishedEvent
}

type publishedEvent struct {
	ctx   context.Context
	event events.Envelope
}

func (b *recordingEventBus) Publish(ctx context.Context, event events.Envelope) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.published = append(b.published, publishedEvent{ctx: ctx, event: event})
	return nil
}

func (b *recordingEventBus) Subscribe(events.EventName, events.Handler) events.Subscription {
	return noopEventSubscription{}
}

type noopEventSubscription struct{}

func (noopEventSubscription) Unsubscribe() {}

func (b *recordingEventBus) eventsByName(name events.EventName) []events.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()
	var result []events.Envelope
	for _, published := range b.published {
		if published.event.Name == name {
			result = append(result, published.event)
		}
	}
	return result
}

func (b *recordingEventBus) singlePublish(t *testing.T, name events.EventName) publishedEvent {
	t.Helper()
	b.mu.Lock()
	defer b.mu.Unlock()
	var result []publishedEvent
	for _, published := range b.published {
		if published.event.Name == name {
			result = append(result, published)
		}
	}
	if len(result) != 1 {
		t.Fatalf("events named %q = %d, want 1", name, len(result))
	}
	return result[0]
}

func decodeRegistrationResponse(t *testing.T, rr *httptest.ResponseRecorder) registrationResponse {
	t.Helper()

	var res registrationResponse
	res.Unmarshal(rr.Body.Bytes())
	return res
}

func TestRegistrationHandlerPublishesRegistrationSucceededEvent(t *testing.T) {
	t.Parallel()

	bus := &recordingEventBus{}
	db := newRegistrationTestDB(nil)
	r := newTestRegistrar(db, WithDefaultRoleName("admin"), WithEventBus(bus))
	rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"email@example.com"}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	published := bus.singlePublish(t, events.EventRegistrationSucceeded)
	payload, ok := published.event.Payload.(events.RegistrationSucceededPayload)
	if !ok {
		t.Fatalf("payload type = %T, want events.RegistrationSucceededPayload", published.event.Payload)
	}
	if payload.UserID == "" {
		t.Fatal("payload.UserID should be populated")
	}
	if payload.Username != "kenton" {
		t.Fatalf("payload.Username = %q, want %q", payload.Username, "kenton")
	}
	if payload.Email != "email@example.com" {
		t.Fatalf("payload.Email = %q, want %q", payload.Email, "email@example.com")
	}
	if payload.Method != "password" || payload.Provider != "password" {
		t.Fatalf("method/provider = %q/%q, want password/password", payload.Method, payload.Provider)
	}
	if payload.InviteUsed {
		t.Fatal("payload.InviteUsed = true, want false")
	}
	if payload.AutoLoginIssued {
		t.Fatal("payload.AutoLoginIssued = true before auto-login is implemented")
	}
}

func TestRegistrationHandlerPublishesFailedEventWithoutSucceededEvent(t *testing.T) {
	t.Parallel()

	bus := &recordingEventBus{}
	db := newRegistrationTestDB(nil)
	r := newTestRegistrar(db, WithDefaultRoleName("admin"), WithEventBus(bus))
	rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"not-an-email"}`)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	if got := bus.eventsByName(events.EventRegistrationSucceeded); len(got) != 0 {
		t.Fatalf("success events = %d, want 0", len(got))
	}

	published := bus.singlePublish(t, events.EventRegistrationFailed)
	payload, ok := published.event.Payload.(events.RegistrationFailedPayload)
	if !ok {
		t.Fatalf("payload type = %T, want events.RegistrationFailedPayload", published.event.Payload)
	}
	if payload.Reason != "invalid_email" {
		t.Fatalf("payload.Reason = %q, want %q", payload.Reason, "invalid_email")
	}
	if strings.Contains(fmt.Sprintf("%+v", payload), "password123") {
		t.Fatal("failed registration event leaked plaintext password")
	}
}

func TestRegistrationHandlerEventContextIncludesRequestMetadata(t *testing.T) {
	t.Parallel()

	bus := &recordingEventBus{}
	db := newRegistrationTestDB(nil)
	req, err := http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(`{"username":"kenton","password":"password123","email":"email@example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "203.0.113.5:12345"
	req.Header.Set("X-Request-Id", "req-123")
	req.Header.Set("Traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00")
	req.Header.Set("X-App-Distinct-ID", "distinct-123")
	req.Header.Set("Authorization", "Bearer secret")

	r := newTestRegistrar(db,
		WithDefaultRoleName("admin"),
		WithEventBus(bus),
		WithRequestEventMetadata(func(r *http.Request) events.ContextMetadata {
			return events.ContextMetadata{Values: map[string]string{"app_distinct_id": r.Header.Get("X-App-Distinct-ID")}}
		}),
	)

	rr := httptest.NewRecorder()
	r.ServeRegisterAPI(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	published := bus.singlePublish(t, events.EventRegistrationSucceeded)
	if published.event.RequestID != "req-123" {
		t.Fatalf("RequestID = %q, want %q", published.event.RequestID, "req-123")
	}
	if published.event.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("TraceID = %q, want parsed traceparent trace id", published.event.TraceID)
	}
	if published.event.Metadata["app_distinct_id"] != "distinct-123" {
		t.Fatalf("app_distinct_id = %q, want %q", published.event.Metadata["app_distinct_id"], "distinct-123")
	}
	if _, ok := published.event.Metadata["authorization"]; ok {
		t.Fatal("event metadata must not include authorization header")
	}
	if published.ctx.Value("ip_address") == "" {
		t.Fatal("published context should include ip_address")
	}
}

func TestRegistrationHandlerHintedRegistrationViaAuthorizationHeader(t *testing.T) {
	t.Parallel()

	hints := newTestHintServiceForHTTP(t)
	hint := registrationhints.Hint{
		ProviderName: "google",
		Email:        "rostered@example.com",
		Issuer:       "https://issuer.example.com",
		Subject:      "sub-123",
		Rosterable:   true,
	}
	token, err := hints.Create(hint)
	if err != nil {
		t.Fatalf("create hint token: %v", err)
	}

	db := newRegistrationTestDB(nil)
	req, err := http.NewRequest(http.MethodGet, "/api/register?hinted", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "RegistrationHint "+token)

	r := newTestRegistrarWithHint(db, hints, WithDefaultRoleName("admin"), WithDisablePublicRegistration(true))

	rr := httptest.NewRecorder()
	r.ServeRegisterAPI(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if _, found := db.FindOne("users", map[string]interface{}{"username": hint.Email}); !found {
		t.Fatal("expected hinted user to be created")
	}
	if len(db.Tables["auth_links"]) != 1 {
		t.Fatalf("auth_links count = %d, want 1", len(db.Tables["auth_links"]))
	}
}

func TestRegistrationHandlerHintedRegistrationViaCookie(t *testing.T) {
	t.Parallel()

	hints := newTestHintServiceForHTTP(t)
	hint := registrationhints.Hint{
		ProviderName: "github",
		Email:        "cookie-user@example.com",
		Issuer:       "https://issuer.example.com",
		Subject:      "sub-456",
		Rosterable:   true,
	}
	token, err := hints.Create(hint)
	if err != nil {
		t.Fatalf("create hint token: %v", err)
	}

	db := newRegistrationTestDB(nil)
	req, err := http.NewRequest(http.MethodGet, "/api/register?hinted", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: registrationhints.CookieName, Value: token})

	r := newTestRegistrarWithHint(db, hints, WithDefaultRoleName("admin"))

	rr := httptest.NewRecorder()
	r.ServeRegisterAPI(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if _, found := db.FindOne("users", map[string]interface{}{"username": hint.Email}); !found {
		t.Fatal("expected hinted user to be created")
	}
}

func TestRegistrationHandlerHintReplayCannotRecreateExistingUser(t *testing.T) {
	t.Parallel()

	hints := newTestHintServiceForHTTP(t)
	hint := registrationhints.Hint{
		ProviderName: "google",
		Email:        "replay@example.com",
		Issuer:       "https://issuer.example.com",
		Subject:      "sub-replay",
		Rosterable:   true,
	}
	token, err := hints.Create(hint)
	if err != nil {
		t.Fatalf("create hint token: %v", err)
	}

	db := newRegistrationTestDB(nil)
	r := newTestRegistrarWithHint(db, hints, WithDefaultRoleName("admin"))

	makeRequest := func() *httptest.ResponseRecorder {
		req, err := http.NewRequest(http.MethodGet, "/api/register?hinted", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "RegistrationHint "+token)
		rr := httptest.NewRecorder()
		r.ServeRegisterAPI(rr, req)
		return rr
	}

	first := makeRequest()
	if first.Code != http.StatusOK {
		t.Fatalf("first request status = %d, want %d: %s", first.Code, http.StatusOK, first.Body.String())
	}
	if len(db.Tables["users"]) != 1 {
		t.Fatalf("users count after first request = %d, want 1", len(db.Tables["users"]))
	}

	replay := makeRequest()
	if replay.Code != http.StatusConflict {
		t.Fatalf("replay status = %d, want %d: %s", replay.Code, http.StatusConflict, replay.Body.String())
	}
	if len(db.Tables["users"]) != 1 {
		t.Fatalf("users count after replay = %d, want 1 (no duplicate user created)", len(db.Tables["users"]))
	}
	if len(db.Tables["auth_links"]) != 1 {
		t.Fatalf("auth_links count after replay = %d, want 1 (no duplicate auth link created)", len(db.Tables["auth_links"]))
	}
}

func TestRegistrationHandlerRejectsInvalidHintToken(t *testing.T) {
	t.Parallel()

	hints := newTestHintServiceForHTTP(t)
	db := newRegistrationTestDB(nil)
	req, err := http.NewRequest(http.MethodGet, "/api/register?hinted", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "RegistrationHint not-a-real-token")

	r := newTestRegistrarWithHint(db, hints, WithDefaultRoleName("admin"))

	rr := httptest.NewRecorder()
	r.ServeRegisterAPI(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	if len(db.Tables["users"]) != 0 {
		t.Fatal("no user should be created for an invalid hint")
	}
}

func TestRegistrationHandlerAutoLoginIssuedWhenNoVerificationRequired(t *testing.T) {
	t.Parallel()

	db := newRegistrationTestDB(nil)
	tm := &recordingTokenManager{}
	r := newTestRegistrar(db, WithDefaultRoleName("admin"), WithTokenManager(tm))
	rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"email@example.com"}`)

	if tm.createCalls != 1 {
		t.Fatalf("CreateAuthToken calls = %d, want 1", tm.createCalls)
	}
	if tm.passCalls != 1 {
		t.Fatalf("PassToClient calls = %d, want 1", tm.passCalls)
	}
	if tm.passedToken == nil || tm.passedToken.AuthToken != "session-token" {
		t.Fatal("expected PassToClient to receive the created token")
	}

	var sawSessionCookie bool
	var sawExpiresCookie bool
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "token" {
			sawSessionCookie = true
		}
		if cookie.Name == "ls_expires_at" {
			sawExpiresCookie = true
		}
	}
	if !sawSessionCookie {
		t.Fatal("expected token cookie to be set")
	}
	if !sawExpiresCookie {
		t.Fatal("expected ls_expires_at cookie to be set")
	}
}

func TestRegistrationHandlerNoAutoLoginWhenEmailVerificationRequired(t *testing.T) {
	t.Parallel()

	sender := &recordingVerificationSender{}
	db := newRegistrationTestDB(nil)
	tm := &recordingTokenManager{}
	r := newTestRegistrar(db,
		WithDefaultRoleName("admin"),
		WithTokenManager(tm),
		WithRequiresEmailVerification(func(context.Context, database.DatabaseAccessor, users.LocksmithUserInterface, textvalidation.ValidationResultEvaluator) bool {
			return true
		}),
		WithAccountVerifier(verificationcodes.NewVerifier(db, sender)),
	)

	rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"email@example.com"}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if tm.createCalls != 0 {
		t.Fatalf("CreateAuthToken calls = %d, want 0", tm.createCalls)
	}
	if tm.passCalls != 0 {
		t.Fatalf("PassToClient calls = %d, want 0", tm.passCalls)
	}
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "token" {
			t.Fatal("token cookie should not be set when email verification is required")
		}
	}
}

func TestRegistrationHandlerBackwardCompatibleWithoutTokenManager(t *testing.T) {
	t.Parallel()

	db := newRegistrationTestDB(nil)
	r := newTestRegistrar(db, WithDefaultRoleName("admin"))
	rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"email@example.com"}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "token" {
			t.Fatal("token cookie should not be set without a configured TokenManager")
		}
	}
}

func TestRegistrationHandlerStrictJSONRejectsUnsupportedFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "rejects role escalation field",
			payload: `{"username":"kenton","password":"password123","email":"email@example.com","role":"admin"}`,
		},
		{
			name:    "rejects removed pwn ok field",
			payload: `{"username":"kenton","password":"password123","email":"email@example.com","pwnOK":true}`,
		},
		{
			name:    "rejects oauth restriction injection",
			payload: `{"username":"kenton","password":"password123","email":"email@example.com","restrictToOauthSrc":"github"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newRegistrationTestDB(nil)
			r := newTestRegistrar(db, WithDefaultRoleName("admin"))
			rr := performRegistrationRequest(t, r, tt.payload)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
			}
			if got := decodeRegistrationResponse(t, rr).Error; got != "could not unmarshal" {
				t.Fatalf("error = %q, want %q", got, "could not unmarshal")
			}
			if _, found := db.FindOne("users", map[string]interface{}{"username": "kenton"}); found {
				t.Fatal("unsupported field request inserted a user")
			}
		})
	}
}

func TestRegistrationHandlerRejectsDuplicateUsernameOrEmailCaseInsensitively(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		seed    map[string]interface{}
		payload string
	}{
		{
			name: "duplicate username with different case",
			seed: map[string]interface{}{
				"existing": map[string]interface{}{"username": "kenton", "email": "other@example.com"},
			},
			payload: `{"username":"KENTON","password":"password123","email":"new@example.com"}`,
		},
		{
			name: "duplicate email with different case",
			seed: map[string]interface{}{
				"existing": map[string]interface{}{"username": "someone", "email": "email@example.com"},
			},
			payload: `{"username":"newuser","password":"password123","email":"EMAIL@example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newTestRegistrar(newRegistrationTestDB(tt.seed), WithDefaultRoleName("admin"))
			rr := performRegistrationRequest(t, r, tt.payload)

			if rr.Code != http.StatusConflict {
				t.Fatalf("status = %d, want %d", rr.Code, http.StatusConflict)
			}
			if got := decodeRegistrationResponse(t, rr).Error; got != "taken" {
				t.Fatalf("error = %q, want %q", got, "taken")
			}
		})
	}
}

func TestRegistrationHandlerRejectsInviteEmailMismatchWithoutCreatingUser(t *testing.T) {
	t.Parallel()

	inviteCode := "jyTeL3RiH-9RgjLDt42CfTKJOVu9G16KebdGfVRygiu2Qf2Qkcb2QRRCQQDJVb210J2ZCz8v2PVJaDL56wuYPOHqiubfOk8M"
	hasher := sha256.New()
	hasher.Write([]byte(inviteCode))
	hashedCode := fmt.Sprintf("%x", hasher.Sum(nil))
	db := database.TestDatabase{Tables: map[string]map[string]interface{}{
		"invites": {
			"invite": map[string]interface{}{
				"code":    hashedCode,
				"email":   "invited@example.com",
				"role":    "admin",
				"inviter": "a-uuid",
				"sentAt":  time.Now().Unix(),
				"userid":  "invited-user-id",
			},
		},
		"users": {},
	}}

	r := newTestRegistrar(db, WithDefaultRoleName("admin"), WithDisablePublicRegistration(true))
	rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"other@example.com","code":"`+inviteCode+`"}`)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if got := decodeRegistrationResponse(t, rr).Error; got != "invalid email" {
		t.Fatalf("error = %q, want %q", got, "invalid email")
	}
	if _, found := db.FindOne("users", map[string]interface{}{"username": "kenton"}); found {
		t.Fatal("invite email mismatch inserted a user")
	}
	if _, found := db.FindOne("invites", map[string]interface{}{"code": hashedCode}); !found {
		t.Fatal("invite email mismatch expired the invite")
	}
}

func TestRegistrationHandlerEmailVerificationBehavior(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		requiresVerify    bool
		wantVerifierCalls int
	}{
		{name: "sends verification when required", requiresVerify: true, wantVerifierCalls: 1},
		{name: "does not send verification when not required", requiresVerify: false, wantVerifierCalls: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sender := &recordingVerificationSender{}
			db := newRegistrationTestDB(nil)
			r := newTestRegistrar(db,
				WithDefaultRoleName("admin"),
				WithRequiresEmailVerification(func(context.Context, database.DatabaseAccessor, users.LocksmithUserInterface, textvalidation.ValidationResultEvaluator) bool {
					return tt.requiresVerify
				}),
				WithAccountVerifier(verificationcodes.NewVerifier(db, sender)),
			)
			rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"email@example.com"}`)

			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
			}
			if sender.calls != tt.wantVerifierCalls {
				t.Fatalf("verification sends = %d, want %d", sender.calls, tt.wantVerifierCalls)
			}
			if tt.wantVerifierCalls == 1 && sender.forValue != "email@example.com" {
				t.Fatalf("verification destination = %q, want %q", sender.forValue, "email@example.com")
			}
		})
	}
}

func TestRegistrationHandlerEmailValidationResponses(t *testing.T) {
	t.Parallel()

	didYouMean := "user@example.com"
	tests := []struct {
		name        string
		validator   staticEmailValidator
		payload     string
		wantStatus  int
		wantConfirm bool
		wantBlocked bool
		wantDYM     string
	}{
		{
			name:        "asks client to confirm suspicious email",
			validator:   staticEmailValidator{result: textvalidation.ValidationResult_CONFIRM, dym: &didYouMean},
			payload:     `{"username":"kenton","password":"password123","email":"uesr@example.com"}`,
			wantStatus:  http.StatusBadRequest,
			wantConfirm: true,
			wantDYM:     didYouMean,
		},
		{
			name:        "blocks rejected email",
			validator:   staticEmailValidator{result: textvalidation.ValidationResult_REJECT},
			payload:     `{"username":"kenton","password":"password123","email":"blocked@example.com"}`,
			wantStatus:  http.StatusBadRequest,
			wantBlocked: true,
		},
		{
			name:       "allows confirmed suspicious email",
			validator:  staticEmailValidator{result: textvalidation.ValidationResult_CONFIRM, dym: &didYouMean},
			payload:    `{"username":"kenton","password":"password123","email":"uesr@example.com","validationok":true}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newRegistrationTestDB(nil)
			r := newTestRegistrar(db, WithDefaultRoleName("admin"), WithEmailValidation(tt.validator))
			rr := performRegistrationRequest(t, r, tt.payload)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
			res := decodeRegistrationResponse(t, rr)
			if res.ConfirmEmail != tt.wantConfirm {
				t.Fatalf("confirmEmail = %v, want %v", res.ConfirmEmail, tt.wantConfirm)
			}
			if res.EmailBlocked != tt.wantBlocked {
				t.Fatalf("emailBlocked = %v, want %v", res.EmailBlocked, tt.wantBlocked)
			}
			if res.DidYouMean != tt.wantDYM {
				t.Fatalf("didYouMean = %q, want %q", res.DidYouMean, tt.wantDYM)
			}
		})
	}
}

func TestRegistrationHandlerCustomUserHookCanPersistAdditionalFields(t *testing.T) {
	t.Parallel()

	db := newRegistrationTestDB(nil)
	r := newTestRegistrar(db, WithDefaultRoleName("admin"), WithConfigureCustomUser(func(lui users.LocksmithUser, db database.DatabaseAccessor) users.LocksmithUserInterface {
		return customUser{LocksmithUser: lui, CustomObject: "configured"}
	}))

	rr := performRegistrationRequest(t, r, `{"username":"kenton","password":"password123","email":"email@example.com"}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	inserted, found := db.FindOne("users", map[string]interface{}{"username": "kenton"})
	if !found {
		t.Fatal("registered user was not persisted")
	}
	user := inserted.(map[string]interface{})
	if user["customObject"] != "configured" {
		t.Fatalf("customObject = %v, want %q", user["customObject"], "configured")
	}
	if user["role"] != "admin" {
		t.Fatalf("role = %v, want %q", user["role"], "admin")
	}
}
