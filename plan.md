# Registration/Auth Refactor Plan

## Session 3 Summary (this pass)

This pass completed the core of goals 2, 3, and part of 5 from "Remaining Goals" below:

1. **Fixed the vendor-specific metadata regression.** `authentication/register/register_http.go` no longer hardcodes PostHog (or any vendor) header names. Default event metadata is limited to `ip_address` and `user_agent`. Applications that want additional metadata (analytics distinct IDs, etc.) must supply `RegistrationHandler.RequestEventMetadata func(*http.Request) events.ContextMetadata`, which is merged on top of the generic defaults. `routes.LocksmithRoutesOptions.RegistrationEventMetadata` exposes this at the route-composition layer.
2. **Implemented hinted/background registration end-to-end**, in `authentication/register/registrar.go` and `register_http.go`:
   - `RegistrationHandler.HintService registrationhints.Service` parses a hint from `Authorization: RegistrationHint <jwt>` or the `registration_hint` cookie, mirroring `internal/method_hint`.
   - A parsed hint is only trusted after `registrationhints.Service.Parse` succeeds (valid signature, `registration` audience, expiration, `Rosterable`, non-empty `ProviderName`/`Issuer`/`Subject`). Unsigned/invalid hints never reach the registrar.
   - `Registrar.Register` treats `req.Hint != nil` as authoritative: it skips password/length/email-as-username rules, uses `hint.Email` for username/email, sets `OAuthRestrictedSource` from the hint provider, and validates username/email as an email address.
   - Hinted registration is permitted even when `DisablePublicRegistration` is true (public **password** registration remains blocked without an invite).
   - `RegistrationRequest`/`RegistrationResult` (in `register.go`) now carry `Method`, `Provider`, `Background`, and `CreatedAuthLink`, and registration events (`requested`/`failed`/`succeeded`) use these instead of hardcoded `"password"`/`"password"`/`false`.
3. **Auth-link creation for hinted registration.** On successful hinted registration, `Registrar.Register` inserts a row into `auth_links` using only the signed hint's `ProviderName`/`Issuer`/`Subject` (never client JSON) plus the newly created user ID. This mirrors the shape used by `authenticator_domain.LinkedIdentity.ToMap()`.
4. **Shared session cookie helper.** Added `tokens.SetBaseCookies(w, token)` in `authentication/tokens/cookies.go`. `authenticator.authorizers.setBaseCookies` is now a thin wrapper around it, eliminating duplicate cookie logic between login and registration.
5. **Auto-login after registration.** `RegistrationHandler.TokenManager tokens.TokenManager` is optional; when set and the newly registered user does **not** require email verification, `ServeHTTP` calls `CreateAuthToken`, sets shared base cookies via `tokens.SetBaseCookies`, and calls `PassToClient`. If email verification is required, or no token manager is configured, registration still returns `200 OK` with no session issued (backward compatible).
6. **Route wiring.** `routes.LocksmithRoutesOptions` gained `RegistrationEventBus`, `RegistrationHintService`, `RegistrationTokenManager`, and `RegistrationEventMetadata`, all passed through to `register.RegistrationHandler`. `main.go` was intentionally left on legacy behavior (no event bus/hint service/token manager wired yet) since wiring those is an application-level decision outside this library change; see "Suggested Remaining Implementation Order" below.
7. **Tests added** (all passing, including under `-race`): hinted registration user/auth-link creation (`registrar_test.go`), disabled-public-registration bypass only for valid hints, rejection of incomplete/non-rosterable hints, duplicate-email handling, missing-email fail-closed behavior, hinted event method/provider/background correctness, HTTP-level hint registration via header and cookie, invalid-hint-token rejection, auto-login issuance and cookie assertions, no-auto-login-when-verification-required, backward compatibility without a token manager, and `tokens.SetBaseCookies` cookie-shape tests.

Validation run this pass:

```sh
go build ./...
go vet ./...
go test ./...
go test -race ./authentication/register/... ./authentication/events/... ./authentication/authenticator/... ./authentication/tokens/... ./authentication/registrationhints/...
```

All passed. The two `go vet` warnings that remain (`administration/invitations/invitations_test.go`, `database/mongo.go`) are pre-existing and untouched by this work.

### Still open after this pass

- The event-driven adapters/deprecation path for the legacy `NewRegistrationEvent`/`RequestNewRegistrationEvent` callbacks was not built; those fields still exist on `RegistrationHandler` but are simply unused by `ServeHTTP` now (they are not invoked at all, so any consumer still relying on them will silently stop receiving callbacks — this is a breaking behavior change that should be called out to consumers).
- `authentication/packets` and OAuth `registration_pkt` producers still exist in the codebase; registration no longer consumes them, but they have not been removed. **`authentication/oauth` (and anything under it, including `registration_pkt` producers there) is explicitly out of scope and deprecated — do not touch it.**
- Security-audit findings that touch `authentication/oauth` (OAuth login callback payload sanitization) are **out of scope per explicit instruction** and will not be fixed as part of this refactor.
- Remaining non-OAuth security-audit findings (auth/register/login rate limiting, invite-endpoint permission checks and response hardening, OIDC state/CSRF hardening, rate limiter TTL eviction, WebAuthn `FinishRegistration` response semantics) are **not yet fixed** — see "Security Audit Findings (Session 3)" below.

