package authentication

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"golang.org/x/crypto/argon2"
)

type PasswordSessions []PasswordSession

func (p PasswordSessions) ToMap() []map[string]interface{} {
	out := make([]map[string]interface{}, len(p))

	for i, pw := range p {
		out[i] = pw.ToMap()
	}

	return out
}

func (p PasswordSessions) FromMap(input []map[string]interface{}) []PasswordSession {
	out := make([]PasswordSession, len(input))

	for i, m := range input {
		out[i] = PasswordSession{}.FromMap(m)
	}

	return out
}

type PasswordSession struct {
	Token     string `json:"token" bson:"token"`
	ExpiresAt int64  `json:"expire" bson:"expire"`
}

func (p PasswordSession) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"token":  p.Token,
		"expire": p.ExpiresAt,
	}
}

func (p PasswordSession) FromMap(input map[string]interface{}) PasswordSession {
	return PasswordSession{
		Token:     input["token"].(string),
		ExpiresAt: input["expire"].(int64),
	}
}

func (p PasswordSession) Marshal() ([]byte, error) {
	jsonData, err := json.Marshal(p)

	if err != nil {
		return []byte{}, err
	}

	return jsonData, nil
}

func (p PasswordSession) IsExpired() bool {
	// 60 second period before its "invalid"
	// Maybe add an extra return value (another bool)
	// that defines whether or not the token is
	// "near" expiring
	return p.ExpiresAt-time.Now().Unix() <= 60
}

func PasswordSessionFromJson(data string) (PasswordSession, error) {
	var session PasswordSession
	err := json.Unmarshal([]byte(data), &session)
	if err != nil {
		return PasswordSession{}, err
	}
	return session, nil
}

type PasswordInfo struct {
	Password            string                `json:"password" bson:"password"`
	Salt                string                `json:"salt" bson:"salt"`
	WebAuthnCredentials []webauthn.Credential `json:"webauthn" bson:"webauthn"`
}

func (p PasswordInfo) ToMap() map[string]interface{} {
	out := make(map[string]interface{})

	out["password"] = p.Password
	out["salt"] = p.Salt
	out["webauth"] = map[string]interface{}{} // TODO

	return out
}

func PasswordInfoFromMap(passinfo map[string]interface{}) PasswordInfo {
	return PasswordInfo{
		Password:            passinfo["password"].(string),
		Salt:                passinfo["salt"].(string),
		WebAuthnCredentials: []webauthn.Credential{},
	}
}

func ValidatePassword(locksmithPassword PasswordInfo, inputPassword string) (bool, error) {
	if len(locksmithPassword.Password) == 0 {
		return false, fmt.Errorf("locksmith password not presented")
	}

	if len(locksmithPassword.Salt) == 0 {
		return false, fmt.Errorf("locksmith salt not presented")
	}

	if len(inputPassword) == 0 {
		return false, fmt.Errorf("no input password passed")
	}

	generatedPassword, err := CompileLocksmithPassword(inputPassword, locksmithPassword.Salt)

	if err != nil {
		return false, fmt.Errorf("failed to generate hashed version of password")
	}

	if generatedPassword.Password != locksmithPassword.Password {
		return false, nil
	}

	return true, nil
}

func CompileLocksmithPassword(password string, saltString ...string) (PasswordInfo, error) {
	var salt []byte

	if len(saltString) != 0 {
		var decodeErr error
		salt, decodeErr = hex.DecodeString(saltString[0])

		if decodeErr != nil {
			return PasswordInfo{}, fmt.Errorf("error parsing provided salt")
		}
	} else {
		var saltError error
		salt, saltError = GenerateRandomBytes(16)

		if saltError != nil {
			return PasswordInfo{}, fmt.Errorf("error generating salt")
		}
	}

	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	encodedToHex := hex.EncodeToString(key)
	encodedSaltToHex := hex.EncodeToString(salt)

	return PasswordInfo{
		Password: encodedToHex,
		Salt:     encodedSaltToHex,
	}, nil
}

type Token struct {
	Token    string
	Username string
}

func ParseToken(cookieValue string) (Token, error) {
	decodedCookie, err := base64.StdEncoding.DecodeString(cookieValue)

	if err != nil {
		return Token{}, err
	}

	splitValue := strings.Split(string(decodedCookie), ":")

	if len(splitValue) != 2 {
		return Token{}, fmt.Errorf("invalid token")
	}

	token := splitValue[0]
	username := splitValue[1]

	if len(token) != 64 {
		return Token{}, fmt.Errorf("invalid token length")
	}

	return Token{
		Token:    token,
		Username: username,
	}, nil
}
