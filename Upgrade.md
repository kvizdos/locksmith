# Locksmith Upgrade Guide: `alpha-ui2-v1.2.0` → current `main`

This guide is for consumers who wired up Locksmith using the `main.go`
boilerplate as it existed at tag `alpha-ui2-v1.2.0`, and now need to move to
the architecture on `main` (18 commits ahead, unreleased as of this writing —
`git describe` reports `alpha-ui2-v1.2.0-18-g2afd6f6`).

This is **not** a drop-in, backward-compatible update. The authentication
core was rewritten: login, registration, OAuth/OIDC, and CSRF protection all
moved from scattered options/packages into two explicit, composable
orchestrators (`authenticator` and `register`) built with functional options,
plus a new event bus for observability. Several packages were deleted
outright (`authentication/login`, `authentication/xsrf`,
`authentication/oauth/oidc`).

Treat this as a **security-relevant migration**, not a mechanical refactor.
Section 5 lists behavior changes that affect your threat model (brute-force
protection, CSRF, user enumeration). Read it before you ship.

---

## 1. TL;DR — what moved where

| Old (`alpha-ui2-v1.2.0`) | New (`main`) |
|---|---|
| `authentication/login` package, `login.LoginOptions{LockoutPolicy: ...}` | **Deleted.** No direct replacement — see [§5.1](#51-brute-force--account-lockout-protection-was-removed). |
| `authentication/xsrf`, `xsrf.XSRFSigningPackage.{Anonymous,Authenticated}` | **Deleted.** No global XSRF signer to configure — see [§5.3](#53-global-xsrf-signing-package-is-gone). |
| `authentication/oauth/oidc` (`oauth_google_oidc.NewOIDCConnection`) | `authenticator.AllowMethodOIDC(authenticator_methods.WithOIDC(...))`, registered via `authenticator.WithMethods(...)` |
| `routes.LocksmithRoutesOptions.OAuthProviders []oauth.OAuthProvider` | Removed. OAuth providers are now methods passed to `authenticator.NewAuthorizer(...)`. |
| `routes.LocksmithRoutesOptions.LoginSettings *login.LoginOptions` | Removed. See `authenticator.WithMinimumResponseTime(...)`. |
| `routes.LocksmithRoutesOptions.RequiresEmailVerification func(...)` | Moved to `register.WithRequiresEmailVerification(...)` on the registrar. |
| `routes.LocksmithRoutesOptions.NewRegistrationEvent func(users.LocksmithUserInterface)` | Replaced by subscribing to `events.EventRegistrationSucceeded` on an `events.Bus`. |
| `routes.LocksmithRoutesOptions.LoginInfoCallback func(method string, user map[string]any)` | Replaced by subscribing to `events.EventLoginSucceeded` / `EventLoginFailed`. (Field still exists on the options struct for back-compat, but the new orchestrators don't call it — don't rely on it.) |
| *(nothing — new concept)* | `routes.LocksmithRoutesOptions.Authorizer authenticator.AuthorizerHandler` — **required** |
| *(nothing — new concept)* | `routes.LocksmithRoutesOptions.Registrar register.RegistrarHandler` — **required** |
| *(nothing — new concept)* | `routes.LocksmithRoutesOptions.Bus events.Bus` |
| `sendWelcomeEmailExample(u users.LocksmithUserInterface)` pattern | `events.EventRegistrationSucceeded` payload carries `Email`/`Username`/`UserID` directly — no DB round-trip needed |

If your `main.go` still compiles against current `main` without touching
`routes.InitializeLocksmithRoutes`, it's because `Authorizer`/`Registrar` are
interface-typed fields that default to `nil` — and `main` calls
`options.Authorizer.ServeLoginAPI` and `options.Registrar.ServeRegisterAPI`
**unconditionally** whenever `DisableAPI` is false, and
`options.Authorizer.GetAdditionalLoginMethods()` unconditionally whenever
`DisableUI` is false. Leaving either nil is a **guaranteed nil-pointer panic
at first request**, not a graceful no-op.

---

## 2. New concepts glossary

- **`authenticator.NewAuthorizer(db, opts...)`** (`authentication/authenticator`) —
  the new home for `/api/login`, `/api/login/{provider}`, and
  `/api/login/{provider}/logo`. Built from a list of *methods*
  (`AllowMethodPassword`, `AllowMethodOIDC`, ...). Implements
  `authenticator.AuthorizerHandler`.
- **`register.NewRegistrar(db, opts...)`** (`authentication/register`) — the
  new home for `/api/register`. Built from a list of *methods*
  (`AllowMethodPassword`, `AllowMethodHint`). Implements
  `register.RegistrarHandler`.
- **`events.Bus`** (`authentication/events`) — a pub/sub interface
  (`Publish`/`Subscribe`) used by both orchestrators to report what happened
  (login succeeded/failed, registration succeeded/failed, rostering,
  account linking, email verification, sign-out). `events.NewMemoryBus()` is
  an in-process implementation good enough for the demo and small
  deployments; `events.NoopBus{}` is the zero-value default if you don't
  wire one in. There's no queue/retry — subscribers run synchronously in the
  request path, so keep them fast and non-blocking (fire-and-forget to a
  goroutine/queue if you need to call slow external systems).
- **`tokens.TokenManager`** (`authentication/tokens`) — abstracts how an
  authenticated session is represented and returned to the client.
  `tokens.NewCookieManager(db, redirectPath)` is the cookie-based
  implementation that replaces the old implicit cookie logic. Both the
  authorizer and registrar take the *same* `TokenManager` instance so a
  freshly registered user is logged in consistently with a normal login.
- **`registrationhints.Service`** (`authentication/registrationhints`) — signs
  and verifies short-lived "registration hint" tokens used when an OIDC
  login doesn't match an existing user and the provider is marked
  `Rosterable: true` (i.e., "log in with Google" should implicitly register
  the user). This must share the **same signing package** between the
  authorizer (which issues the hint) and the registrar's hint method (which
  redeems it) — see [§4.7](#47-wire-the-registrar).
- **Method options packages** — `authenticator_methods` and
  `register_methods` hold the `With*`/`Allow*`-style functional options for
  each authentication method (password length, OIDC config, hint service).

---

## 3. Step-by-step migration

Work through these in order — later steps depend on values created in
earlier ones (`sp`, `authEvents`, `tm`).

### 3.1 Update imports

Remove:
```go
"github.com/kvizdos/locksmith/authentication/login"
"github.com/kvizdos/locksmith/authentication/oauth"
oauth_google_oidc "github.com/kvizdos/locksmith/authentication/oauth/oidc"
"github.com/kvizdos/locksmith/authentication/xsrf"
```

Add:
```go
"log/slog"

"github.com/kvizdos/locksmith/authentication/authenticator"
"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_methods"
"github.com/kvizdos/locksmith/authentication/events"
"github.com/kvizdos/locksmith/authentication/register"
"github.com/kvizdos/locksmith/authentication/register/register_methods"
"github.com/kvizdos/locksmith/authentication/registrationhints"
"github.com/kvizdos/locksmith/authentication/tokens"
"github.com/kvizdos/locksmith/authentication/verificationcodes"
google_icon "github.com/kvizdos/locksmith/oidc-icons/oidc-google"
```

`authentication/oauth` is **not** fully gone — it now only exposes small
static-asset handlers (`oauth.KeepAliveJSRoute`, `oauth.GoogleFCMJSRoute`),
which `routes.InitializeLocksmithRoutes` wires up for you. You don't need to
import it directly in `main.go` unless you're doing something custom.

### 3.2 Move signing package construction earlier

The signing package (`sp`) used to be decoded at the very bottom of `main()`,
right before it was assigned to `magic.MagicSigningPackage` and
`xsrf.XSRFSigningPackage`. It's now a **required constructor argument** for
the authorizer (`WithSigningPackage`) and is also shared with the
registrar's hint method, so decode it near the top of `main()`, before you
build anything that depends on it:

```go
sp, _ := signing.DecodePrivateKey("...")
```

Keep the `magic.MagicSigningPackage = &sp` assignment later in `main()` —
that's unrelated to this migration and still required for Magic Tokens.

**Do not** carry over `xsrf.XSRFSigningPackage.Anonymous = &sp` /
`xsrf.XSRFSigningPackage.Authenticated = &sp` — the package no longer
exists. See [§5.3](#53-global-xsrf-signing-package-is-gone).

### 3.3 Stand up an event bus and (optionally) a logger

```go
slog.SetLogLoggerLevel(slog.LevelDebug) // optional, but recommended so authenticator/register slog output is visible

authEvents := events.NewMemoryBus()
```

If you want console visibility into every auth event during development,
subscribe a logger to the events you care about. See the working example in
`subscribePrettyPrintedAuthEvents` in this repo's `main.go` (`main.go:89`) —
it's safe to copy verbatim.

If you skip this, pass `events.NoopBus{}` (or simply don't set
`authenticator.WithEventBus`/`register.WithEventBus`/`routes...Bus` at all —
both orchestrators default to a no-op bus) — but you lose the only
replacement for the old `LoginInfoCallback`/`NewRegistrationEvent` hooks, so
most consumers should wire at least one real subscriber (e.g. your welcome
email, your analytics pipeline).

### 3.4 Build the token manager

```go
tm := tokens.NewCookieManager(db, "/app") // second arg: default post-auth redirect path
```

### 3.5 Replace `RequiresEmailVerification`'s old home... wait, build it first

You'll pass this into the registrar in step 3.7, but it's easiest to define
the closure first since it's identical in signature to before:

```go
requiresEmailVerification := func(ctx context.Context, da database.DatabaseAccessor, lui users.LocksmithUserInterface, validationRes textvalidation.ValidationResultEvaluator) bool {
	_, res := validationRes.Result(true)
	if res != textvalidation.ValidationResult_VALID {
		return true
	}
	return false
}
```

This function's signature is unchanged from `alpha-ui2-v1.2.0`. What changed
is *where it's registered* (`register.WithRequiresEmailVerification(...)`
instead of `routes.LocksmithRoutesOptions.RequiresEmailVerification`).

### 3.6 Build the authorizer

This replaces: `oauth_google_oidc.NewOIDCConnection(...)`,
`routes.LocksmithRoutesOptions.OAuthProviders`, `.LoginSettings`, and
`.LoginInfoCallback`.

```go
googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

authorizer := authenticator.NewAuthorizer(
	db,
	authenticator.WithTokenManager(tm),
	authenticator.WithEventBus(authEvents),

	// Timing-based mitigation for online brute-force/enumeration attacks.
	// Read §5.1 before picking a value — this does NOT replace lockout/captcha.
	authenticator.WithMinimumResponseTime(2*time.Second),

	authenticator.WithSigningPackage(sp), // required if any method can roster new users
	authenticator.WithEmailAsUsername(),  // matches UseEmailAsUsername: true below

	authenticator.WithMethods(
		authenticator.AllowMethodPassword(
			authenticator_methods.RequireMinPasswordLength(8), // keep in sync, see §5.5
		),
		authenticator.AllowMethodOIDC(
			authenticator_methods.WithOIDC(authenticator_methods.OIDConfig{
				Issuer:       "https://accounts.google.com",
				BaseURL:      "https://your-domain.example", // used to build the OAuth2 redirect URL: BaseURL + "/api/login"
				ProviderName: "google",                       // must stay unique across methods
				ClientID:     googleClientID,
				ClientSecret: googleClientSecret,
				LogoBytes:    google_icon.GoogleIcon,
				Rosterable:   true, // if true, unknown Google users are silently sent into registration
			}),
		),
	),
)
```

Notes:
- `authenticator.NewAuthorizer` **panics** at construction time if
  `WithTokenManager`, `WithMethods`, or `WithSigningPackage` weren't
  supplied. Fail fast is intentional — don't wrap this in a `recover()`.
- `authenticator_methods.WithOIDC` calls `oidc.NewProvider` **synchronously**
  during `NewAuthorizer`, and panics if provider discovery fails (bad
  issuer URL, no network, etc.). This mirrors the old
  `oauth_google_oidc.NewOIDCConnection` + `panic(err)` behavior, but it's no
  longer optional — there's no error return to check yourself.
- If you had a second/custom OAuth provider under the old
  `oauth.OAuthProvider` interface, it needs to be re-implemented as an
  `authenticator_domain.Handler` (see `authentication/authenticator/internal/method_oidc`
  for the reference implementation) and passed to `WithMethods`. There is no
  generic adapter from the old interface.

### 3.7 Wire the registrar

This replaces: `routes.LocksmithRoutesOptions.RequiresEmailVerification`,
`.NewRegistrationEvent`.

```go
registrar := register.NewRegistrar(
	db,
	register.WithDefaultRoleName("user"), // REQUIRED — omitting this makes every registration 500
	register.WithEmailAsUsername(true),
	register.WithMinimumLengthRequirement(8), // keep in sync, see §5.5
	register.WithRequiresEmailVerification(requiresEmailVerification),
	register.WithAccountVerifier(verificationcodes.NewVerifier(db, nil)),
	register.WithEventBus(authEvents),
	register.WithTokenManager(tm),
	register.WithMethods(
		// Must share the SAME signing package as authenticator.WithSigningPackage
		// above, or hints issued by the authorizer will fail to verify here.
		register.AllowMethodHint(register_methods.WithHintService(registrationhints.Service{Signer: sp})),
		register.AllowMethodPassword(),
	),
)
```

Gotchas:
- `register.WithDefaultRoleName("user")` must reference a role that exists
  (i.e. created via `roles.CreatePermissionSet`/`roles.AddPermissionsToRole`
  before requests start flowing) — a missing/nonexistent role causes
  `ServeRegisterAPI` to return `500` for every registration attempt, not a
  clear startup error.
- Only add `register.AllowMethodHint(...)` if you have a `Rosterable: true`
  OIDC method on the authorizer. If you don't need "log in with Google to
  auto-register," you can drop this method and the `registrationhints`
  import entirely.

### 3.8 Replace the `NewRegistrationEvent` welcome-email hook with an event subscriber

Old:
```go
func sendWelcomeEmailExample(u users.LocksmithUserInterface) {
	fmt.Printf("Sending welcome email to %s\n", u.GetEmail())
}
// ...
routes.InitializeLocksmithRoutes(mux, db, routes.LocksmithRoutesOptions{
	// ...
	NewRegistrationEvent: sendWelcomeEmailExample,
})
```

New:
```go
func sendWelcomeEmailExampleByEmail(email string) {
	fmt.Printf("Sending welcome email to %s\n", email)
}
// ...
authEvents.Subscribe(events.EventRegistrationSucceeded, func(ctx context.Context, envelope events.Envelope) error {
	payload, ok := envelope.Payload.(events.RegistrationSucceededPayload)
	if !ok {
		return nil
	}
	sendWelcomeEmailExampleByEmail(payload.Email)
	return nil
})
```

`events.RegistrationSucceededPayload` also carries `UserID`, `Username`,
`Method`, `Provider`, `InviteUsed`, `RequiresEmailVerification`,
`AutoLoginIssued`, and `SelectBy` — you likely no longer need to look the
user back up from the DB inside this callback.

If you relied on `LoginInfoCallback` for analytics/audit logging, subscribe
to `events.EventLoginSucceeded` / `events.EventLoginFailed` the same way.
Also consider `events.EventRosterStarted/Succeeded/Failed` (OIDC
auto-registration), `events.EventAccountLinked`, and `events.EventSignOut`
— these didn't have equivalents in the old callback-based API at all.

### 3.9 Update the `routes.InitializeLocksmithRoutes` call

Remove:
```go
RequiresEmailVerification: ...,
LoginInfoCallback:         ...,
OAuthProviders:            []oauth.OAuthProvider{googleOIDC},
LoginSettings:             &login.LoginOptions{ ... },
NewRegistrationEvent:      sendWelcomeEmailExample,
```

Add:
```go
Authorizer: authorizer,
Registrar:  registrar,
Bus:        authEvents,
```

`MinimumPasswordLength: 8` and `HIBPIntegrationOptions` are unaffected —
keep them as-is (but see [§5.5](#55-minimum-password-length-now-lives-in-three-places)).

### 3.10 (Optional, Google-specific) Add the Google Identity Services script tag

The OIDC method now supports Google's newer "FedCM/One Tap" credential flow
in addition to classic redirect-based OAuth2. To enable the credential flow
in your UI, inject the bundled script with your client ID:

```go
Styling: pages.LocksmithPageStyling{
	LogoURL: "/components/locksmith.svg",
	InjectHeader: template.HTML(fmt.Sprintf(
		`<script>console.log("Loaded page.")</script>
		<script src="/api/auth/oauth/google_fcm.js?client_id=%s" async defer></script>`,
		googleClientID,
	)),
},
```

This is additive — the classic redirect flow (`/api/login/google`) keeps
working without it. See [§5.4](#54-google-login-now-has-a-fedcmone-tap-path-with-csrf-implications)
for the security rationale and a required Google Cloud Console change.

### 3.11 Remove the now-invalid `xsrf` assignments

Delete:
```go
xsrf.XSRFSigningPackage.Anonymous = &sp
xsrf.XSRFSigningPackage.Authenticated = &sp
```

There is nothing to replace these with in your `main.go` — CSRF protection
for the new OIDC credential flow is handled internally via a per-browser
nonce cookie (see [§5.4](#54-google-login-now-has-a-fedcmone-tap-path-with-csrf-implications)).
If your **own** custom endpoints (outside Locksmith) read/wrote
`xsrf.XSRFSigningPackage` directly, you'll need to build your own CSRF
mitigation for them — see [§5.3](#53-global-xsrf-signing-package-is-gone).

---

## 4. `routes.LocksmithRoutesOptions` field reference

| Field | Status |
|---|---|
| `Authorizer authenticator.AuthorizerHandler` | **New, required.** Nil causes a panic on first `/api/login*` or `/login` request. |
| `Registrar register.RegistrarHandler` | **New, required.** Nil causes a panic on first `/api/register` or `/register` request. |
| `Bus events.Bus` | **New, optional.** Passed straight to `sign_out.SignOutHTTP{EventBus: options.Bus}` — set it if you want `EventSignOut` fired. |
| `OAuthProviders []oauth.OAuthProvider` | **Removed.** Configure via `authenticator.WithMethods(authenticator.AllowMethodOIDC(...))` instead. |
| `LoginSettings *login.LoginOptions` | **Removed.** No direct replacement; see §5.1. |
| `RequiresEmailVerification func(...)` | **Removed** from this struct. Moved to `register.WithRequiresEmailVerification(...)`. |
| `NewRegistrationEvent func(users.LocksmithUserInterface)` | **Removed.** Moved to `events.EventRegistrationSucceeded` subscriber. |
| `LoginInfoCallback func(method string, user map[string]any)` | Field still exists on the struct but is **not called** by the new authorizer/registrar. Treat as dead — migrate to event subscribers. |
| `RequestLoginInfoCallback func(*http.Request, string, map[string]any)` | Same as above — present, but not invoked by the new code path. |
| `MinimumPasswordLength int` | **Unchanged**, but now one of three places controlling password length — see §5.5. |
| `AccountVerifier`, `EmailValidation`, `SharedMemory`, `SAMLConfig`, `Styling`, `HIBPIntegrationOptions`, `LaunchpadSettings`, `InactivityLockDuration`, `WithErrors`, `ResetPasswordOptions`, `DisablePublicRegistration`, `UseEmailAsUsername`, `OnboardPath`, `InviteUsedRedirect`, `AppName`, `Disable*` flags | **Unchanged.** Carry over as-is. |

---

## 5. Security-relevant behavior changes (read before shipping)

### 5.1 Brute-force / account lockout protection was removed

`alpha-ui2-v1.2.0` shipped `login.LoginOptions.LockoutPolicy`
(`CaptchaAfter`, `LockoutAfter`, `ResetAfter`, `OnLockout`) — a
counting-based defense that required a captcha after N failed attempts and
locked the account after M. **This entire mechanism, and the `login`
package it lived in, has been deleted.** There is no `LockoutPolicy`
equivalent in `authenticator`.

What replaced it:
- `authenticator.WithMinimumResponseTime(d)` — enforces a minimum wall-clock
  time for every `/api/login` response (success or failure), which slows
  down (but does not stop) online brute-forcing and reduces
  timing-side-channel username enumeration.
- `authenticator.DisableUserEnumerationProtection()` — opt-in, disabled by
  default; leaving it disabled means failed logins return a generic
  "Username or Password is incorrect" instead of "Email not found" /
  "Password is incorrect" (also mitigates enumeration).
- The `shared-memory/objects.UserLoginAttempt` machinery (attempt counting)
  still exists in the codebase but is **no longer wired into the
  authorizer** — it's dead code from the authorizer's perspective as of this
  version.

**Action for you:** if your threat model depended on hard lockouts and
captchas (e.g. compliance requirement, previously-mitigated credential
stuffing), you must now provide this yourself — options include:
- A reverse proxy / WAF / API gateway rate limit in front of `/api/login`
  and `/api/login/{provider}` (per-IP and, if feasible, per-username).
  `ratelimits.NewRateLimiter` is available in this codebase and used
  elsewhere (see `TestAppHandler` route in `main.go`), but it's on you to
  apply it to `/api/login` — `InitializeLocksmithRoutes` does not do this
  for you.
- Subscribing to `events.EventLoginFailed` and implementing your own
  counting + temporary lock (e.g. flip a `locked` flag on the user record,
  or feed a rate-limiter keyed by presented username/IP).
- Re-enabling a captcha via `captchaproviders.UseProvider` (this hook still
  exists) in front of login, independent of attempt counts.

Do not treat `WithMinimumResponseTime` alone as equivalent protection — it
raises the cost of brute-forcing, it doesn't cap it.

### 5.2 Passwordless (magic/OAuth-only) accounts can no longer trigger password reset

A later commit on `main` ("DiD: Don't allow passwordless to reset")
hardens the reset-password flow against issuing password reset flows for
accounts that don't have passwords. If your app or support tooling relied
on being able to "reset" a purely-OAuth or magic-link account into having a
password, verify that flow still behaves the way you expect after
upgrading, and update any documentation/support scripts accordingly.

### 5.3 Global XSRF signing package is gone

`authentication/xsrf` (including `xsrf.XSRFSigningPackage`) was deleted.
Frontend templates (`pages/login2.html`, `pages/launchpad*.html`,
`components/signin.component.js`) still reference an `xsrf`/`LoginXSRF`
field, but **nothing on the server populates it anymore** — it will render
as an empty string. This is expected; don't spend time trying to "fix" a
missing `LoginXSRF` value.

If you built custom, non-Locksmith endpoints that reused
`xsrf.XSRFSigningPackage` for your own CSRF protection, that global no
longer exists at all (not deprecated — deleted, will fail to compile).
You'll need your own CSRF mechanism for those endpoints (e.g. Go's
`net/http` `SameSite` cookies plus a double-submit token, or a
purpose-built middleware).

### 5.4 Google login now has a FedCM/One-Tap path, with CSRF implications

The new `method_oidc` handler supports both:
1. The classic server-redirect OAuth2 flow (`/api/login/{provider}` →
   provider → `/api/login` callback), functionally similar to the old
   `oauth_google_oidc` package.
2. A newer Google Identity Services credential flow ("One Tap"/FedCM),
   enabled by serving `/api/auth/oauth/google_fcm.js` (see §3.10). This
   flow is protected against **login CSRF** — a forged cross-site POST that
   would otherwise log a victim into an attacker-controlled account — by
   binding a per-browser, HttpOnly nonce cookie (`ls_oidc_fcm_nonce`,
   `SameSite=Lax`, `Secure`) to the credential's `nonce` claim, verified
   server-side before the credential is trusted.

**Action for you:**
- If you adopt the `google_fcm.js` script tag, add your production
  domain(s) as **Authorized JavaScript origins** (not just Authorized
  redirect URIs) in the Google Cloud Console OAuth client configuration —
  FedCM/One Tap requires this in addition to whatever redirect URI you
  already configured for the classic flow.
- Because the nonce cookie is `Secure`, this flow will not function over
  plain HTTP — confirm your local dev setup either uses HTTPS or that you
  only exercise the classic redirect flow locally.

### 5.5 Minimum password length now lives in three places

Previously, `routes.LocksmithRoutesOptions.MinimumPasswordLength` was the
single source of truth. Now there are three independent settings that
should be kept equal:

| Setting | Effect |
|---|---|
| `routes.LocksmithRoutesOptions.MinimumPasswordLength` | Drives the *displayed* minimum length on the register/reset-password pages, and enforcement on `/api/reset-password`. |
| `register.WithMinimumLengthRequirement(n)` | Enforced during `/api/register`. |
| `authenticator_methods.RequireMinPasswordLength(n)` | Enforced during `/api/login` (re-validated on every login per the HIBP/continuous-verification feature). |

Setting these to different values won't cause an error — it will cause
confusing, inconsistent UX and inconsistent enforcement (e.g., UI says 8,
but login enforces 10). Pick one value and set it in all three places.

### 5.6 OIDC-rostered ("auto-register via Google") users go through a signed hint token

When an OIDC login doesn't match an existing user and the method is
`Rosterable: true`, the authorizer no longer improvises — it issues a
short-lived (~60s), signed `registrationhints` token (via the shared
signing package, `sp`), sets it as a cookie, and redirects to
`/api/register?hinted`. The registrar's `AllowMethodHint` method verifies
this signed token before creating the account. This closes a class of bugs
where an unauthenticated request could otherwise claim to be "the user
Google just verified." If you add other `Rosterable: true` OIDC providers,
double-check `WithSigningPackage(sp)` on the authorizer and
`register_methods.WithHintService(registrationhints.Service{Signer: sp})`
on the registrar use the exact same signing package instance, or hint
verification will always fail.

---

## 6. Non-breaking but worth adopting

- `slog.SetLogLoggerLevel(slog.LevelDebug)` — the new packages use
  structured `slog` logging extensively (`a.log.DebugContext`, etc.); turn
  this up during migration testing to see what's happening on every
  login/registration attempt, then dial back to `Info`/`Warn` for
  production.
- `subscribePrettyPrintedAuthEvents` — a ready-made console logger for
  every event name the bus supports; useful scaffolding even if you replace
  it with real telemetry later.
- `map[string]interface{}` → `map[string]any` in the launchpad `Custom`
  fields — purely cosmetic (Go 1.18+ alias), no behavior change; update if
  you want to match upstream style, not required.

---

## 7. Post-migration validation checklist

- [ ] `go build ./...` succeeds with no references to
      `authentication/login`, `authentication/oauth/oidc`, or
      `authentication/xsrf` remaining in your code.
- [ ] `go vet ./...` clean.
- [ ] Password login succeeds and fails correctly (wrong password, unknown
      user) and returns the expected generic-vs-specific error message
      based on your `DisableUserEnumerationProtection()` choice.
- [ ] Google OAuth classic redirect flow logs in an existing linked user.
- [ ] Google OAuth logs in with a **new** Google account against a
      `Rosterable: true` provider and correctly redirects into
      `/api/register?hinted`, and completes registration.
- [ ] If you added `google_fcm.js`: One Tap / FedCM prompt appears and
      completes login over HTTPS; confirm Google Cloud Console has your
      origin listed under Authorized JavaScript origins.
- [ ] Registration with a role that does **not** yet exist in `roles.*`
      fails loudly in your own testing (before prod) rather than silently
      500ing in front of users — verify `WithDefaultRoleName` matches a
      role created via `roles.CreatePermissionSet`/`AddPermissionsToRole`.
- [ ] `events.EventRegistrationSucceeded` subscriber fires and welcome
      email/analytics hook still works.
- [ ] `events.EventLoginSucceeded`/`EventLoginFailed` subscribers (if any)
      fire as expected.
- [ ] Password length is identical across `MinimumPasswordLength`,
      `register.WithMinimumLengthRequirement`, and
      `authenticator_methods.RequireMinPasswordLength`.
- [ ] Explicit decision made and documented for brute-force protection
      (rate limiting, captcha, or custom lockout) now that `LockoutPolicy`
      no longer exists — do not ship without one.
- [ ] Password-reset flow rejects passwordless (magic/OAuth-only) accounts
      as expected (§5.2).
- [ ] Sign-out (`/sign-out`) still works and, if wired, publishes
      `events.EventSignOut`.