---

## Session 4 Summary

This pass wired the event bus into the remaining authenticator flows and closed out the registration-hint replay concern:

1. **Login success/failure events.** `authentication/authenticator/methods.go` gained an `eventBus events.Bus` field and `WithEventBus(bus)` option (defaults to `events.NoopBus{}`, mirroring the registrar's pattern). `attempt_login.go` adds a `publishAuthEvent` helper and now publishes:
   - `events.EventLoginSucceeded` (with `events.LoginSucceededPayload{UserID, Method, Provider, Passwordless}`) after a token is successfully created and passed to the client.
   - `events.EventLoginFailed` (new `events.LoginFailedPayload{PresentedUsername, Method, Reason}`, added to `authentication/events/auth_events.go`) for unknown-user, invalid-password, and passwordless-required failure branches. Reasons are stable enum-like strings (`user_not_found`, `invalid_password`, `passwordless_required`); no plaintext password is ever included.
   - `events.EventRosterStarted` (new `events.RosterStartedPayload{Provider}`) when `beginRegistrationRostering` issues a signed registration hint.
2. **Account-link events.** `authentication/authenticator/link_account.go`'s `LinkAccount` now publishes `events.EventAccountLinked` (`events.AccountLinkedPayload{UserID, Provider, Issuer, Subject}`) after a successful `auth_links` insert; no event is published if the insert fails.
3. **Hint replay / single-use protection, scoped as requested.** Rather than building token-level single-use tracking, verified that the existing pre-insert duplicate username/email check in `Registrar.Register` already prevents an already-registered user from being recreated by replaying a signed registration hint (the second attempt fails with `ErrRegistrationTaken` before any DB write, and no duplicate `auth_links` row is created). Added explicit regression tests at both the registrar level (`TestRegistrarHintedRegistrationReplayCannotRecreateExistingUser`) and the HTTP level (`TestRegistrationHandlerHintReplayCannotRecreateExistingUser`) proving: first use creates exactly one user and one auth link; replay of the identical signed token returns `409 Conflict`/`ErrRegistrationTaken`; user/auth-link counts remain at 1 after replay; the original user is unaffected.
4. **Tests added** in `authentication/authenticator/events_test.go` (new file) and `authentication/register/registrar_test.go` / `register_behavior_test.go`, all passing under `go test` and `go test -race`.

Validation run this pass:

```sh
go build ./...
go test ./...
go test -race ./authentication/register/... ./authentication/authenticator/... ./authentication/events/... ./authentication/tokens/...
```

All passed.

### Explicit scope boundary (per instruction)

`authentication/oauth` (including everything under it) is deprecated and **out of scope** for this refactor. No files under that package were read, modified, or otherwise touched in Session 4, and the OAuth-login-callback-payload security finding from Session 3 will **not** be remediated as part of this work.

### Still open after Session 4

- True single-use/consumed-hint tracking (e.g. a `jti` denylist) was intentionally not built; the duplicate-user check is the agreed-upon replay defense per explicit instruction, not a token-consumption mechanism. If two *different*, not-yet-registered emails share a hint's provider/issuer/subject (which shouldn't happen for a correctly-issued hint), that scenario is not separately guarded — this is considered out of scope unless a concrete threat model requires it.
- `EventRosterSucceeded`/`EventRosterFailed` are still not published anywhere (only `EventRosterStarted` was wired, since success/failure of rostering is really just registration succeeding/failing afterward, which is already covered by `EventRegistrationSucceeded`/`EventRegistrationFailed`).
- Non-OAuth security-audit backlog items (rate limiting, invite endpoint hardening, OIDC CSRF/state, rate-limiter TTL eviction, WebAuthn error handling) remain open.

---

## Remaining Goals

1. ~~Refactor `authentication/register` into the same handler/session/method format as the new `authentication/authenticator` package.~~ Scaffolding + core orchestration done (see Phase 4/5 below); full `register_domain.Session`-based dispatch through `Registrar` is still not wired (the HTTP handler still builds the core `RegistrationRequest` itself rather than delegating to `register_domain.Handler`/`Session` implementations).
2. ~~Use the completed neutral `authentication/registrationhints` package to support hinted/background registration for rosterable/federated users.~~ Done this pass (see Session 3 Summary).
3. ~~Automatically log users in after successful registration when policy allows it.~~ Done this pass (see Session 3 Summary).
4. Finish replacing the old `packets` capability with registration hints and, if needed, a purpose-built future registration grant system. `authentication/packets` still exists and is unused by registration; removal/replacement elsewhere is still open. (Note: `authentication/oauth` and any `registration_pkt` producers under it are out of scope/deprecated — do not touch.)
5. ~~Wire the completed `authentication/events` bus into registration/login/roster/account-link flows and deprecate callback-style hooks.~~ Registration, login success/failure, roster-started, and account-link events are wired (see Session 3/4 Summaries). Legacy callback deprecation adapters were not built (the old registration callback fields are simply dead now).

