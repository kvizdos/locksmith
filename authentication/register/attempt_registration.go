package register

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/administration/invitations"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/events"
	"github.com/kvizdos/locksmith/authentication/register/register_domain"
	"github.com/kvizdos/locksmith/authentication/textvalidation"
	"github.com/kvizdos/locksmith/authentication/tokens"
	"github.com/kvizdos/locksmith/authentication/verificationcodes"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/users"
)

// ServeRegisterAPI is the public HTTP entrypoint for registration, mirroring
// authenticator.(*authorizers).ServeLoginAPI: it detects which registration
// method should handle the request, dispatches to it, and turns the result
// into an HTTP response (including optional auto-login).
func (r *registrar) ServeRegisterAPI(w http.ResponseWriter, req *http.Request) {
	if r.defaultRoleName == "" {
		r.writeRegistrationError(w, register_domain.ErrRegistrationRoleMissing)
		return
	}
	if !roles.RoleExists(r.defaultRoleName) {
		r.writeRegistrationError(w, register_domain.ErrRegistrationRoleInvalid)
		return
	}

	ctx := r.registrationEventContext(req)
	req = req.WithContext(ctx)

	// First, detect which method this session will utilize. Each method
	// runs its own heuristic to determine if it can handle this request.
	handler, err := r.getHandler(req)
	if err != nil {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*req), req.URL.Path)
		r.writeRegistrationError(w, register_domain.ErrRegistrationMissingFields)
		return
	}

	result, err := r.attemptRegistration(req.Context(), handler, req)
	if err != nil {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*req), req.URL.Path)
		r.writeRegistrationError(w, err)
		return
	}

	token, err := r.tm.CreateAuthToken(result.User)
	if err != nil {
		r.log.ErrorContext(req.Context(), "failed to create registration auth token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tokens.SetBaseCookies(w, token)
	if err := r.tm.PassToClient(w, req, token); err != nil {
		r.log.ErrorContext(req.Context(), "failed to pass registration auth token to client", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// attemptRegistration loads the normalized register_domain.Request from the
// chosen method's session and runs it through the shared registration
// business rules. This mirrors authenticator.(*authorizers).attemptLogin.
func (r *registrar) attemptRegistration(ctx context.Context, handler register_domain.Handler, req *http.Request) (*register_domain.Result, error) {
	session := handler.Session(r.db)

	if err := session.LoadRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", register_domain.ErrFailedToParse, err)
	}

	return r.register(ctx, session.RegistrationRequest())
}

// register contains the shared registration business rules: validation,
// invite handling, duplicate detection, user creation, auth-link creation,
// and email verification. It is the registration equivalent of the core of
// authenticator.(*authorizers).attemptLogin, deliberately separated from
// HTTP/session concerns so it can be exercised directly in tests.
func (r *registrar) register(ctx context.Context, req register_domain.Request) (result *register_domain.Result, err error) {
	isHinted := req.Hint != nil

	if r.emailAsUsername && !isHinted {
		req.Email = req.Username
	}
	if isHinted {
		req.Username = req.Hint.Email
		req.Email = req.Hint.Email
	}

	method := "password"
	provider := "password"
	background := false
	selectBy := ""
	if isHinted {
		method = "hint"
		provider = req.Hint.ProviderName
		background = true
		selectBy = req.Hint.SelectBy
	}

	r.publishRegistrationEvent(ctx, events.EventRegistrationRequested, events.RegistrationRequestedPayload{
		Username:       strings.ToLower(req.Username),
		Email:          strings.ToLower(req.Email),
		Method:         method,
		Provider:       provider,
		InviteProvided: req.Code != "",
		Background:     background,
		SelectBy:       selectBy,
	})
	defer func() {
		if err != nil {
			r.publishRegistrationEvent(ctx, events.EventRegistrationFailed, events.RegistrationFailedPayload{
				Username:       strings.ToLower(req.Username),
				Email:          strings.ToLower(req.Email),
				Method:         method,
				Provider:       provider,
				InviteProvided: req.Code != "",
				Background:     background,
				Reason:         registrationFailureReason(err),
				SelectBy:       selectBy,
			})
		}
	}()

	if r.defaultRoleName == "" {
		return nil, register_domain.ErrRegistrationRoleMissing
	}
	if !roles.RoleExists(r.defaultRoleName) {
		return nil, register_domain.ErrRegistrationRoleInvalid
	}
	// Public registration lockdown only permits requests carrying either a
	// valid invite code or a server-signed, already-verified registration
	// hint (req.Hint is only non-nil after successful signature/claim
	// verification upstream, inside the hint method's session). Unsigned or
	// unverified hints never reach here.
	if r.disablePublicRegistration && len(req.Code) == 0 && !isHinted {
		return nil, register_domain.ErrPublicRegistrationDisabled
	}
	if isHinted && (req.Hint.ProviderName == "" || req.Hint.Issuer == "" || req.Hint.Subject == "" || !req.Hint.Rosterable) {
		return nil, register_domain.ErrRegistrationInvalidHint
	}
	if !req.HasRequiredFields() {
		return nil, register_domain.ErrRegistrationMissingFields
	}

	var emailValidationResult textvalidation.ValidationResultEvaluator
	if r.emailValidation != nil {
		validated, verr := r.emailValidation.Validate(ctx, req.Email)
		if verr == nil {
			emailValidationResult = validated
			validated.DebugPrint(req.Email)
			didYouMean, res := validated.Result(req.ValidationOK)
			switch res {
			case textvalidation.ValidationResult_CONFIRM:
				if !req.ValidationOK {
					dym := ""
					if didYouMean != nil {
						dym = *didYouMean
					}
					return nil, &register_domain.ConfirmEmailError{DidYouMean: dym}
				}
			case textvalidation.ValidationResult_REJECT:
				return nil, register_domain.ErrRegistrationEmailBlocked
			}
		} else {
			r.log.WarnContext(ctx, "skipped email validation", "error", verr)
		}
	}

	if !isHinted && r.minimumLengthRequirement != 0 && r.minimumLengthRequirement > len(req.Password) {
		return nil, register_domain.ErrRegistrationPasswordTooShort
	}

	usernamePattern := "^[a-zA-Z0-9]+$"
	if r.emailAsUsername || isHinted {
		usernamePattern = `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	}
	validUsername, _ := regexp.MatchString(usernamePattern, req.Username)
	if !validUsername {
		return nil, register_domain.ErrRegistrationIllegalUsername
	}

	validEmail, _ := regexp.MatchString(`^[^\s@]+@[^\s@]+\.[^\s@]+$`, req.Email)
	if !validEmail {
		return nil, register_domain.ErrRegistrationInvalidEmail
	}

	useRole := r.defaultRoleName
	useID := uuid.New().String()
	var invite invitations.Invitation
	inviteUsed := false
	if len(req.Code) > 0 {
		if len(req.Code) != 96 {
			return nil, register_domain.ErrRegistrationBadInviteCode
		}
		var inviteErr error
		invite, inviteErr = invitations.GetInviteFromCode(r.db, req.Code)
		if inviteErr != nil {
			return nil, register_domain.ErrRegistrationInvalidInviteCode
		}
		if invite.Email != req.Email {
			return nil, register_domain.ErrRegistrationInvalidEmail
		}
		useRole = invite.Role
		useID = invite.AttachUserID
		inviteUsed = true
	}

	matches, _ := r.db.Find("users", map[string]interface{}{
		"$or": []map[string]interface{}{{"username": strings.ToLower(req.Username)}, {"email": strings.ToLower(req.Email)}},
	})
	if len(matches) != 0 {
		return nil, register_domain.ErrRegistrationTaken
	}

	password, err := authentication.CompileLocksmithPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("compile password: %w", err)
	}

	oauthRestriction := ""
	if isHinted {
		oauthRestriction = req.Hint.ProviderName
	}

	var lsu users.LocksmithUserInterface = users.LocksmithUser{
		ID:                    useID,
		Username:              strings.ToLower(req.Username),
		Email:                 strings.ToLower(req.Email),
		PasswordInfo:          password,
		WebAuthnSessions:      []webauthn.SessionData{},
		PasswordSessions:      []authentication.PasswordSession{},
		Role:                  useRole,
		OAuthRestrictedSource: oauthRestriction,
	}

	if r.configureCustomUser != nil {
		lsu = r.configureCustomUser(lsu.(users.LocksmithUser), r.db)
	}
	if r.requiresEmailVerification != nil && req.Hint == nil {
		// Always hand the callback a non-nil evaluator, even if no
		// EmailValidator is configured, so callbacks can safely call its
		// methods without a nil check. A zero-value EmailValidationResult
		// evaluates to a non-VALID result (since it can't confirm anything),
		// which is the conservative "we don't know, so ask the app" default.
		evaluator := emailValidationResult
		if evaluator == nil {
			evaluator = textvalidation.EmailValidationResult{}
		}
		lsu = lsu.SetRequiresEmailVerification(r.requiresEmailVerification(ctx, r.db, lsu, evaluator))
	}

	_, err = r.db.InsertOne("users", lsu.ToMap())
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	if inviteUsed {
		invite.Expire(r.db)
	}

	createdAuthLink := false
	if isHinted {
		// Auth link identity is derived exclusively from the signed hint
		// (provider, issuer, subject), never from client-supplied JSON.
		_, err = r.db.InsertOne("auth_links", map[string]interface{}{
			"provider":  req.Hint.ProviderName,
			"issuer":    req.Hint.Issuer,
			"subject":   req.Hint.Subject,
			"user_id":   lsu.GetID(),
			"linked_at": time.Now().UTC().Unix(),
		})
		if err != nil {
			return nil, fmt.Errorf("insert auth link: %w", err)
		}
		createdAuthLink = true
	}
	if r.accountVerifier != nil && lsu.RequiresEmailVerification() {
		if err := r.accountVerifier.SendVerification(ctx, lsu, verificationcodes.VerifierMethod_EMAIL, lsu.GetEmail()); err != nil {
			return nil, fmt.Errorf("send verification code: %w", err)
		}
		r.publishRegistrationEvent(ctx, events.EventEmailVerificationSent, events.EmailVerificationSentPayload{
			UserID:   lsu.GetID(),
			Username: lsu.GetUsername(),
			Email:    lsu.GetEmail(),
			Method:   "email",
			Target:   lsu.GetEmail(),
		})
	}

	result = &register_domain.Result{
		User:                      lsu,
		InviteUsed:                inviteUsed,
		RequiresEmailVerification: lsu.RequiresEmailVerification(),
		Method:                    method,
		Provider:                  provider,
		CreatedAuthLink:           createdAuthLink,
		Background:                background,
	}
	r.publishRegistrationEvent(ctx, events.EventRegistrationSucceeded, events.RegistrationSucceededPayload{
		UserID:                    lsu.GetID(),
		Username:                  lsu.GetUsername(),
		Email:                     lsu.GetEmail(),
		Method:                    method,
		Provider:                  provider,
		InviteUsed:                inviteUsed,
		RequiresEmailVerification: lsu.RequiresEmailVerification(),
		Background:                background,
		AutoLoginIssued:           false,
		SelectBy:                  selectBy,
	})

	return result, nil
}

func (r *registrar) publishRegistrationEvent(ctx context.Context, name events.EventName, payload any) {
	if r.eventBus == nil {
		return
	}
	envelope := events.EnrichEnvelope(ctx, events.Envelope{
		ID:         uuid.New().String(),
		Name:       name,
		OccurredAt: time.Now(),
		Source:     "authentication/register",
		Payload:    payload,
	})
	if err := r.eventBus.Publish(ctx, envelope); err != nil && r.log != nil {
		r.log.WarnContext(ctx, "registration event publish failed", "event", string(name), "error", err)
	}
}

func (r *registrar) registrationEventContext(req *http.Request) context.Context {
	ip := logger.GetIPFromRequest(*req)
	ctx := context.WithValue(req.Context(), "ip_address", ip)
	metadata := events.ContextMetadata{
		RequestID: firstHeader(req, "X-Request-Id", "X-Request-ID", "X-Correlation-Id"),
		TraceID:   traceIDFromRequest(req),
		Source:    "authentication/register",
		Values: map[string]string{
			"ip_address": ip,
			"user_agent": req.UserAgent(),
		},
	}
	if r.requestEventMetadata != nil {
		metadata = mergeContextMetadata(metadata, r.requestEventMetadata(req))
	}
	return events.WithContextMetadata(ctx, metadata)
}

func mergeContextMetadata(base, app events.ContextMetadata) events.ContextMetadata {
	if app.RequestID != "" {
		base.RequestID = app.RequestID
	}
	if app.TraceID != "" {
		base.TraceID = app.TraceID
	}
	if app.Source != "" {
		base.Source = app.Source
	}
	if len(app.Values) > 0 {
		if base.Values == nil {
			base.Values = map[string]string{}
		}
		for key, value := range app.Values {
			base.Values[key] = value
		}
	}
	return base
}

func firstHeader(r *http.Request, names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(r.Header.Get(name)); value != "" {
			return value
		}
	}
	return ""
}

func traceIDFromRequest(r *http.Request) string {
	if traceID := firstHeader(r, "X-Trace-Id", "X-Trace-ID"); traceID != "" {
		return traceID
	}
	traceparent := strings.TrimSpace(r.Header.Get("Traceparent"))
	parts := strings.Split(traceparent, "-")
	if len(parts) >= 2 && len(parts[1]) == 32 {
		return parts[1]
	}
	return ""
}

func registrationFailureReason(err error) string {
	var confirmErr *register_domain.ConfirmEmailError
	switch {
	case errors.Is(err, register_domain.ErrRegistrationRoleMissing):
		return "role_missing"
	case errors.Is(err, register_domain.ErrRegistrationRoleInvalid):
		return "role_invalid"
	case errors.Is(err, register_domain.ErrPublicRegistrationDisabled):
		return "public_registration_disabled"
	case errors.Is(err, register_domain.ErrRegistrationMissingFields):
		return "missing_fields"
	case errors.Is(err, register_domain.ErrRegistrationPasswordTooShort):
		return "password_too_short"
	case errors.Is(err, register_domain.ErrRegistrationIllegalUsername):
		return "illegal_username"
	case errors.Is(err, register_domain.ErrRegistrationInvalidEmail):
		return "invalid_email"
	case errors.Is(err, register_domain.ErrRegistrationBadInviteCode):
		return "bad_invite_code"
	case errors.Is(err, register_domain.ErrRegistrationInvalidInviteCode):
		return "invalid_invite_code"
	case errors.Is(err, register_domain.ErrRegistrationTaken):
		return "taken"
	case errors.As(err, &confirmErr):
		return "confirm_email"
	case errors.Is(err, register_domain.ErrRegistrationEmailBlocked):
		return "email_blocked"
	case errors.Is(err, register_domain.ErrRegistrationInvalidHint):
		return "invalid_hint"
	case errors.Is(err, register_domain.ErrFailedToParse):
		return "failed_to_parse"
	default:
		return "internal"
	}
}

func (r *registrar) writeRegistrationError(w http.ResponseWriter, err error) {
	var confirmErr *register_domain.ConfirmEmailError

	switch {
	case errors.Is(err, register_domain.ErrRegistrationRoleMissing):
		fmt.Println("Registration role name must be set!")
		w.WriteHeader(http.StatusInternalServerError)
	case errors.Is(err, register_domain.ErrRegistrationRoleInvalid):
		fmt.Println("Registration role name is invalid!")
		w.WriteHeader(http.StatusInternalServerError)
	case errors.Is(err, register_domain.ErrPublicRegistrationDisabled):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, register_domain.ErrFailedToParse):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "could not unmarshal"}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationMissingFields):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "missing fields"}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationPasswordTooShort):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "password too short"}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationIllegalUsername):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "illegal username characters"}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationInvalidEmail):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "invalid email"}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationBadInviteCode):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "bad invite code"}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationInvalidInviteCode):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "invalid code"}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationTaken):
		w.WriteHeader(http.StatusConflict)
		w.Write(registrationResponse{Error: "taken"}.Marshal())
	case errors.As(err, &confirmErr):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{ConfirmEmail: true, DidYouMean: confirmErr.DidYouMean}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationEmailBlocked):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{EmailBlocked: true}.Marshal())
	case errors.Is(err, register_domain.ErrRegistrationInvalidHint):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(registrationResponse{Error: "invalid registration hint"}.Marshal())
	default:
		fmt.Println("Error registering user:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
