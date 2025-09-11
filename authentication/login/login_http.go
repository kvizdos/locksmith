package login

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/authentication/oauth"
	captchaproviders "github.com/kvizdos/locksmith/captcha-providers"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/observability"
	"github.com/kvizdos/locksmith/pages"
	sharedmemory "github.com/kvizdos/locksmith/shared-memory"
	"github.com/kvizdos/locksmith/shared-memory/objects"
	"github.com/kvizdos/locksmith/shared-memory/providers"
	"github.com/kvizdos/locksmith/users"
)

type LockoutPolicy struct {
	CaptchaAfter int
	LockoutAfter int
	ResetAfter   time.Duration
	OnLockout    func(username string)
}

type LoginOptions struct {
	// OnboardPath string
	// InactivityLockDuration map[string]time.Duration
	LockoutPolicy   LockoutPolicy
	CaptchaProvider captchaproviders.CAPTCHAProvider
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	XSRF     string `json:"xsrf"`
	PwnOK    bool   `json:"pwnok,omitempty"`
}

func (r loginRequest) HasRequiredFields() bool {
	return !(r.Username == "" || r.Password == "")
}

type LoginHandler struct {
	HIBP hibp.HIBPSettings
	// Set this to be longer than your session
	// duration. Session durations do not
	// change the last login date, which is
	// used for comparison.
	// It is set for every new session made,
	// so once refresh is enabled, it will
	// update the last login once refreshed.
	LockInactivityAfter map[string]time.Duration
	SharedMemory        sharedmemory.MemoryProvider
	Options             LoginOptions
	LoginInfoCallback   func(method string, user map[string]any)
}

type LoginHTTPResponse struct {
	Error           string `json:"error"`
	CaptchaRequired bool   `json:"captcha"`
}

func (l LoginHTTPResponse) Marshal() []byte {
	js, _ := json.Marshal(l)
	return js
}

func formatDuration(milliseconds int) string {
	totalSeconds := milliseconds / 1000

	hours := totalSeconds / 3600
	totalSeconds %= 3600
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutes", minutes))
	}
	if hours == 0 && minutes == 0 {
		parts = append(parts, fmt.Sprintf("%d seconds", seconds))
	}

	return strings.Join(parts, " ")
}

func (lh LoginHandler) generateInvalidUsernamePasswordError(attemptsRemaining int, timeTillUnlock int64) string {
	if attemptsRemaining == 0 {
		return fmt.Sprintf("Account locked for %s.", formatDuration(int(timeTillUnlock)))
	}

	if attemptsRemaining <= 3 {
		return fmt.Sprintf("Invalid username or password. %d attempts remaining.", attemptsRemaining)
	}

	return "Invalid username or password."
}