---

## Current State

### `authentication/authenticator`

The new authenticator package has a good structure:

```text
authentication/authenticator/
  attempt_login.go
  methods.go
  start_provider.go
  link_account.go

  authenticator_domain/
    handler.go
    errors.go
    token.go
    linked_identity.go

  authenticator_methods/
    oidc.go
    password.go

  internal/
    method_oidc/
    method_password/
    sessions/
```

Important concepts already in place:

- `authenticator_domain.Handler`
- `authenticator_domain.Session`
- `authenticator_domain.Beginnable`
- `authenticator_domain.FederatedIdentity`
- `authenticator_domain.VerifiedContact`
- `authenticator_domain.Rosterable` returning `*registrationhints.Hint`
- `tokens.TokenManager`
- `signing.SigningPackageInterface`

### `authentication/register`

Current register package is mostly one large HTTP handler:

```text
authentication/register/
  register_http.go
  register_http_test.go
  README.md
```

`register_http.go` currently handles:

- HTTP method validation
- request parsing
- public registration checks
- email-as-username normalization
- strict request parsing with unknown JSON fields rejected
- field validation
- email validation
- password length checks
- username/email format checks
- invite validation
- duplicate user detection
- password compilation
- user construction
- custom user configuration
- email verification decisions
- DB insertion
- invite expiration
- verification email sending
- registration callbacks

This should still be split into smaller domain, method, and orchestration components.

Already removed from the current registration path:

- HIBP registration password checks and `PwnOK`
- old `Packet` authorization parsing
- old auto-login-token generation through `RegistrationJWTService`
- `registration_pkt` page-cookie consumption

---

## Handoff Context for Future Agents

Completed in the first implementation pass:

1. Added `authentication/registrationhints` with `Hint`, `Service`, secure cookie helpers, and tests.
   - `Service.Create`/`Parse` require valid signature, `registration` audience, expiration, `Rosterable`, `ProviderName`, `Issuer`, and `Subject`.
   - Cookie helpers use `HttpOnly`, `Secure`, `SameSite=Lax`, `/register`, and explicit expiry/max-age.
2. Updated `authentication/authenticator` roster flow to use `registrationhints.Hint` instead of `authenticator_domain.RegistrationHint`.
   - OIDC roster hints now populate `Issuer` and `Subject`.
3. Removed HIBP from the registration flow.
   - `authentication/hibp` remains because it is still used outside registration.
4. Removed `authentication/packets` from the registration path and route wiring.
   - `authentication/packets` and some OAuth code still exist, but registration no longer consumes `registration_pkt` or `Packet` authorization.
5. Added `authentication/events` with event names, envelope, context metadata helpers, `Bus`, `NoopBus`, synchronous `MemoryBus`, and tests.
   - Registration lifecycle events are now partially wired through `Registrar`.
   - Login/roster/account-link flows are not wired yet.
6. Hardened current `/api/register` request parsing with `json.Decoder.DisallowUnknownFields()` and a 1 MiB body limit.
7. Added Phase 4 register domain/method scaffolding:
   - `authentication/register/register_domain`
   - `authentication/register/register_methods`
   - `authentication/register/internal/method_password`
   - `authentication/register/internal/method_hint`
8. Added Phase 5 core `Registrar` extraction:
   - `authentication/register/register.go`
   - `authentication/register/registrar.go`
   - `authentication/register/responses.go`
   - `RegistrationHandler.ServeHTTP` now delegates successful registration work to `Registrar.Register`.
9. Added behavior tests for strict JSON/mass-assignment rejection, duplicates, invites, email validation/verification, custom user hooks, registrar success, method sessions, registration event publication, and event request metadata propagation.
10. Added `events.Bus` and `*slog.Logger` options to `Registrar`/`RegistrationHandler`.
    - `Registrar.Register` publishes registration requested/succeeded/failed and email-verification-sent events using curated payloads.
    - Direct `logger.LOGGER.Log(logger.REGISTRATION_SUCCESS)` and `logger.LOGGER.Log(logger.INVITE_CODE_USED)` calls were removed from the registration success path.
    - Deprecated registration callback fields remain on `RegistrationHandler`, but are no longer invoked directly by the handler; future compatibility adapters should subscribe to `events.Bus` if needed.

Validation from the latest pass:

```sh
go test ./authentication/register/...
go test ./authentication/registrationhints ./authentication/events ./authentication/authenticator/... ./authentication/register ./routes
go test ./...
```

