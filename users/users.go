package users

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
)

type LocksmithUserInterface interface {
	GetUsername() string
	GetPasswordInfo() authentication.PasswordInfo
	ValidatePassword(string) (bool, error)
	GeneratePasswordSession() (authentication.PasswordSession, error)
	SavePasswordSession(authentication.PasswordSession, database.DatabaseAccessor) error

	// Read from Database
	ReadFromMap(*LocksmithUserInterface, map[string]interface{})

	WebAuthnID() []byte
	WebAuthnDisplayName() string
	WebAuthnName() string
	WebAuthnCredentials() []webauthn.Credential

	// WebAuthnIcon is a deprecated option.
	// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
	WebAuthnIcon() string

	AddNewWebAuthnCredential(*webauthn.Credential)

	GetSessions() []webauthn.SessionData
}

type LocksmithUser struct {
	ID               string                           `bson:"id"`
	Username         string                           `json:"username" bson:"username"`
	PasswordInfo     authentication.PasswordInfo      `json:"-" bson:"password"`
	WebAuthnSessions []webauthn.SessionData           `json:"-" bson:"websessions"`
	PasswordSessions []authentication.PasswordSession `json:"-" bson:"sessions"`
}

func (u LocksmithUser) ReadFromMap(writeTo *LocksmithUserInterface, user map[string]interface{}) {
	sessions := []authentication.PasswordSession{}

	if user["sessions"] != nil {
		for _, v := range user["sessions"].([]interface{}) {
			session, ok := v.(authentication.PasswordSession)
			if !ok {
				newSession := v.(map[string]interface{})
				sessions = append(sessions, authentication.PasswordSession{
					Token:     newSession["token"].(string),
					ExpiresAt: newSession["expire"].(int64),
				})

				continue
			}

			sessions = append(sessions, session)
		}
	}

	var passinfo authentication.PasswordInfo
	switch user["password"].(type) {
	case authentication.PasswordInfo:
		passinfo = user["password"].(authentication.PasswordInfo)
	case map[string]interface{}:
		passinfo = authentication.PasswordInfoFromMap(user["password"].(map[string]interface{}))
	}
	*writeTo = LocksmithUser{
		ID:               user["id"].(string),
		Username:         user["username"].(string),
		PasswordInfo:     passinfo,
		PasswordSessions: sessions,
	}
}

func (u LocksmithUser) GetUsername() string {
	return u.Username
}

func (u LocksmithUser) GetPasswordInfo() authentication.PasswordInfo {
	return u.PasswordInfo
}

func (u LocksmithUser) ValidatePassword(inputPassword string) (bool, error) {
	return authentication.ValidatePassword(u.PasswordInfo, inputPassword)
}

func (u LocksmithUser) GeneratePasswordSession() (authentication.PasswordSession, error) {
	token, err := authentication.GenerateRandomString(64)

	if err != nil {
		return authentication.PasswordSession{}, err
	}

	now := time.Now()
	futureDate := now.AddDate(0, 0, 30)

	timestamp := futureDate.Unix()

	session := authentication.PasswordSession{
		Token:     token,
		ExpiresAt: timestamp,
	}

	return session, nil
}

func (u LocksmithUser) GenerateCookieValueFromSession(session authentication.PasswordSession) string {
	token := fmt.Sprintf("%s:%s", session.Token, u.Username)
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func (u LocksmithUser) SavePasswordSession(session authentication.PasswordSession, db database.DatabaseAccessor) error {
	_, err := db.UpdateOne("users", map[string]interface{}{
		"username": u.Username,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": session,
		},
	})

	return err
}

func (u LocksmithUser) ValidateSessionToken(token string, db database.DatabaseAccessor) bool {
	found := false

	nonexpiredTokens := []authentication.PasswordSession{}

	for _, session := range u.PasswordSessions {
		// Maybe renew the token here if its getting soon to expiring..
		if !session.IsExpired() {
			nonexpiredTokens = append(nonexpiredTokens, session)
			if session.Token == token {
				found = true
			}
		}
	}

	if len(nonexpiredTokens) != len(u.PasswordSessions) {
		// Create a new slice of type []interface{}
		interfaces := make([]interface{}, len(nonexpiredTokens))

		// Convert each PasswordSession to an interface{}
		for i, session := range nonexpiredTokens {
			interfaces[i] = session
		}

		db.UpdateOne("users", map[string]interface{}{
			"username": u.Username,
		}, map[database.DatabaseUpdateActions]map[string]interface{}{
			database.SET: {
				"sessions": interfaces,
			},
		})
	}
	// fmt.Printf("%d - %d\n", len(nonexpiredTokens), len(u.PasswordSessions))

	return found
}

func (u LocksmithUser) WebAuthnID() []byte {
	return []byte(u.ID)
}
func (u LocksmithUser) WebAuthnDisplayName() string {
	return u.Username
}
func (u LocksmithUser) WebAuthnName() string {
	return u.Username
}
func (u LocksmithUser) WebAuthnCredentials() []webauthn.Credential {
	return u.PasswordInfo.WebAuthnCredentials
}
func (u LocksmithUser) WebAuthnIcon() string {
	// IS DEPRECATED.
	return ""
}
func (u LocksmithUser) GetSessions() []webauthn.SessionData {
	return u.WebAuthnSessions
}
func (u LocksmithUser) AddNewWebAuthnCredential(cred *webauthn.Credential) {
	fmt.Println("Adding credential", cred)
}