func (lh LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	failedLoginResponse := LoginHTTPResponse{
		Error:           "",
		CaptchaRequired: false,
	}

	if r.Method != "POST" {
		// Make it more of a pain to detect this login_xsrf cookie
		// if you aren't careful paying attention.
		cookieXSRF := http.Cookie{
			Name:     "login_xsrf",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
			Path:     "/api/login",
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookieXSRF)

		observability.LoginFailures.WithLabelValues("bad_method").Inc()
		failedLoginResponse.Error = "Incorrect HTTP method"
		logger.LOGGER.Log(logger.INVALID_METHOD, logger.GetIPFromRequest(*r), r.URL.Path, "POST", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(failedLoginResponse.Marshal())
		return
	}

	start := time.Now()
	const minDuration = 3 * time.Second

	delayIfNeeded := func() {
		if elapsed := time.Since(start); elapsed < minDuration {
			time.Sleep(minDuration - elapsed)
		}
	}

	// loginXSRFCookie, err := r.Cookie("login_xsrf")

	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	observability.LoginFailures.WithLabelValues("bad_request").Inc()
	// 	return
	// }

	if r.Body == nil {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		observability.LoginFailures.WithLabelValues("bad_body").Inc()
		failedLoginResponse.Error = "Bad request sent to server."
		w.WriteHeader(http.StatusBadRequest)
		w.Write(failedLoginResponse.Marshal())
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		failedLoginResponse.Error = "Something went wrong. Please try again later."
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(failedLoginResponse.Marshal())
		return
	}

	var loginReq loginRequest
	err = json.Unmarshal(body, &loginReq)

	if err != nil {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !loginReq.HasRequiredFields() {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		failedLoginResponse.Error = "Missing required field. Please try again."
		w.WriteHeader(http.StatusBadRequest)
		w.Write(failedLoginResponse.Marshal())
		return
	}

	if lh.SharedMemory == nil {
		lh.SharedMemory = providers.NewRamSharedMemoryProvider()
	}

	// if loginReq.XSRF != loginXSRFCookie.Value {
	// 	fmt.Println("Bad XSRF!", loginReq.XSRF, loginXSRFCookie.Value)
	// 	observability.LoginFailures.WithLabelValues("bad_xsrf").Inc()
	// 	logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// sidCookie, err := r.Cookie("sid")

	// if err != nil {
	// 	fmt.Println("No SID present on login request")
	// 	observability.LoginFailures.WithLabelValues("no_sid").Inc()
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// if !xsrf.Confirm(loginReq.XSRF, sidCookie.Value) {
	// 	fmt.Println("bad xsrf used")
	// 	observability.LoginFailures.WithLabelValues("xsrf_confirmation_error").Inc()
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	db := r.Context().Value("database").(database.DatabaseAccessor)
	hibpIsPwnedChan := make(chan bool)

	if lh.HIBP.Enabled && !(lh.HIBP.Enforcement == hibp.LOOSE && loginReq.PwnOK) {
		httpClient := &http.Client{}
		if lh.HIBP.HTTPClient != nil {
			httpClient = lh.HIBP.HTTPClient
		}

		go hibp.CheckPassword(lh.HIBP.AppName, loginReq.Password, hibpIsPwnedChan, httpClient)
	}

	// Before doing any database connections,
	// check if CAPTCHA is:
	// - [ ] required
	// - [ ] present
	// - [ ] needed on next login attempt
	// Dig into User-Specific Rate Limiting
	var captchaAttempts objects.UserLoginAttempt
	if hasMemory, found := lh.SharedMemory.GetFromMemory(objects.USER_LOGIN_ATTEMPTS, logger.GetIPFromRequest(*r)); found {
		captchaAttempts = hasMemory.(objects.UserLoginAttempt)
	} else {
		captchaAttempts = objects.NewUserLoginAttempt()
	}
	if captchaAttempts.Attempts >= lh.Options.LockoutPolicy.CaptchaAfter-1 {
		failedLoginResponse.CaptchaRequired = true
	}

	// Dig into User-Specific Rate Limiting
	var loginAttempts objects.UserLoginAttempt
	if hasMemory, found := lh.SharedMemory.GetFromMemory(objects.USER_LOGIN_ATTEMPTS, loginReq.Username); found {
		loginAttempts = hasMemory.(objects.UserLoginAttempt)
	} else {
		loginAttempts = objects.NewUserLoginAttempt()
	}

	attemptsRemaining := int(math.Max(0, float64(lh.Options.LockoutPolicy.LockoutAfter-loginAttempts.Attempts)))

	timeTillLockoutReset := lh.Options.LockoutPolicy.ResetAfter.Milliseconds() - (time.Now().UnixMilli() - loginAttempts.LastAttempt)

	// Reset Counter if time has elapsed.
	if timeTillLockoutReset <= 0 {
		loginAttempts.Attempts = 0
	}

	if loginAttempts.Attempts > lh.Options.LockoutPolicy.LockoutAfter {
		if loginAttempts.Attempts == lh.Options.LockoutPolicy.LockoutAfter {
			go lh.Options.LockoutPolicy.OnLockout(loginReq.Username)
			logger.LOGGER.Log(logger.LOGIN_LOCKOUT, loginReq.Username, logger.GetIPFromRequest(*r))
		} else {
			logger.LOGGER.Log(logger.LOGIN_LOCKED, loginReq.Username, logger.GetIPFromRequest(*r))
		}

		go lh.SharedMemory.Increment(objects.USER_LOGIN_ATTEMPTS, loginReq.Username, loginAttempts)
		go lh.SharedMemory.Increment(objects.USER_LOGIN_ATTEMPTS, logger.GetIPFromRequest(*r), captchaAttempts)

		observability.LoginFailures.WithLabelValues("locked_account").Inc()
		w.WriteHeader(http.StatusLocked)

		failedLoginResponse.Error = lh.generateInvalidUsernamePasswordError(attemptsRemaining, timeTillLockoutReset)
		w.Write(failedLoginResponse.Marshal())
		return
	} else {
		// Only update the last attempt time IF its not already locked out
		loginAttempts.LastAttempt = time.Now().UTC().UnixMilli()
		go lh.SharedMemory.Increment(objects.USER_LOGIN_ATTEMPTS, loginReq.Username, loginAttempts)
		go lh.SharedMemory.Increment(objects.USER_LOGIN_ATTEMPTS, logger.GetIPFromRequest(*r), captchaAttempts)
	}

	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": strings.ToLower(loginReq.Username),
	})

	if !usernameExists {
		logger.LOGGER.Log(logger.LOGIN_INVALID_USERNAME, logger.GetIPFromRequest(*r), loginReq.Username)
		observability.LoginFailures.WithLabelValues("invalid_username").Inc()
		delayIfNeeded()
		w.WriteHeader(http.StatusUnauthorized)
		failedLoginResponse.Error = lh.generateInvalidUsernamePasswordError(attemptsRemaining, timeTillLockoutReset)
		w.Write(failedLoginResponse.Marshal())
		return
	}

	var tmpUser users.LocksmithUserInterface
	users.LocksmithUser{}.ReadFromMap(&tmpUser, dbUser.(map[string]interface{}))
	user := tmpUser.(users.LocksmithUser)

	passwordValidated, err := user.ValidatePassword(loginReq.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !passwordValidated {
		go lh.SharedMemory.Increment(objects.USER_LOGIN_ATTEMPTS, user.GetID(), loginAttempts)
		go lh.SharedMemory.Increment(objects.USER_LOGIN_ATTEMPTS, logger.GetIPFromRequest(*r), captchaAttempts)

		logger.LOGGER.Log(logger.LOGIN_FAIL_BAD_PASSWORD, loginReq.Username, logger.GetIPFromRequest(*r))
		observability.LoginFailures.WithLabelValues("invalid_password").Inc()
		failedLoginResponse.Error = lh.generateInvalidUsernamePasswordError(attemptsRemaining, timeTillLockoutReset)
		delayIfNeeded()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(failedLoginResponse.Marshal())
		return
	}

	// Confirm user is not locked from inactivity
	var lockAccountsAfter time.Duration
	role, _ := user.GetRole()
	if lockAfter, ok := lh.LockInactivityAfter[role.Name]; ok {
		// Use programmer-defined role lock-out period
		lockAccountsAfter = lockAfter
	} else if defaultValue, ok := lh.LockInactivityAfter["default"]; ok {
		// Use the default value if it is not found
		lockAccountsAfter = defaultValue
	} else {
		// If no Default is specified, use 100 years and throw a log message.
		fmt.Println("WARNING: No default LockInactivityAfter period set. Using 100 years.")
		lockAccountsAfter = 24 * 365 * 100 * time.Hour
	}
	if time.Now().UTC().After(user.GetLastLoginDate().Add(lockAccountsAfter)) {
		logger.LOGGER.Log(logger.LOGIN_LOCKED, loginReq.Username, logger.GetIPFromRequest(*r))
		observability.LoginFailures.WithLabelValues("locked_account").Inc()
		w.WriteHeader(http.StatusLocked)
		failedLoginResponse.Error = "Account Locked. Please contact support."
		w.Write(failedLoginResponse.Marshal())
		return
	}

	// before a session is made,
	// confirm the HIBP status.
	if lh.HIBP.Enabled && !(lh.HIBP.Enforcement == hibp.LOOSE && loginReq.PwnOK) {
		passwordIsPwned := <-hibpIsPwnedChan
		if passwordIsPwned && (lh.HIBP.Enforcement == hibp.STRICT || (lh.HIBP.Enforcement == hibp.LOOSE && !loginReq.PwnOK)) {
			http.Redirect(w, r, "/reset-password?hibp=true", http.StatusTemporaryRedirect)
			return
		}
	}

	session, err := user.GeneratePasswordSession()

	if err != nil {
		fmt.Println("Error generating session token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = user.SavePasswordSession(session, db)

	if err != nil {
		fmt.Println("Error saving session token to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Reset Attempts
	lh.SharedMemory.DeleteMemory(objects.USER_LOGIN_ATTEMPTS, loginReq.Username)
	lh.SharedMemory.DeleteMemory(objects.USER_LOGIN_ATTEMPTS, logger.GetIPFromRequest(*r))

	logger.LOGGER.Log(logger.LOGIN, loginReq.Username, logger.GetIPFromRequest(*r))

	observability.LoginSuccess.Inc()

	cookieValue := user.GenerateCookieValueFromSession(session)

	// Expire Login XSRF cookie
	cookieXSRF := http.Cookie{Name: "login_xsrf", Value: "", Expires: time.Unix(0, 0), HttpOnly: true, Secure: true, Path: "/api/login", SameSite: http.SameSiteStrictMode}

	// Attach Session Cookie
	cookie := http.Cookie{Name: "token", Value: cookieValue, Expires: time.Unix(session.ExpiresAt, 0), HttpOnly: true, Secure: true, Path: "/"}

	sessionExpiresAtCookie := http.Cookie{Name: "ls_expires_at", Value: fmt.Sprintf("%d", session.ExpiresAt), Expires: time.Unix(session.ExpiresAt, 0), HttpOnly: false, Secure: true, Path: "/"}

	oauthprovidercookie := http.Cookie{Name: "ls_oauth_provider", Value: "", Expires: time.Unix(0, 0), HttpOnly: false, Secure: true, Path: "/"}

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &sessionExpiresAtCookie)
	http.SetCookie(w, &oauthprovidercookie)
	http.SetCookie(w, &cookieXSRF)

	if lh.LoginInfoCallback != nil {
		lh.LoginInfoCallback("password", dbUser.(map[string]any))
	}
}

type LoginPageHandler struct {
	AppName string
	// Only allow users with an invite code to register
	DisablePublicRegistration bool
	Styling                   pages.LocksmithPageStyling
	EmailAsUsername           bool
	OnboardingPath            string
	CaptchaProvider           captchaproviders.CAPTCHAProvider
	OAuthProviders            oauth.OAuthProviders
}

func (lr LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	loginXSRF, ok := r.Context().Value("login_xsrf").(string)

	if !ok || (ok && loginXSRF == "") {
		w.Write([]byte("Login Handler must be wrapped in LoginPageMiddleware"))
		return
	}

	tmpl, err := template.New("login.html").Parse(string(pages.LoginPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	type PageData struct {
		Title                     string
		Styling                   pages.LocksmithPageStyling
		EmailAsUsername           bool
		OnboardingPath            string
		LoginXSRF                 string
		OAuthProviders            string
		CaptchaProvider           captchaproviders.CAPTCHAProvider
		DisablePublicRegistration bool
	}

	providers := ""

	if lr.OAuthProviders != nil {
		js, _ := json.Marshal(lr.OAuthProviders.GetNames())
		providers = string(js)
	}

	data := PageData{
		Title:                     lr.AppName,
		Styling:                   lr.Styling,
		EmailAsUsername:           lr.EmailAsUsername,
		OnboardingPath:            lr.OnboardingPath,
		LoginXSRF:                 loginXSRF,
		CaptchaProvider:           lr.CaptchaProvider,
		OAuthProviders:            providers,
		DisablePublicRegistration: lr.DisablePublicRegistration,
	}

	if data.Styling.SubmitColor == "" {
		data.Styling.SubmitColor = "#476ade"
	}

	if data.Styling.StartGradient == "" {
		data.Styling.StartGradient = "#476ade"
	}

	if data.Styling.EndGradient == "" {
		data.Styling.EndGradient = "#2744a3"
	}

	if data.Title == "" {
		data.Title = "Locksmith"
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