Both commands passed.

Important remaining security constraints:

- Do not auto-login users requiring email verification unless the registration method supplied a server-verified email identity.
- Keep registration request DTOs explicit; never map request bodies directly into `users.LocksmithUser` or role fields.
- Registration hints are short-lived but not single-use yet; add replay protection if threat model requires it.
- Event payloads currently include some PII-capable fields by design. Do not log whole envelopes by default, and consider redaction/options before external forwarding.
- Event request metadata must remain generic and app-controlled. Do not hardcode vendor-specific analytics concepts (for example PostHog header names) into `authentication/register` or `authentication/events`; applications should enrich context before invoking handlers, or this package should expose generic middleware/hooks/options for allowlisted metadata extraction.
- Public login/register/provider-start rate limiting is still not implemented in this refactor.
- `Registrar.Register` supports password/invite registration only. Hinted registration method/session scaffolding exists, but core hinted registration, auth-link creation, and auto-login are not implemented yet.

---

## Target Package Layout

```text
authentication/register/
  register.go                    # Registrar constructor/options
  serve_register_api.go           # HTTP API orchestration
  serve_register_page.go          # Registration page handler, initially migrated from register_http.go
  registrar.go                    # Core registration orchestration
  responses.go                    # API response structs

  register_domain/
    handler.go                    # Handler/session interfaces
    request.go                    # Normalized registration request
    result.go                     # Registration result
    errors.go                     # Typed registration errors
    validators.go                 # Optional domain validation contracts

  register_methods/
    password.go                   # Public password registration method constructor
    hint.go                       # Public hinted/background registration method constructor

  internal/
    method_password/
      method_password.go
      password_session.go
      password_validation.go

    method_hint/
      method_hint.go
      hint_session.go
      hint_loader.go
```

New shared packages:

```text
authentication/registrationhints/
  hint.go
  service.go
  errors.go
  README.md

authentication/events/
  bus.go
  event.go
  noop.go
  memory.go
  auth_events.go
  registration_events.go
```

Potential shared token/session helper:

```text
authentication/sessioncookies/
  cookies.go
```

---

## Phase 1 — Move Registration Hints to a Neutral Package ✅ Complete

Registration hints have been moved to `authentication/registrationhints`.

Implemented files:

```text
authentication/registrationhints/
  hint.go
  service.go
  errors.go
  service_test.go
```

Current behavior:

- `registrationhints.Hint` embeds `jwt.RegisteredClaims` and carries provider, email, display name, issuer, subject, and rosterability.
- `registrationhints.Service` creates/parses signed JWTs with a default 60 second TTL.
- Parsing validates signature, `registration` audience, required expiration, expiry time, `Rosterable`, `ProviderName`, `Issuer`, and `Subject`.
- `FromRequest` loads from `X-Registration-Hint` or the `registration_hint` cookie.
- `SetCookie`, `SetCookieWithTTL`, and `ClearCookie` set secure cookie attributes.

Authenticator updates are also complete:

- `authenticator_domain.Rosterable` now returns `*registrationhints.Hint`.
- `beginRegistrationRostering` delegates JWT creation/cookies to `registrationhints`.
- OIDC rosterable sessions populate `Issuer` and `Subject`.

---

## Phase 2 — Remove HIBP from Registration for Now ✅ Complete

HIBP has been removed from the registration flow.

Completed removals:

- `RegistrationHandler.HIBP`
- registration `hibp` imports
- `PwnOK` request handling
- HIBP goroutine/channel logic
- pwned-password registration responses/tests
- route wiring of `HIBPIntegrationOptions` into registration handlers

`authentication/hibp` was intentionally kept because it is still used outside registration.

---

## Phase 3 — Remove the Current Packets System 🟡 Partially Complete

The current registration path no longer consumes the old `packets` system.

Completed removals from registration/route wiring:

- `authentication/register` no longer imports `authentication/packets`.
- `RegistrationHandler.RegistrationJWTService` was removed.
- `validateTrustedJWT` and `Authorization: Packet ...` handling were removed.
- `RegistrationJWTService.CreateAutoLoginToken(...)` was removed from registration.
- `routes.LocksmithRoutesOptions.RegistrationJWTService` was removed.
- `RegistrationPageHandler` no longer consumes the `registration_pkt` cookie.

Still present outside the registration path:

- `authentication/packets`
- OAuth code that can still create/set `registration_pkt`

Future agents should either delete those remaining packet producers if no longer needed, or replace them with the new `registrationhints`/future `registrationgrants` design. Registration itself no longer depends on them.

### Replacement design

Replace Packets with two separate concepts:

1. **Registration hints**
   - generated by authenticator when a verified identity is rosterable but no local user exists
   - short-lived
   - browser-oriented
   - carried in `registration_hint` cookie or explicit header
   - used to create a local user and auth link

2. **Registration grants**
   - future server-side/API equivalent for trusted app-to-app registration
   - not implemented in the first pass unless needed immediately
   - explicitly scoped and auditable
   - not called `Packet`

