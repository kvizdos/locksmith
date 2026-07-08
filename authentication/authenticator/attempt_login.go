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
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/users"
)

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

func (a *authorizers) beginRegistrationRostering(w http.ResponseWriter, r *http.Request, hint *authenticator_domain.RegistrationHint) {
	hint.ID = uuid.NewString()
	token, err := a.sp.CreateJWT(hint)
	if err != nil {
		a.log.ErrorContext(r.Context(), "failed to create registration hint jwt", "err", err)
		api_helpers.WriteResponse(w, api_helpers.APIResponseError{
			Reason: "Failed to create JWT.",
		}, http.StatusInternalServerError)
		return
	}

	hintCookie := http.Cookie{
		Name:     "registration_hint",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/register",
		MaxAge:   60,
	}
	http.SetCookie(w, &hintCookie)

	api_helpers.WriteResponse(w, map[string]any{
		"mode": "roster",
	}, http.StatusAccepted)
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

	token, err := a.attemptLogin(ctx, handler, r)
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
			a.writeAuthError(waitMinimum, w, fmt.Sprintf("%s not found.", a.getUsernameNoun()))
			return
		case errors.As(err, &invalidPassword):
			a.log.InfoContext(ctx, "invalid password presented", "user", invalidPassword.Username)
			a.writeAuthError(waitMinimum, w, "Password is incorrect.")
			return
		case errors.Is(err, authenticator_domain.ErrPasswordlessRequired):
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

	// Pass to client is expected to finish.
}

func (a *authorizers) attemptLogin(ctx context.Context, handler authenticator_domain.Handler, r *http.Request) (*authenticator_domain.Token, error) {
	// Next, create a session for this handler.
	// This session will include handler-specific state.
	session := handler.Session(a.db)

	// Load the request into the session to initialize it.
	// This will parse the request and set any handler-specific state.
	if err := session.LoadRequest(r); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// If required, resolve the identity of the user
	// from the session. This is particularly
	// useful for OAuth login scenarios.
	if resolver, ok := session.(authenticator_domain.IdentityResolver); ok {
		if err := resolver.ResolveIdentity(ctx); err != nil {
			return nil, fmt.Errorf("resolve identity: %w", err)
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
				return nil, &authenticator_domain.UserNotFoundError{
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
		return nil, err
	}

	rawMap, ok := rawUser.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid user result type %T", rawUser)
	}

	var user users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&user, rawMap)
	if user == nil {
		return nil, errors.New("failed to read user from map")
	}

	// Confirm the user is authorized to use this handler.
	if user.Passwordless() && !handler.Passwordless() {
		return nil, fmt.Errorf("handler %q does not support passwordless: %w", handler.Name(), authenticator_domain.ErrPasswordlessRequired)
	}

	if user.RequiresEmailVerification() {
		return nil, fmt.Errorf("account email not verified: %w", authenticator_domain.ErrPasswordlessRequired)
	}

	if lu, ok := user.(interface{ GetOAuthRestrictedSource() string }); ok {
		if source := lu.GetOAuthRestrictedSource(); source != "" && source != handler.Name() {
			a.log.WarnContext(ctx, "user attempted login via restricted provider", "allowed", source, "attempted", handler.Name(), "user", user.GetID())
			return nil, &authenticator_domain.UserNotFoundError{
				PresentedUsername: session.GetPresentedUser(),
			}
		}
	}

	// Confirm the user is authorized...
	err := session.IsAuthorized(user)
	if err != nil {
		if errors.Is(err, authenticator_domain.ErrInvalidPassword) {
			return nil, &authenticator_domain.InvalidPasswordError{
				Username: user.GetUsername(),
				UserID:   user.GetID(),
			}
		}
		return nil, fmt.Errorf("failed to validate: %w", err)
	}

	// User is authorized, create the token.
	// This can be a cookie, a JWT, or any other method supported by the token manager.
	token, err := a.tm.CreateAuthToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth token: %w", err)
	}

	a.log.InfoContext(ctx, "user logged in", "lookup", lookup, "id", id, "handler", handler.Name(), "user", user.GetID())

	return token, nil
}

func (a *authorizers) setBaseCookies(w http.ResponseWriter, token *authenticator_domain.Token) error {
	sessionExpiresAtCookie := http.Cookie{Name: "ls_expires_at", Value: fmt.Sprintf("%d", token.ExpiresAt.Unix()), Expires: time.Unix(token.ExpiresAt.Unix(), 0), HttpOnly: false, Secure: true, Path: "/"}

	if token.OAuthProvider == "" {
		oauthProviderCookie := http.Cookie{Name: "ls_oauth_provider", Value: token.OAuthProvider, Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"}
		http.SetCookie(w, &oauthProviderCookie)
	} else {
		// For OAuth to auto-login after expiration, we need to
		// set a few things..
		oauthProviderCookie := http.Cookie{Name: "ls_oauth_provider", Value: token.OAuthProvider, Expires: time.Now().UTC().AddDate(10, 0, 0), HttpOnly: false, Secure: true, Path: "/"}
		oauthHintCookie := http.Cookie{Name: "ls_oauth_hint", Value: token.OAuthHint, Expires: time.Now().UTC().AddDate(10, 0, 0), HttpOnly: true, Secure: true, Path: "/"}
		http.SetCookie(w, &oauthProviderCookie)
		http.SetCookie(w, &oauthHintCookie)
	}

	http.SetCookie(w, &sessionExpiresAtCookie)
	return nil
}
