---
name: locksmith-secure-upgrade
description: Migrate a Go application's Locksmith wiring (main.go-style setup) from the alpha-ui2-v1.2.0 authentication API to the current authenticator/register/events API, without silently dropping security controls. Use when a user asks to "upgrade Locksmith", "migrate off alpha-ui2-v1.2.0", or reports compile errors after updating the kvizdos/locksmith module referencing authentication/login, authentication/xsrf, or authentication/oauth/oidc.
---

# Locksmith Secure Upgrade Skill

You are performing a **security-relevant** dependency migration, not a
cosmetic refactor. The old (`alpha-ui2-v1.2.0`) authentication wiring and the
new (`main`) wiring are not behaviorally equivalent — some old protections
(brute-force lockout) have no automatic replacement, and some new protections
(FedCM nonce binding, hint-token signing) only work if wired correctly. Speed
is not the priority here; correctness and an explicit security decision trail
are.

Full technical background lives in `Upgrade.md` (same directory as this
skill, in the `webauthn-bootstrap` / `kvizdos/locksmith` repo). Read it fully
before touching code — this skill file is the *procedure*, `Upgrade.md` is
the *reference*.

If `Upgrade.md` is not present in the user's project, fetch/read it from the
`kvizdos/locksmith` repository first, or ask the user for it. Do not attempt
this migration from memory of this skill alone.

## When to use this skill

Trigger this skill when:
- The user asks to upgrade/migrate their Locksmith integration off
  `alpha-ui2-v1.2.0` (or any pre-rewrite tag) to current `main`.
- The user shows compile errors referencing `authentication/login`,
  `authentication/oauth/oidc`, `authentication/xsrf`,
  `login.LoginOptions`, `login.LockoutPolicy`, or
  `oauth_google_oidc.NewOIDCConnection` after bumping their `go.mod`
  reference to `kvizdos/locksmith`.
- The user asks "why did my Locksmith login/registration stop compiling
  after updating."

Do not use this skill for unrelated Locksmith bugs on an already-migrated
(post-rewrite) codebase.

## Ground rules

1. **Never delete a security control without an explicit, recorded
   replacement decision.** If the user's old `main.go` had a
   `login.LoginOptions{LockoutPolicy: ...}` block, you must not just delete
   it and move on — see Step 5.
2. **Never fabricate secrets, client IDs, or signing keys.** If the user's
   existing code hardcodes a real (non-placeholder) OAuth client secret or
   signing key, flag it and recommend moving it to an environment variable
   as part of this migration — don't just copy it verbatim into new code
   without comment, and don't invent one if one isn't present.
3. **Preserve behavior the user didn't ask you to change.** Password length
   requirements, role names, redirect paths, HIBP settings, SAML config,
   launchpad config, etc. should carry over unchanged unless the user asks
   otherwise.
4. **Compile and smoke-test before declaring done.** `go build ./...` (and
   `go vet ./...`) must pass. If the project has tests covering auth flows,
   run them. If it doesn't, say so explicitly rather than implying coverage
   that doesn't exist.
5. **Surface every item in Upgrade.md §5 to the user by name**, even if you
   think the answer is obvious for their app. These are the
   security-relevant behavior changes (brute-force protection removal, XSRF
   package removal, FedCM CSRF binding, passwordless-reset hardening,
   password-length triad, hint-token signing). Do not silently pick a
   default on the user's behalf for §5.1 (brute-force protection) — this is
   the one item in the whole migration with no automatic replacement, and
   picking wrong has real security consequences.

## Procedure

### Step 1 — Locate and read the target code

1. Find the user's Locksmith wiring code (usually a `main.go` calling
   `routes.InitializeLocksmithRoutes`, but confirm — some apps wrap this in
   their own `internal/auth` package).