Potential future package:

```text
authentication/registrationgrants/
  grant.go
  service.go
  errors.go
```

Potential grant claims:

```go
type Grant struct {
	jwt.RegisteredClaims

	Subject string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name,omitempty"`

	ProviderName string `json:"provider_name,omitempty"`
	Issuer       string `json:"issuer,omitempty"`
	ProviderSub  string `json:"provider_sub,omitempty"`

	Scopes []string `json:"scopes"`
}
```

Required scopes could include:

- `registration:create`
- `registration:login_after_create`
- `auth_link:create`

Do not keep compatibility with the old `Packet` auth format unless a consuming app still requires it. If compatibility is needed, implement a temporary adapter that converts old packets into new grants and mark it deprecated.

---

## Phase 4 — Refactor Register into Handler/Session Methods ✅ Scaffold Complete

The register domain/method scaffolding now exists and compiles. The current HTTP API is not yet dispatched through these method handlers, but the interfaces and password/hint sessions are available for the next integration step.

Implemented packages:

```text
authentication/register/register_domain/
authentication/register/register_methods/
authentication/register/internal/method_password/
authentication/register/internal/method_hint/
```

Current behavior:

- Password method handles `POST` JSON registration input using explicit DTO fields and `DisallowUnknownFields`.
- Hint method handles `registration_hint` cookie and `Authorization: RegistrationHint <jwt>`.
- Hint method uses `registrationhints.Service` and creates a normalized passwordless domain request with `Hint` populated.

Still remaining:

- Wire method dispatch into a future API orchestration layer, or keep current `RegistrationHandler` as compatibility wrapper.
- Convert `Registrar.Register` from the legacy `RegistrationRequest` type to `register_domain.Request` when ready.
- Implement actual hinted registration behavior in core registrar logic.

### Domain interfaces

Implemented:

```text
authentication/register/register_domain/handler.go
```

Suggested interfaces:

```go
type Handler interface {
	CanHandle(r *http.Request) bool
	Name() string
	Session(db database.DatabaseAccessor) Session
}

type Session interface {
	LoadRequest(r *http.Request) error
	RegistrationRequest() Request
}
```

### Normalized request

Create:

```text
authentication/register/register_domain/request.go
```

Suggested request:

```go
type Request struct {
	Username string
	Password string
	Email    string
	Code     string

	ValidationOK bool

	Hint *registrationhints.Hint
}
```

No `PwnOK` once HIBP is removed.

### Result

Create:

```text
authentication/register/register_domain/result.go
```

Suggested result:

```go
type Result struct {
	User users.LocksmithUserInterface

	Method   string
	Provider string

	CreatedAuthLink bool
	InviteUsed      bool
	RequiresEmailVerification bool
}
```

### Registrar constructor

Create:

```text
authentication/register/register.go
```

Suggested shape:

```go
type Registrar struct {
	db  database.DatabaseAccessor
	log *slog.Logger

	methods []register_domain.Handler
	tm      tokens.TokenManager
	hints   registrationhints.Service
	events  events.Bus

	defaultRoleName           string
	disablePublicRegistration bool
	configureCustomUser       RegisterCustomUserFunc
	requiresEmailVerification func(...)
	accountVerifier          verificationcodes.Verifier
	emailValidation           textvalidation.EmailValidator
	emailAsUsername           bool
	minimumLengthRequirement int
}

type Option func(*Registrar)

func NewRegistrar(db database.DatabaseAccessor, opts ...Option) *Registrar
```

### Methods

Implement at least two registration methods:

1. Password registration

```text
authentication/register/register_methods/password.go
authentication/register/internal/method_password/
```

Handles normal JSON body registration.

2. Hint registration

```text
authentication/register/register_methods/hint.go
authentication/register/internal/method_hint/
```

Handles requests with:

- `registration_hint` cookie
- or `Authorization: RegistrationHint <jwt>`

This method should support background registration without requiring a password.

---

## Phase 5 — Registration Orchestration 🟡 Password/Invite Core Complete, Hint/Event Work Remaining

The actual password/invite registration operation has been extracted from HTTP into `Registrar.Register`.

Implemented files:

```text
authentication/register/register.go
authentication/register/registrar.go
authentication/register/responses.go
```

Completed behavior in `Registrar.Register`:

- role config validation
- disabled public registration + invite requirement for current password path
- email-as-username normalization
- required field validation
- email validation behavior
- password length validation
- username/email regex validation
- invite validation
- duplicate username/email detection
- password compilation
- user construction
- custom user hook
- email verification decision
- DB insertion
- invite expiration
- verification email send

`RegistrationHandler.ServeHTTP` now owns only HTTP concerns:

- method/body validation
- strict JSON DTO decoding
- pre-DB compatibility checks for legacy tests
- calling `Registrar.Register`
- mapping typed registration errors to HTTP responses
- legacy callback invocation

Still remaining:

- Convert core registrar API to consume `register_domain.Session` / `register_domain.Request` directly.
- Add hinted/background registration to `Registrar.Register`.
- Create `auth_links` rows during hinted registration.
- Publish events through `authentication/events.Bus`.
- Remove or adapt legacy callbacks once event bus wiring exists.

### Auth link creation

Hinted/federated registration should create an `auth_links` row when the hint includes:

- provider name
- issuer
- subject
- user ID

Use the same shape as `authenticator_domain.LinkedIdentity` currently inserted by:

```text
authentication/authenticator/link_account.go
```

Consider moving `LinkedIdentity` to a neutral package too if both authenticator and register need to construct it.

Potential package:

```text
authentication/linkedidentities/
  linked_identity.go
