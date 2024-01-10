package login

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/observability"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/users"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	XSRF     string `json:"xsrf"`
	PwnOK    bool   `json:"pwnok,omitempty"`
}

func (r loginRequest) HasRequiredFields() bool {
	return !(r.Username == "" || r.Password == "" || r.XSRF == "")
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
}

func (lh LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		logger.LOGGER.Log(logger.INVALID_METHOD, logger.GetIPFromRequest(*r), r.URL.Path, "POST", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return
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

	dbUser, usernameExists := db.FindOne("users", map[string]interface{}{
		"username": strings.ToLower(loginReq.Username),
	})

	if !usernameExists {
		logger.LOGGER.Log(logger.LOGIN_INVALID_USERNAME, logger.GetIPFromRequest(*r), loginReq.Username)
		observability.LoginFailures.WithLabelValues("invalid_username").Inc()
		w.WriteHeader(http.StatusNotFound)
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
		logger.LOGGER.Log(logger.LOGIN_FAIL_BAD_PASSWORD, loginReq.Username, logger.GetIPFromRequest(*r))
		observability.LoginFailures.WithLabelValues("invalid_password").Inc()
		w.WriteHeader(http.StatusUnauthorized)
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

	if err != nil {
		fmt.Println("Error marshalling session token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = user.SavePasswordSession(session, db)

	if err != nil {
		fmt.Println("Error saving session token to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.LOGGER.Log(logger.LOGIN, loginReq.Username, logger.GetIPFromRequest(*r))

	observability.LoginSuccess.Inc()

	cookieValue := user.GenerateCookieValueFromSession(session)

	// Expire Login XSRF cookie
	cookieXSRF := http.Cookie{Name: "login_xsrf", Value: "", Expires: time.Unix(0, 0), HttpOnly: true, Secure: true, Path: "/api/login", SameSite: http.SameSiteStrictMode}

	// Attach Session Cookie
	cookie := http.Cookie{Name: "token", Value: cookieValue, Expires: time.Unix(session.ExpiresAt, 0), HttpOnly: true, Secure: true, Path: "/"}
	http.SetCookie(w, &cookie)
	http.SetCookie(w, &cookieXSRF)
}

type LoginPageHandler struct {
	AppName string
	// Only allow users with an invite code to register
	DisablePublicRegistration bool
	Styling                   pages.LocksmithPageStyling
	EmailAsUsername           bool
	OnboardingPath            string
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
		Title           string
		Styling         pages.LocksmithPageStyling
		EmailAsUsername bool
		OnboardingPath  string
		LoginXSRF       string
	}

	data := PageData{
		Title:           lr.AppName,
		Styling:         lr.Styling,
		EmailAsUsername: lr.EmailAsUsername,
		OnboardingPath:  lr.OnboardingPath,
		LoginXSRF:       loginXSRF,
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
