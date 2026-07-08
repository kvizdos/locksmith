package authenticator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/api_helpers"
	"github.com/kvizdos/locksmith/authentication/authenticator/authenticator_domain"
	"github.com/kvizdos/locksmith/authentication/events"
	"github.com/kvizdos/locksmith/authentication/registrationhints"
	"github.com/kvizdos/locksmith/authentication/tokens"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/users"
)

func (a *authorizers) publishAuthEvent(ctx context.Context, name events.EventName, payload any) {
	if a.eventBus == nil {
		return
	}
	envelope := events.EnrichEnvelope(ctx, events.Envelope{
		ID:         uuid.New().String(),
		Name:       name,
		OccurredAt: time.Now(),
		Source:     "authentication/authenticator",
		Payload:    payload,
	})
	if err := a.eventBus.Publish(ctx, envelope); err != nil {
		a.log.WarnContext(ctx, "auth event publish failed", "event", string(name), "error", err)
	}
}

func (a *authorizers) writeAuthError(minCtx func(), w http.ResponseWriter, protectionDisabled string) {
	if minCtx != nil {
		minCtx()
	}
	if a.disableUserEnumerationProtection {
		api_helpers.WriteResponse(w, map[string]any{
			"error": protectionDisabled,
		}, http.StatusUnauthorized)
		return
	}

	api_helpers.WriteResponse(w, map[string]string{
		"error": "Username or Password is incorrect.",
	}, http.StatusUnauthorized)
}

func (a *authorizers) enrichCtx(r *http.Request) context.Context {
	ctx := r.Context()
	return context.WithValue(ctx, "ip", logger.GetIPFromRequest(*r))
}

func (a *authorizers) getUsernameNoun() string {
	if a.emailAsUsername {
		return "Email"
	}
	return "Username"
}

func (a *authorizers) beginRegistrationRostering(w http.ResponseWriter, r *http.Request, hint *registrationhints.Hint) {
	token, err := registrationhints.Service{Signer: a.sp}.Create(*hint)
	if err != nil {
		a.log.ErrorContext(r.Context(), "failed to create registration hint")
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "Failed to start registration.",
		}, http.StatusInternalServerError)
		return
	}

	registrationhints.SetCookie(w, token)

	a.publishAuthEvent(r.Context(), events.EventRosterStarted, events.RosterStartedPayload{Provider: hint.ProviderName, SelectBy: hint.SelectBy})

	http.Redirect(w, r, "/api/register?hinted", http.StatusSeeOther)
	// api_helpers.WriteResponse(w, map[string]any{
	// 	"mode": "roster",
	// }, http.StatusAccepted)
}

func (a *authorizers) GetAdditionalLoginMethods() []string {
	var methods []string
	for _, method := range a.methods {
		if method.Passwordless() {
			methods = append(methods, method.Name())
		}
	}
	return methods
}

