package administration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"text/template"
	"time"

	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/pages"
	"github.com/kvizdos/locksmith/users"
)

type AdministrationListUsersHandler struct {
	UserInterface users.LocksmithUserInterface
}

func (h AdministrationListUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.UserInterface == nil {
		h.UserInterface = users.LocksmithUser{}
	}

	queryValues := r.URL.Query()
	roles, roleExists := queryValues["role"]

	parsedRoles := []string{}

	if roleExists {
		for _, role := range roles {
			parsedRoles = append(parsedRoles, role)
		}
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	lsu := reflect.Zero(reflect.TypeOf(h.UserInterface)).Interface().(users.LocksmithUserInterface)
	users, err := ListUsers(db, ListUsersOptions{
		CustomInterface: lsu,
		GetRoles:        parsedRoles,
	})

	if err != nil {
		fmt.Println("failed to serve listing:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsond, err := json.Marshal(users)

	if err != nil {
		fmt.Println("failed to marshal user json:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsond)
}

type AdministrationLockStatusAPI struct {
	UserInterface       users.LocksmithUserInterface
	LockInactivityAfter time.Duration
}

func (h AdministrationLockStatusAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.UserInterface == nil {
		h.UserInterface = users.LocksmithUser{}
	}

	authUser, _ := r.Context().Value("authUser").(users.LocksmithUser)
	db := r.Context().Value("database").(database.DatabaseAccessor)

	requestingUserID := r.URL.Query().Get("id")

	if requestingUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		rawUser, found := db.FindOne("users", map[string]interface{}{
			"id": requestingUserID,
		})

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var tmpUser users.LocksmithUserInterface
		users.LocksmithUser{}.ReadFromMap(&tmpUser, rawUser.(map[string]interface{}))
		user := tmpUser.(users.LocksmithUser)

		type LastLoginInfo struct {
			Locked    bool  `json:"locked"`
			LastLogin int64 `json:"last_login"`
			LocksAt   int64 `json:"locks_at"`
		}

		// Confirm user is not locked from inactivity
		var lockAccountsAfter time.Duration
		if h.LockInactivityAfter > 0 {
			lockAccountsAfter = h.LockInactivityAfter
		} else {
			// If no Default is specified, use 100 years and throw a log message.
			fmt.Println("WARNING: No default LockInactivityAfter period set. Using 100 years.")
			lockAccountsAfter = 24 * 365 * 100 * time.Hour
		}

		last := LastLoginInfo{
			Locked:    time.Now().UTC().After(user.GetLastLoginDate().Add(lockAccountsAfter)),
			LastLogin: user.GetLastLoginDate().Unix(),
			LocksAt:   user.GetLastLoginDate().Add(lockAccountsAfter).Unix(),
		}

		js, _ := json.Marshal(last)
		w.Write(js)
	} else if r.Method == http.MethodPost {
		if !authUser.HasPermission("users.lock.manage") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		setTo := r.URL.Query().Get("status")

		if setTo != "0" && setTo != "1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var setLoginTimeTo time.Time

		switch setTo {
		case "0": // unlocked
			setLoginTimeTo = time.Now().UTC()
		case "1": // lock
			setLoginTimeTo = time.Date(1970, 1, 1, 1, 1, 1, 0, time.UTC)
		}

		_, err := db.UpdateOne("users", map[string]interface{}{
			"id": requestingUserID,
		}, map[database.DatabaseUpdateActions]map[string]interface{}{
			database.SET: {
				"last_login": setLoginTimeTo.Unix(),
				"sessions":   []interface{}{},
			},
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type AdministrationDeleteUsersHandler struct {
	UserInterface users.LocksmithUserInterface
}

func (h AdministrationDeleteUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type deleteRequest struct {
		Username string `json:"username"`
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("authUser").(users.LocksmithUser)

	if !ok {
		fmt.Println("Delete users endpoint is required to be wrapped in SecureEndpointHTTPMiddleware()")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// handle the error
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var delReq deleteRequest
	err = json.Unmarshal(body, &delReq)

	if err != nil || (err == nil && delReq.Username == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if delReq.Username != user.GetUsername() {
		if !user.HasPermission("user.delete.other") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else if !user.HasPermission("user.delete.self") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db := r.Context().Value("database").(database.DatabaseAccessor)

	deleted, err := DeleteUser(db, delReq.Username)

	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type ServeAdminPage struct{}

func (p ServeAdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("admin.html").Parse(string(pages.AdminPageHTML))

	if err != nil {
		fmt.Println(err)
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