2. Read the whole file(s) involved before editing anything.
3. Identify which optional old features are actually in use:
   - `login.LoginOptions`/`LockoutPolicy`? (→ Step 5 is mandatory)
   - `xsrf.XSRFSigningPackage` referenced anywhere outside Locksmith's own
     setup (i.e., in the user's *own* handlers)? (→ flag, needs bespoke
     replacement, out of scope for automatic migration)
   - How many OAuth/OIDC providers, and are any *not* Google (i.e., custom
     `oauth.OAuthProvider` implementations)? (→ Step 4 needs a
     hand-written `authenticator_domain.Handler` per custom provider —
     do not attempt to silently drop a provider)
   - Is `RequiresEmailVerification` doing anything non-trivial (side
     effects, external calls)? Preserve it exactly, just relocate it.
   - Is `NewRegistrationEvent`/`LoginInfoCallback` doing anything beyond
     logging (e.g., sending real emails, writing to a real analytics
     system)? Preserve the call, just relocate it to an event subscriber.

### Step 2 — Confirm module version

Check the user's `go.mod` for their `kvizdos/locksmith` reference (or, if
this *is* the locksmith module itself, confirm which commit/tag they're
building from). If they're not actually on a post-rewrite commit yet
(i.e., `authentication/authenticator` package doesn't exist in their
vendored/downloaded module version), the migration steps below will not
compile — tell the user they need to bump the module version first, and
confirm which version/commit to target (ask, if not stated by the user).

### Step 3 — Rebuild imports

Remove imports for `authentication/login`, `authentication/oauth/oidc`,
`authentication/xsrf`, and any `oauth.OAuthProvider`-based custom provider
imports (unless the provider will be reimplemented — see Step 4). Add
imports per `Upgrade.md` §3.1, adjusted for only the methods the user
actually uses (e.g., skip `registrationhints`/`register_methods` entirely
if the user has no `Rosterable: true` OIDC provider).

### Step 4 — Rebuild the authenticator and registrar

Follow `Upgrade.md` §3.2–3.7 in order:
1. Move signing package decoding earlier if a signing package is needed
   (only required if any OIDC method is `Rosterable: true`, or the app
   uses Magic Tokens — check for `magic.MagicSigningPackage` usage).
2. Build (or reuse) an `events.Bus`.
3. Build a `tokens.TokenManager`.
4. Build `authenticator.NewAuthorizer(...)` with every method the user had
   before (password + each OAuth/OIDC provider). For non-Google custom
   providers, this requires writing a new `authenticator_domain.Handler` —
   this is real implementation work, not config; scope it explicitly with
   the user rather than attempting a blind mechanical translation.
5. Build `register.NewRegistrar(...)`, carrying over
   `RequiresEmailVerification` and role/length settings unchanged.
6. If a `Rosterable: true` OIDC method exists, wire
   `register.AllowMethodHint` with the **same** signing package instance
   used by the authorizer. Double-check this — it's the single most common
   way to break the "log in with Google auto-registers you" flow silently
   (verification fails, user sees a generic error, nothing crashes).

### Step 5 — Handle brute-force/lockout protection explicitly (mandatory decision point)

Before finishing, check whether the user's old code set
`LoginSettings.LockoutPolicy`. If it did:

1. Tell the user plainly: this policy has no automatic equivalent, and the
   underlying package (`authentication/login`) was deleted upstream.
2. Present the options from `Upgrade.md` §5.1 (reverse-proxy/WAF rate
   limiting, a custom `events.EventLoginFailed`-driven lockout, or
   re-enabling `captchaproviders.UseProvider`).
3. Ask the user which they want, or recommend
   `authenticator.WithMinimumResponseTime(...)` **plus** at least one of the
   above — do not present `WithMinimumResponseTime` alone as sufficient
   unless the user explicitly accepts that tradeoff after being told it
   only slows, not stops, brute-forcing.
4. Implement whatever the user chooses. If they explicitly decline any
   additional protection beyond `WithMinimumResponseTime`, record that
   decision in a code comment near the `authenticator.NewAuthorizer(...)`
   call so it's visible to the next reader, e.g.:
   ```go
   // NOTE: LockoutPolicy from the pre-rewrite Locksmith API has no
   // replacement here. Team decided (2026-xx-xx) to rely on
   // WithMinimumResponseTime + upstream WAF rate limiting instead of an
   // in-app counting lockout. See Upgrade.md §5.1.
   ```