func (a *authorizers) ServeLoginAPI(w http.ResponseWriter, r *http.Request) {
	ctx := a.enrichCtx(r)
	var waitMinimum func()

	if a.minimumResponseTime > 0 {
		c, cancel := context.WithTimeout(ctx, a.minimumResponseTime)
		defer cancel()
		waitMinimum = func() {
			deadline, ok := c.Deadline()
			if ok {
				a.log.DebugContext(ctx, "waiting for minimum response time", "remaining", time.Until(deadline))
			}
			<-c.Done()
		}
	} else {
		waitMinimum = func() {}
	}

	// First, detect which handler this session will
	// utilize. Each handler runs its own heuristic to
	// determine if it can handle this request.
	handler, err := a.getHandler(r)
	if err != nil {
		a.log.ErrorContext(ctx, "no handler supports this request", "error", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "no handler supports this request",
		}, http.StatusInternalServerError)
		return
	}

	ctx = context.WithValue(ctx, "login_handler", handler.Name())
	ctx = context.WithValue(ctx, "login_passwordless", handler.Passwordless())

	token, lastSelectBy, err := a.attemptLogin(ctx, handler, r)
	if err != nil {
		var (
			invalidPassword *authenticator_domain.InvalidPasswordError
			invalidUser     *authenticator_domain.UserNotFoundError
		)
		switch {
		case errors.As(err, &invalidUser):
			if invalidUser.RegistrationHint != nil && invalidUser.RegistrationHint.Rosterable {
				a.log.InfoContext(ctx, "starting registration rostering", "source", invalidUser.RegistrationHint.ProviderName)
				a.beginRegistrationRostering(w, r, invalidUser.RegistrationHint)
				return
			}
			a.log.InfoContext(ctx, "user not found", "presented_user_id", invalidUser.PresentedUsername)
			a.publishAuthEvent(ctx, events.EventLoginFailed, events.LoginFailedPayload{
				PresentedUsername: invalidUser.PresentedUsername,
				Method:            handler.Name(),
				Reason:            "user_not_found",
			})
			a.writeAuthError(waitMinimum, w, fmt.Sprintf("%s not found.", a.getUsernameNoun()))
			return
		case errors.As(err, &invalidPassword):
			a.log.InfoContext(ctx, "invalid password presented", "user", invalidPassword.Username)
			a.publishAuthEvent(ctx, events.EventLoginFailed, events.LoginFailedPayload{
				PresentedUsername: invalidPassword.Username,
				Method:            handler.Name(),
				Reason:            "invalid_password",
			})
			a.writeAuthError(waitMinimum, w, "Password is incorrect.")
			return
		case errors.Is(err, authenticator_domain.ErrPasswordlessRequired):
			a.publishAuthEvent(ctx, events.EventLoginFailed, events.LoginFailedPayload{
				Method: handler.Name(),
				Reason: "passwordless_required",
			})
			a.writeAuthError(waitMinimum, w, "Passwordless login required.")
			return
		case errors.Is(err, authenticator_domain.ErrFailedToParse):
			a.log.DebugContext(ctx, "invalid request", "err", err)
			api_helpers.WriteResponse(w, api_helpers.APIResponseError{
				Reason: "invalid request body",
			}, http.StatusBadRequest)
			return
		case errors.Is(err, authenticator_domain.ErrMethodNotSupported):
			a.log.DebugContext(ctx, "unsupported method", "error", err, "stage", "attempt_login", "handler", fmt.Sprintf("%T", handler))
			api_helpers.WriteResponse(w, api_helpers.APIResponseError{
				Reason: "unsupported method",
			}, http.StatusMethodNotAllowed)
			return
		}

		a.log.ErrorContext(ctx, "failed to attempt login", "err", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "an unexpected error occurred",
		}, http.StatusInternalServerError)
		return
	}
	if token == nil {
		a.log.ErrorContext(ctx, "token is missing, yet there was no error")
		a.writeAuthError(waitMinimum, w, fmt.Sprintf("%s or Password is incorrect.", a.getUsernameNoun()))
		return
	}

	if _, ok := handler.(authenticator_domain.Beginnable); ok {
		a.log.DebugContext(ctx, "setting oauth hints")
		token.OAuthHint = token.User.Email
		token.OAuthProvider = handler.Name()
	}

	// Some Locksmith functionality requires specific cookies
	// to be set..
	if err := a.setBaseCookies(w, token); err != nil {
		a.log.ErrorContext(ctx, "failed to set base cookies", "error", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "something went wrong",
		}, http.StatusInternalServerError)
		return
	}

	if err := a.tm.PassToClient(w, r, token); err != nil {
		a.log.ErrorContext(ctx, "failed to pass token to client", "error", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "something went wrong",
		}, http.StatusInternalServerError)
		return
	}

	userID := ""
	if token.User != nil {
		userID = token.User.GetID()
	}
	a.publishAuthEvent(ctx, events.EventLoginSucceeded, events.LoginSucceededPayload{
		UserID:       userID,
		Method:       handler.Name(),
		Provider:     token.OAuthProvider,
		Passwordless: handler.Passwordless(),
		SelectBy:     lastSelectBy,
	})

	// Pass to client is expected to finish.
}

