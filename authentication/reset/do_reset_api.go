package reset

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/hibp"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/logger"
	"github.com/kvizdos/locksmith/users"
)

type ResetPasswordAPIHandler struct {
	MinimumLengthRequirement int
	HIBP                     hibp.HIBPSettings
}

type resetPasswordRequest struct {
	Password string `json:"password"`
	PwnOK    bool   `json:"pwnok,omitempty"`
}

func (r resetPasswordRequest) HasRequiredFields() bool {
	return r.Password != ""
}

type resetResponse struct {
	Error     string `json:"error,omitempty"`
	PwnStatus bool   `json:"pwned,omitempty"`
}

func (r resetResponse) Marshal() []byte {
	js, _ := json.Marshal(r)
	return js
}

func (r *resetResponse) Unmarshal(err []byte) {
	json.Unmarshal(err, r)
}

func (h ResetPasswordAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authUser := r.Context().Value("authUser").(users.LocksmithUser)
	db := r.Context().Value("database").(database.DatabaseAccessor)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resetResponse{
			Error: "could not unmarshal",
		}.Marshal())
		return
	}

	var resetReq resetPasswordRequest
	err = json.Unmarshal(body, &resetReq)

	if err != nil || (err == nil && !resetReq.HasRequiredFields()) {
		logger.LOGGER.Log(logger.BAD_REQUEST, logger.GetIPFromRequest(*r), r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resetResponse{
			Error: "missing fields",
		}.Marshal())
		return
	}

	// Start HIBP Check
	hibpIsPwnedChan := make(chan bool)
	if h.HIBP.Enabled && !(h.HIBP.Enforcement == hibp.LOOSE && resetReq.PwnOK) {
		httpClient := &http.Client{}
		if h.HIBP.HTTPClient != nil {
			httpClient = h.HIBP.HTTPClient
		}

		go hibp.CheckPassword(h.HIBP.AppName, resetReq.Password, hibpIsPwnedChan, httpClient)
	}

	// Confirm Password Length Requirements
	if h.MinimumLengthRequirement != 0 && h.MinimumLengthRequirement > len(resetReq.Password) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resetResponse{
			Error: "password too short",
		}.Marshal())
		return
	}

	fmt.Println(h.HIBP.Enabled)
	if h.HIBP.Enabled && !(h.HIBP.Enforcement == hibp.LOOSE && resetReq.PwnOK) {
		passwordIsPwned := <-hibpIsPwnedChan
		fmt.Println("Password Pwned", passwordIsPwned)
		if passwordIsPwned && (h.HIBP.Enforcement == hibp.STRICT || (h.HIBP.Enforcement == hibp.LOOSE && !resetReq.PwnOK)) {
			w.WriteHeader(http.StatusConflict)
			w.Write(resetResponse{
				Error:     "password pwned",
				PwnStatus: true,
			}.Marshal())
			return
		}
	}

	password, err := authentication.CompileLocksmithPassword(resetReq.Password)

	if err != nil {
		fmt.Println("Error compiling password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.UpdateOne("users", map[string]interface{}{
		"id": authUser.ID,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: {
			"password": password.ToMap(),
			"sessions": []interface{}{},
		},
	})

	cookie, err := r.Cookie("magic")

	if err == nil {
		magic.ExpireOld(db, authUser.ID, cookie.Value)
		c := &http.Cookie{
			Name:    "magic",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),

			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}

	if err != nil {
		fmt.Println("Password Reset Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
