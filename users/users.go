package users

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
)

type LocksmithUserInterface interface {
	ValidatePassword(string) (bool, error)
	ValidateSessionToken(token string, db database.DatabaseAccessor) bool
	GeneratePasswordSession() (authentication.PasswordSession, error)
	SavePasswordSession(authentication.PasswordSession, database.DatabaseAccessor) error

	GetLastLoginDate() time.Time

	// Read from Database
	ReadFromMap(*LocksmithUserInterface, map[string]interface{})
	ToMap() map[string]interface{}

	// Convert to "public" interface
	// Slimmed down version of this interface
	// with less sensitive information
	ToPublic() (PublicLocksmithUserInterface, error)

	GetRole() (roles.Role, error)

	// Magic Auth stuff
	CleanupOldMagicTokens(database.DatabaseAccessor)
	SetMagicPermissions([]string) LocksmithUserInterface
	SetMagic() LocksmithUser // Denotes the user as a "magic only" user
	CreateMagicAuthenticationCode(database.DatabaseAccessor, magic.MagicAuthenticationVariables) (string, error)
	GetMagicPermissions() []string
	IsMagic() bool

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
	ID               string                          `bson:"id"`
	Username         string                          `json:"username" bson:"username"`
	Email            string                          `json:"email" bson:"email"`
	PasswordInfo     authentication.PasswordInfo     `json:"-" bson:"password"`
	WebAuthnSessions []webauthn.SessionData          `json:"-" bson:"websessions"`
	PasswordSessions authentication.PasswordSessions `json:"-" bson:"sessions"`
	Magics           magic.MagicAuthentications      `json:"-" bson:"magic"`
	Role             string                          `json:"role" bson:"role"`
	MagicPermissions []string                        `json:"-" bson:"-"`
	ImMagic          bool                            `json:"-" bson:"-"`
	ImRegular        bool                            `json:"-" bson:"-"`
	LastLogin        time.Time                       `json:"-" bson:"-"`
}

func (u LocksmithUser) GetLastLoginDate() time.Time {
	return u.LastLogin.UTC()
}

func (u LocksmithUser) GetMagics() []magic.MagicAuthentication {
	return u.Magics
}

func (u LocksmithUser) GetRole() (roles.Role, error) {
	role, err := roles.GetRole(u.Role)

	if err != nil {
		return roles.Role{}, err
	}

	if len(u.MagicPermissions) > 0 {
		role.Permissions = u.MagicPermissions
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

func (u LocksmithUser) ToMap() map[string]interface{} {
	out := make(map[string]interface{})

	out["id"] = u.ID
	out["username"] = u.Username
	out["email"] = u.Email
	out["password"] = u.PasswordInfo.ToMap()
	out["websessions"] = map[string]interface{}{} // TODO
	out["sessions"] = u.PasswordSessions.ToMap()
	out["role"] = u.Role
	out["magic"] = u.Magics.ToMap()

	if u.GetLastLoginDate().IsZero() {
		out["last_login"] = time.Now().UTC().Unix()
	} else {
		out["last_login"] = u.GetLastLoginDate().Unix()
	}

	return out
}

func (u LocksmithUser) ReadFromMap(writeTo *LocksmithUserInterface, user map[string]interface{}) {
	sessions := []authentication.PasswordSession{}

	if user["sessions"] != nil {
		rawSessions := user["sessions"]
		if mappedSessions, ok := rawSessions.([]map[string]interface{}); ok {
			sessions = authentication.PasswordSessions{}.FromMap(mappedSessions)
		} else {

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
	}

	var passinfo authentication.PasswordInfo
	switch user["password"].(type) {
	case authentication.PasswordInfo:
		passinfo = user["password"].(authentication.PasswordInfo)
	case map[string]interface{}:
		passinfo = authentication.PasswordInfoFromMap(user["password"].(map[string]interface{}))
	}

	var magics []magic.MagicAuthentication
	if magicValue, ok := user["magic"]; ok {
		switch magicValue.(type) {
		case []magic.MagicAuthentication:
			magics = magicValue.([]magic.MagicAuthentication)
		case []interface{}:
			magics = magic.MagicsFromMap(magicValue.([]interface{}))
		}
	}

	var loginTime time.Time
	if user["last_login"] != nil {
		loginTime = time.Unix(user["last_login"].(int64), 0)
	} else {
		// If they've yet to login, set to now.
		loginTime = time.Now().UTC()
	}

	*writeTo = LocksmithUser{
		ID:               user["id"].(string),
		Username:         user["username"].(string),
		Email:            user["email"].(string),
		Role:             user["role"].(string),
		PasswordInfo:     passinfo,
		PasswordSessions: sessions,
		Magics:           magics,
		LastLogin:        loginTime,
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

	_, err = db.UpdateOne("users", map[string]interface{}{
		"username": u.Username,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: {
			"last_login": time.Now().UTC().Unix(),
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
			// Helps prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(session.Token), []byte(hashedToken)) == 1 {
				found = true
			}
		}
	}

	updateMap := map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: {},
	}

	if len(nonexpiredTokens) != len(u.PasswordSessions) {
		// Create a new slice of type []interface{}
		interfaces := make([]interface{}, len(nonexpiredTokens))

		// Convert each PasswordSession to an interface{}
		for i, session := range nonexpiredTokens {
			interfaces[i] = session
		}

		updateMap[database.SET]["sessions"] = interfaces
	}

	if len(updateMap[database.SET]) > 0 {
		db.UpdateOne("users", map[string]interface{}{
			"id": u.ID,
		}, updateMap)
	}

	return found
}

func (u LocksmithUser) SetMagicPermissions(permissions []string) LocksmithUserInterface {
	u.MagicPermissions = permissions
	return u
}

func (u LocksmithUser) SetMagic() LocksmithUser {
	u.ImMagic = true
	return u
}

func (u LocksmithUser) IsMagic() bool {
	return u.ImMagic && !u.ImRegular
}

func (u LocksmithUser) CreateMagicAuthenticationCode(db database.DatabaseAccessor, vars magic.MagicAuthenticationVariables) (string, error) {
	mac, identifier, err := magic.CreateMagicAuthentication(vars)

	if err != nil {
		return "", err
	}

	db.UpdateOne("users", map[string]interface{}{
		"id": u.ID,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"magic": mac.ToMap(),
		},
	})

	return identifier, nil
}

func (u LocksmithUser) GetMagicPermissions() []string {
	return u.MagicPermissions
}

func (u LocksmithUser) CleanupOldMagicTokens(db database.DatabaseAccessor) {
	if magic.MagicSigningPackage == nil {
		return
	}
	activeMagics := make(chan magic.MagicAuthentications)
	go magic.FilterActive(activeMagics, u.Magics)
	keep := <-activeMagics

	if len(keep) != len(u.Magics) {
		db.UpdateOne("users", map[string]interface{}{
			"id": u.ID,
		}, map[database.DatabaseUpdateActions]map[string]interface{}{
			database.SET: {
				"magic": keep.ToMap(),
			},
		})
	}
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