func (a *authorizers) attemptLogin(ctx context.Context, handler authenticator_domain.Handler, r *http.Request) (*authenticator_domain.Token, string, error) {
	// Next, create a session for this handler.
	// This session will include handler-specific state.
	session := handler.Session(a.db)

	// Load the request into the session to initialize it.
	// This will parse the request and set any handler-specific state.
	if err := session.LoadRequest(r); err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	// If the session can identify which UI surface produced the
	// credential (e.g. Google Identity Services' "select_by"), capture it
	// now so every return path below can report it.
	selectBy := ""
	if fs, ok := session.(authenticator_domain.FlowSource); ok {
		selectBy = fs.GetSelectBy()
	}

	// If required, resolve the identity of the user
	// from the session. This is particularly
	// useful for OAuth login scenarios.
	if resolver, ok := session.(authenticator_domain.IdentityResolver); ok {
		if err := resolver.ResolveIdentity(ctx); err != nil {
			return nil, selectBy, fmt.Errorf("resolve identity: %w", err)
		}
	}

	// Request the User ID from the session state.
	// This will return the user ID and any handler-specific state.
	var rawUser any
	var found bool
	userAlreadyFetched := false

	lookup := "username"
	if a.emailAsUsername {
		lookup = "email"
	}
	id := session.GetPresentedUser()

	if fi, ok := session.(authenticator_domain.FederatedIdentity); ok {
		provider := fi.GetProvider()
		subject := fi.GetSubject()

		a.log.DebugContext(ctx, "checking federated identity", "provider", provider, "subject", subject)

		rawUser, found = a.db.FindOne("auth_links", map[string]any{
			"provider": provider,
			"subject":  subject,
		})

		if found {
			linkedIdentity := authenticator_domain.LinkedIdentityFromMap(rawUser.(map[string]any))
			a.log.DebugContext(ctx, "linked identity found", "issuer", linkedIdentity.Issuer, "subject", linkedIdentity.Subject, "user_id", linkedIdentity.UserID)

			if linkedIdentity.Issuer != fi.GetIssuer() {
				a.log.WarnContext(ctx, "linked identity issuer does not match session issuer", "issuer", linkedIdentity.Issuer, "session_issuer", fi.GetIssuer())
				return nil, selectBy, &authenticator_domain.UserNotFoundError{
					PresentedUsername: session.GetPresentedUser(),
				}
			}

			lookup = "id"
			id = linkedIdentity.UserID
		} else if vc, ok := session.(authenticator_domain.VerifiedContact); ok && vc.EmailVerified() {
			// No existing link, but a verified contact point is available —
			// attempt to auto-link to an existing user account.
			email := vc.GetEmail()
			a.log.DebugContext(ctx, "no linked identity found, attempting auto-link by email", "provider", provider, "subject", subject, "email", email)
			if autoLinkUser, autoLinkFound := a.db.FindOne("users", map[string]any{"email": email}); autoLinkFound {
				autoLinkMap := autoLinkUser.(map[string]any)
				userID, _ := autoLinkMap["id"].(string)
				if err := a.LinkAccount(ctx, userID, provider, fi.GetIssuer(), subject); err != nil {
					a.log.ErrorContext(ctx, "failed to auto-link account", "provider", provider, "subject", subject, "user_id", userID, "err", err)
				} else {
					a.log.InfoContext(ctx, "auto-linked account", "provider", provider, "subject", subject, "user_id", userID)
					rawUser = autoLinkUser
					found = true
					userAlreadyFetched = true
				}
			} else {
				a.log.DebugContext(ctx, "no user found for auto-link email", "email", email)
				id = "" // ensure nothing is found
			}
		} else {
			a.log.DebugContext(ctx, "no linked identity and no verified contact, denying", "provider", provider, "subject", subject)
			id = "" // ensure nothing is found
		}
	}

	if !userAlreadyFetched && id != "" {
		a.log.DebugContext(ctx, "checking user identity", "lookup", lookup, "id", id)

		rawUser, found = a.db.FindOne("users", map[string]any{
			lookup: id,
		})
	}

	if !found {
		a.log.DebugContext(ctx, "user not found", "lookup", lookup, "id", id)
		err := &authenticator_domain.UserNotFoundError{
			PresentedUsername: session.GetPresentedUser(),
		}
		if r, ok := session.(authenticator_domain.Rosterable); ok {
			err.RegistrationHint = r.RegistrationHint()
		}
		return nil, selectBy, err
	}

	rawMap, ok := rawUser.(map[string]any)
	if !ok {
		return nil, selectBy, fmt.Errorf("invalid user result type %T", rawUser)
	}

	var user users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&user, rawMap)
	if user == nil {
		return nil, selectBy, errors.New("failed to read user from map")
	}

	// Confirm the user is authorized to use this handler.
	if user.Passwordless() && !handler.Passwordless() {
		return nil, selectBy, fmt.Errorf("handler %q does not support passwordless: %w", handler.Name(), authenticator_domain.ErrPasswordlessRequired)
	}

	if lu, ok := user.(interface{ GetOAuthRestrictedSource() string }); ok {
		if source := lu.GetOAuthRestrictedSource(); source != "" && source != handler.Name() {
			a.log.WarnContext(ctx, "user attempted login via restricted provider", "allowed", source, "attempted", handler.Name(), "user", user.GetID())
			return nil, selectBy, &authenticator_domain.UserNotFoundError{
				PresentedUsername: session.GetPresentedUser(),
			}
		}
	}

	// Confirm the user is authorized...
	err := session.IsAuthorized(user)
	if err != nil {
		if errors.Is(err, authenticator_domain.ErrInvalidPassword) {
			return nil, selectBy, &authenticator_domain.InvalidPasswordError{
				Username: user.GetUsername(),
				UserID:   user.GetID(),
			}
		}
		return nil, selectBy, fmt.Errorf("failed to validate: %w", err)
	}

	// User is authorized, create the token.
	// This can be a cookie, a JWT, or any other method supported by the token manager.
	token, err := a.tm.CreateAuthToken(user)
	if err != nil {
		return nil, selectBy, fmt.Errorf("failed to create auth token: %w", err)
	}

	// If the session carries a caller-supplied "return to this page after
	// login" target, hand it to the TokenManager so it can redirect there
	// instead of its configured default.
	if rs, ok := session.(authenticator_domain.RedirectSource); ok {
		if redirectTarget := rs.GetRedirectTarget(); redirectTarget != "" {
			token.RedirectPath = redirectTarget
		}
	}

	a.log.InfoContext(ctx, "user logged in", "lookup", lookup, "id", id, "handler", handler.Name(), "user", user.GetID())

	return token, selectBy, nil
}

func (a *authorizers) setBaseCookies(w http.ResponseWriter, token *authenticator_domain.Token) error {
	tokens.SetBaseCookies(w, token)
	return nil
}