```

---

## Phase 6 — Auto-Login After Registration

### Objective

Successful registration should automatically log the user in when allowed.

### Replace old auto-login token system

Do not carry forward:

```go
RegistrationJWTService.CreateAutoLoginToken(...)
```

Instead use the same `tokens.TokenManager` path as authenticator:

```go
token, err := r.tm.CreateAuthToken(user)
if err != nil {
	return err
}

if err := sessioncookies.SetBaseCookies(w, token); err != nil {
	return err
}

if err := r.tm.PassToClient(w, req, token); err != nil {
	return err
}
```

### Shared cookie helper

Move this logic out of `authentication/authenticator/attempt_login.go`:

```go
func (a *authorizers) setBaseCookies(w http.ResponseWriter, token *authenticator_domain.Token) error
```

Into:

```text
authentication/sessioncookies/cookies.go
```

This lets both login and registration set consistent cookies.

### Email verification policy

Important behavior decision:

- If a new user requires email verification, do not auto-login unless the registration method already supplied a verified email.
- For OIDC registration hints with verified email, either:
  - mark the user as not requiring email verification, or
  - let `RequiresEmailVerification` inspect the hint/method and return false.

This keeps registration consistent with login, where users requiring email verification are currently blocked.

---

## Phase 7 — Internal Event Bus 🟡 Registration Partially Wired, App Metadata Hook Remaining

The `authentication/events` package now exists with:

```text
authentication/events/
  bus.go
  event.go
  noop.go
  memory.go
  auth_events.go
  registration_events.go
  memory_test.go
```

Implemented pieces:

- event names for registration, login, roster, account-link, and email-verification flows
- `Envelope` with generic metadata support
- typed context metadata helpers in `events.ContextMetadata`, `events.WithContextMetadata`, `events.MetadataFromContext`, and `events.EnrichEnvelope`
- `Handler`
- `Subscription`
- `Bus`
- `NoopBus`
- synchronous in-memory `MemoryBus`
- payload structs for registration requested/succeeded/failed, email-verification sent, login success, and account linking
- `Registrar` options for `events.Bus` and `*slog.Logger`
- registration requested/succeeded/failed and email-verification-sent publication from `Registrar.Register`
- behavior tests for no-subscriber publish, delivery, multiple subscribers, unsubscribe, deterministic handler errors, registration success/failure events, and generic request metadata propagation

Still remaining:

- Correct the current request metadata extraction so package code stays generic/app-defined. If metadata is extracted in this package, expose a generic option/hook such as `RequestEventMetadata func(*http.Request) events.ContextMetadata` or middleware that apps can compose; do not hardcode product/vendor-specific header names.
- Add `events.Bus` options to authenticator constructors.
- Publish events from login success/failure, roster start/success/failure, and account linking.
- Adapt or deprecate old callback fields:

```go
NewRegistrationEvent
RequestNewRegistrationEvent
LoginInfoCallback
RequestLoginInfoCallback
```

Temporary compatibility path:

- keep callback fields for one release window if public API compatibility matters
- register adapters that subscribe to events and call the old callbacks
- mark fields with `// Deprecated: use authentication/events.Bus instead.`

Security/privacy guidance:

- Keep payloads stable, small, and safe.
- Avoid publishing password hashes, tokens, raw OAuth claims, full user maps, raw invite codes, cookies, authorization headers, or full `*http.Request` values.
- Treat `Email`, `Username`, `Issuer`, and `Subject` as sensitive if forwarded outside process boundaries.
- Request metadata must be allowlisted by the application. The library may provide generic plumbing, but app-specific analytics details belong in the consuming application.

---

## Phase 8 — Tests

### Registration hints ✅ Mostly Complete

Already covered:

- hint creation sets ID, audience, issued-at, expiration
- valid hint parses successfully
- expired hint fails
- wrong audience fails
- non-rosterable hint fails
- missing subject/issuer fails closed
- cookie helpers set secure attributes and expiry

Remaining when hinted registration is implemented:

- missing email behavior for the registration method, if email remains required
- provider mismatch/unknown provider behavior
- replay/single-use behavior if replay protection is added

