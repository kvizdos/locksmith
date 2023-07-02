package users

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
)

type LocksmithUserInterface interface {
	ValidatePassword(string) (bool, error)
	GeneratePasswordSession() (authentication.PasswordSession, error)
	SavePasswordSession(authentication.PasswordSession, database.DatabaseAccessor) error

	// Read from Database
	ReadFromMap(*LocksmithUserInterface, map[string]interface{})

	// Convert to "public" interface
	// Slimmed down version of this interface
	// with less sensitive information
	ToPublic() (PublicLocksmithUserInterface, error)

	GetRole() (roles.Role, error)

	WebAuthnID() []byte
	WebAuthnDisplayName() string
	WebAuthnName() string
	WebAuthnCredentials() []webauthn.Credential

	// WebAuthnIcon is a deprecated option.
	// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
	WebAuthnIcon() string
	AddNewWebAuthnCredential(*webauthn.Credential)

	// Getters
	GetUsername() string
	GetEmail() string
	GetID() string
	GetPasswordInfo() authentication.PasswordInfo
	GetWebAuthnSessions() []webauthn.SessionData
	GetPasswordSessions() []authentication.PasswordSession
}

// This interface is used when user structures are
// sent to the frontend
type PublicLocksmithUserInterface interface {
	FromRegular(LocksmithUserInterface) (PublicLocksmithUserInterface, error)
}

// Only used to show user data to an endpoint
// This should hide any sensitive data like
// password session info, etc
type PublicLocksmithUser struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	ActiveSessionCount int    `json:"sessions"`
	LastActive         int64  `json:"lastActive"`
	Role               string `json:"role"`
}

// Convert a LocksmithUser{} into
// the public equivalent.
func (u PublicLocksmithUser) FromRegular(user LocksmithUserInterface) (PublicLocksmithUserInterface, error) {
	publicUser := PublicLocksmithUser{}

	role, err := user.GetRole()

	if err != nil {
		return PublicLocksmithUser{}, err
	}

	publicUser.Username = user.GetUsername()
	publicUser.Email = user.GetEmail()
	publicUser.ActiveSessionCount = len(user.GetPasswordSessions())
	publicUser.ID = user.GetID()
	publicUser.LastActive = -1
	publicUser.Role = role.Name

	return publicUser, nil
}

type LocksmithUser struct {
	ID               string                           `bson:"id"`
	Username         string                           `json:"username" bson:"username"`
	Email            string                           `json:"email" bson:"email"`
	PasswordInfo     authentication.PasswordInfo      `json:"-" bson:"password"`
	WebAuthnSessions []webauthn.SessionData           `json:"-" bson:"websessions"`
	PasswordSessions []authentication.PasswordSession `json:"-" bson:"sessions"`
	Role             string                           `json:"role" bson:"role"`
}

func (u LocksmithUser) GetRole() (roles.Role, error) {
	role, err := roles.GetRole(u.Role)

	if err != nil {
		return roles.Role{}, err
	}

	return role, nil
}

func (u LocksmithUser) ToPublic() (PublicLocksmithUserInterface, error) {
	publicUser, err := PublicLocksmithUser{}.FromRegular(u)

	return publicUser, err
}

func (u LocksmithUser) GetID() string {
	return u.ID
}

func (u LocksmithUser) GetEmail() string {
	return u.Email
}

func (u LocksmithUser) GetPasswordSessions() []authentication.PasswordSession {
	return u.PasswordSessions
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
		Email:            user["email"].(string),
		Role:             user["role"].(string),
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
	hasher := sha256.New()
	hasher.Write([]byte(session.Token))
	hashedCode := hasher.Sum(nil)
	hashedToken := fmt.Sprintf("%x", hashedCode)

	session.Token = hashedToken

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

	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedCode := hasher.Sum(nil)
	hashedToken := fmt.Sprintf("%x", hashedCode)

	for _, session := range u.PasswordSessions {
		// Maybe renew the token here if its getting soon to expiring..
		if !session.IsExpired() {
			nonexpiredTokens = append(nonexpiredTokens, session)
			if session.Token == hashedToken {
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
func (u LocksmithUser) GetWebAuthnSessions() []webauthn.SessionData {
	return u.WebAuthnSessions
}
func (u LocksmithUser) AddNewWebAuthnCredential(cred *webauthn.Credential) {
	fmt.Println("Adding credential", cred)
}