If the old code never configured `LockoutPolicy` (defaults only), this step
is still worth a one-line mention to the user, but doesn't block the
migration.

### Step 6 — Update `routes.InitializeLocksmithRoutes` call site

Per `Upgrade.md` §3.9 / §4. Remove the deleted/dead fields, add
`Authorizer`, `Registrar`, and `Bus`. Grep the user's codebase for any other
callers of these removed fields (some apps construct
`LocksmithRoutesOptions` in more than one place, e.g. tests).

### Step 7 — Relocate event-driven side effects

For every `NewRegistrationEvent`/`LoginInfoCallback`/`RequestLoginInfoCallback`
usage found in Step 1, add an equivalent `bus.Subscribe(...)` call using the
matching `events.Event*` constant and payload struct (`Upgrade.md` §3.8).
Preserve the exact external behavior (same email content, same analytics
event name/fields) — only the wiring mechanism should change.

### Step 8 — Sync password length settings

Grep for every place minimum password length is configured
(`MinimumPasswordLength`, any hardcoded `8`/other constant near
password-related config) and set `routes.LocksmithRoutesOptions.MinimumPasswordLength`,
`register.WithMinimumLengthRequirement`, and
`authenticator_methods.RequireMinPasswordLength` to the **same** value
(`Upgrade.md` §5.5).

### Step 9 — Google FedCM/One Tap (only if requested)

Only add the `google_fcm.js` script tag / `InjectHeader` change (`Upgrade.md`
§3.10) if the user asks for it or already had equivalent "Sign in with
Google" one-tap UX. If you add it, explicitly remind the user about the
Google Cloud Console "Authorized JavaScript origins" requirement
(`Upgrade.md` §5.4) — this is an out-of-repo configuration step you cannot
perform yourself; call it out as a manual action item.

### Step 10 — Build, vet, and self-review

1. Run `go build ./...` and `go vet ./...` (or the project's equivalent) in
   the affected module. Fix compile errors before proceeding.
2. Re-read the final diff against `Upgrade.md` §4's field table — confirm
   no removed field is still referenced, and both `Authorizer` and
   `Registrar` are set on every `LocksmithRoutesOptions{}` construction
   site you touched.
3. Confirm `register.WithDefaultRoleName(...)` references a role that is
   actually created somewhere in the app's startup path
   (`roles.CreatePermissionSet`/`AddPermissionsToRole`), not just assumed
   to exist.

### Step 11 — Report to the user

In your final summary, explicitly list:
- Every file changed.
- The brute-force protection decision made in Step 5 (or the fact that it's
  still outstanding and needs the user's input if they haven't answered
  yet — do not silently default this).
- Any custom OAuth provider that still needs a hand-written
  `authenticator_domain.Handler` (if you couldn't fully migrate it
  automatically).
- The full checklist from `Upgrade.md` §7, marked as done/not-done/not
  -applicable, and which items still require **manual** testing (OAuth
  flows, FedCM, email delivery) that you could not execute yourself.

## Anti-patterns to avoid

- Do not comment out or delete `LockoutPolicy`/lockout-related code without
  discussing the replacement with the user first (Step 5 is not optional).
- Do not leave `Authorizer`/`Registrar` unset "to make it compile" — both
  are required; a nil value compiles fine and panics at runtime on first
  request, which is worse than a compile error.
- Do not reuse two different signing-package instances between the
  authorizer and the hint-registration method "because it compiles either
  way" — it compiles but silently breaks OIDC auto-registration.
- Do not add the `google_fcm.js` tag without also telling the user about
  the required Google Cloud Console change — a silently broken FedCM prompt
  is a worse outcome than not offering it.
- Do not assume `map[string]interface{}` → `map[string]any` needs to
  happen; it's purely cosmetic and out of scope unless the user asks for
  general modernization.