### Register refactor

Add tests for:

- password registration succeeds
- password registration rejects duplicate username/email
- password registration rejects invalid email
- password registration rejects short password
- invite registration succeeds
- invite email mismatch fails
- public registration disabled blocks normal registration
- public registration disabled allows valid hinted/background registration, if desired
- custom user configuration is applied
- email verification is sent when required

### Background/hinted registration

Add tests for:

- valid registration hint creates user
- valid registration hint creates `auth_links` record
- valid registration hint auto-logs in user
- consumed/expired hint cannot be reused if replay protection is added

### Auto-login

Add tests for:

- successful registration calls `TokenManager.CreateAuthToken`
- successful registration calls `TokenManager.PassToClient`
- base cookies are set
- user requiring email verification does not auto-login unless method is trusted/verified

### Event bus 🟡 Registration Flow Partially Covered, Remaining Flow Tests Needed

Already covered:

- publish with no subscribers succeeds
- subscriber receives event
- multiple subscribers receive event
- unsubscribe works
- handler error behavior is deterministic
- registration success publishes `auth.registration.succeeded`
- registration failure publishes `auth.registration.failed`
- registration event context/envelope receives generic request metadata
- sensitive headers are not copied into registration event metadata by default

Remaining after event wiring:

- registration requested publishes `auth.registration.requested` in all intended parse/validation paths
- email verification send publishes `auth.email_verification.sent`
- app-provided request metadata enrichment works without vendor-specific hardcoding
- login success publishes `auth.login.succeeded`
- login failure publishes `auth.login.failed`
- roster flow publishes roster events
- account link publishes `auth.account_linked`

---

## Phase 9 — Migration/Removal Checklist

### Remove or replace imports ✅ Complete for registration path

Completed for the registration path:

```go
authentication/hibp
authentication/packets
```

Remaining cleanup outside registration:

- Decide whether to delete or replace `authentication/packets` and OAuth `registration_pkt` producers.
- Keep `authentication/hibp` unless/until reset/login/password-change flows no longer use it.

### Remove old fields from new registrar path ✅ Complete for current handler

Already removed from `RegistrationHandler`/route options:

```go
HIBP hibp.HIBPSettings
RegistrationJWTService packets.RegistrationJWTServiceInterface
```

Do not reintroduce these in the new registrar.

Temporarily keep deprecated callback fields only if existing public API compatibility matters:

```go
NewRegistrationEvent
RequestNewRegistrationEvent
```

### Move or neutralize shared auth domain types

Completed:

```go
authenticator_domain.RegistrationHint -> authentication/registrationhints.Hint
```

Consider moving later:

```go
authenticator_domain.Token -> authentication/tokens.Token or authentication/session.Token
authenticator_domain.LinkedIdentity -> authentication/linkedidentities.LinkedIdentity
```

---

## Suggested Remaining Implementation Order

Completed earlier: registration hints package, authenticator roster-hint migration, HIBP removal from registration, registration packet removal, and event bus scaffolding.

Recommended next steps:

1. ~~Convert `Registrar.Register` to use `register_domain.Session` / `register_domain.Request` directly, or add a thin adapter from method sessions to the current core request.~~ Deferred: `RegistrationHandler.ServeHTTP` still builds `RegistrationRequest` directly rather than going through `register_domain.Handler`/`Session`. Not required for functional completeness of hinted/password registration, but still open for the originally planned architecture.
2. ~~Add hinted/background registration to core registrar logic using `registrationhints.Service`.~~ Done.
3. ~~Create auth-link rows during hinted registration using provider + issuer + subject.~~ Done.
4. ~~Add shared session cookie helper.~~ Done (`tokens.SetBaseCookies`).
5. ~~Add auto-login after registration via `tokens.TokenManager`, gated by email-verification policy.~~ Done.
6. ~~Replace hardcoded request metadata extraction with generic app-owned enrichment before continuing event work.~~ Done (`RegistrationHandler.RequestEventMetadata`).
7. Wire `events.Bus` into login/roster/account-link flows. **Still open.**
8. Adapt or deprecate old callbacks. **Partially open** — `NewRegistrationEvent`/`RequestNewRegistrationEvent` are no longer invoked by `ServeHTTP` but remain on the struct undocumented as dead; needs an explicit deprecation notice/adapter or removal.
9. Decide whether to delete/replace remaining `authentication/packets` and OAuth `registration_pkt` producers. **Still open.**
10. ~~Add tests for hinted registration, auth-link creation, auto-login, and event publication.~~ Done.
11. Update `authentication/register/README.md`. **Still open.**
12. Add replay protection (single-use) for registration hints. **Still open**, see Open Decision 3.
13. Add rate limiting to `/api/register`, `/api/login`, and provider-start endpoints. **Still open**, flagged by security audit.
14. Sanitize OAuth login callback (`LoginInfoCallback`/`RequestLoginInfoCallback`) payloads to a safe DTO instead of `user.ToMap()`. **Still open**, flagged by security audit as high severity (potential password hash/session leakage into logs/callbacks).
15. Add explicit permission checks and response hardening (no-store, size limits, `DisallowUnknownFields`) to the invite-issuance admin endpoint. **Still open**, flagged by security audit.
16. Harden OIDC callback state handling with a signed/expiring nonce bound to a server-side cookie, beyond the current relative-path redirect check. **Still open**, flagged by security audit.
17. Add TTL-based eviction to `ratelimits.RateLimiter`'s per-identifier limiter maps to bound memory growth. **Still open**, flagged by security audit.
18. Harden `authentication/webauth/webauthn_http.go` `FinishRegistration` to return explicit status codes/structured errors instead of silently returning. **Still open**, flagged by security audit.

---

## Security Audit Findings (Session 3)

A read-only security review (api-audit / owasp-audit / golang-security oriented) was performed against the registration/authenticator event surface, invite handling, and OAuth callback flow. No fixes from this list were applied in this pass except item 1 (vendor-specific metadata), which is already reflected above. This is a prioritized backlog, not a completed remediation:

1. **[Fixed this pass] Event metadata hardcoding/vendor coupling** — `authentication/register/register_http.go` previously hardcoded PostHog header names. Now generic; app-specific metadata must be supplied via `RequestEventMetadata`.
2. **[High, open] OAuth login callback payload exposure** — `authentication/oauth/login_user.go` passes `user.ToMap()` (potentially including password hash/salt/sessions) to `LoginInfoCallback`/`RequestLoginInfoCallback`, and `main.go` logs the full user struct. Replace with a curated safe DTO (ID/username/email only) before shipping to callbacks or logs.
3. **[High, open] No rate limiting on `/api/register`, `/api/login`, or provider-start endpoints.** `ratelimits.RateLimiter` exists but is only wired into `endpoints.SecureEndpointHTTPMiddleware` for already-authenticated routes. Add per-IP and per-identifier (username/email/invite-code-hash) limits to the unauthenticated auth surface.
4. **[Medium, open] Invite issuance endpoint hardening.** `administration/invitations/invitations_http.go` returns the raw invite token in the response body without an explicit fine-grained permission check (only "authenticated" is verified), without `DisallowUnknownFields()`/body size limit, and without `Cache-Control: no-store`.
5. **[Medium, open] OIDC callback CSRF/state hardening.** `authentication/oauth/oidc/oidc.go` uses the OAuth `state` parameter primarily as a redirect-path carrier; it should also function as a signed, short-lived, single-use CSRF nonce bound to a server-side cookie.
6. **[Medium, open] Unbounded rate-limiter key growth.** `ratelimits/ratelimits.go` stores one `*rate.Limiter` per identifier in `sync.Map` with no eviction; an attacker who can vary the limiter key can grow memory unboundedly.
7. **[Low/Medium, open] `authentication/webauth/webauthn_http.go` `FinishRegistration`** silently returns without writing a status code on several error paths, harming both API consumers and auditability.
8. **[Informational] Deprecated registration callbacks are now dead code.** `RegistrationHandler.NewRegistrationEvent`/`RequestNewRegistrationEvent` are no longer called anywhere. This is intentional per the events-bus migration, but is a breaking behavior change for any existing caller relying on those fields; it should be called out in release notes or bridged via an event-subscriber adapter if backward compatibility is required.

Items 2–7 should be treated as separate, focused follow-up changes (each in its own reviewable diff per the security-audit skill's guidance), not bundled into the registration refactor itself.

---

## Open Decisions

1. Should hinted/background registration be allowed when public registration is disabled?
   - Recommended: yes, if the hint is valid, rosterable, and signed by the server.

2. Should a user requiring email verification be auto-logged in?
   - Recommended: no, unless the method provided verified email identity, such as OIDC.

3. Should registration hints be single-use?
   - Resolved (Session 4): per explicit instruction, replay protection is scoped to "an already-existing user cannot be created again," which the existing pre-insert duplicate username/email check already guarantees (verified with regression tests at the registrar and HTTP layers). True token-level single-use/consumed-hint tracking (e.g. a `jti` denylist) was intentionally not built.

4. Should old `Packet` clients get a compatibility adapter?
   - Recommended: no, unless a known consuming app requires it. Prefer a clean break to registration hints/grants.

5. Should `LinkedIdentity` move to a neutral package now?
   - Recommended: yes if hinted registration creates auth links in the first pass.

6. Should event handlers run synchronously or asynchronously?
   - Recommended: make the bus configurable. Default to non-blocking/noop for app hooks, but allow synchronous audit/security subscribers.

7. How should request metadata be supplied to events?
   - Required direction: app-owned/generic enrichment only. Do not hardcode analytics provider details in this package. Prefer middleware or a `RegistrationHandler` option that lets the application return `events.ContextMetadata` from `*http.Request`.
