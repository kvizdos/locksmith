package invitations

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
)

type Invitation struct {
	Code         string `json:"code,omitempty"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	AttachUserID string `json:"userid"` // Attach THIS user ID once they register
	InvitedBy    string `json:"inviter"`
	SentAt       int64  `json:"sentAt"` // time that invite was sent
}

func (i Invitation) Expire(db database.DatabaseAccessor) {
	db.DeleteOne("invites", map[string]interface{}{
		"code": i.Code,
	})
}

func (i Invitation) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    i.Code,
		"email":   i.Email,
		"sentAt":  i.SentAt,
		"role":    i.Role,
		"inviter": i.InvitedBy,
		"userid":  i.AttachUserID,
	}
}

func InvitationFromMap(inp any) Invitation {
	input := inp.(map[string]interface{})

	return Invitation{
		Code:         input["code"].(string),
		Email:        input["email"].(string),
		SentAt:       input["sentAt"].(int64),
		Role:         input["role"].(string),
		InvitedBy:    input["inviter"].(string),
		AttachUserID: input["userid"].(string),
	}
}

func ListInvites(db database.DatabaseAccessor) []Invitation {
	rawInvite, found := db.Find("invites", map[string]interface{}{})

	if !found {
		return []Invitation{}
	}

	out := make([]Invitation, len(rawInvite))
	for i, raw := range rawInvite {
		inv := InvitationFromMap(raw)
		inv.Code = ""
		out[i] = inv
	}

	return out
}

// InviteUser() is a handler that allows applications to directly
// import users (think through migration, importing, etc). It returns
// a string and an error, where the string is the "invite code" used
// to register an account.
// InvitedBy is the UID of the user who invited this email.
// Returns [inviteCode, attachUserID, error]
func InviteUser(db database.DatabaseAccessor, email string, role string, invitedBy string) (string, string, error) {
	if !roles.RoleExists(role) {
		return "", "", fmt.Errorf("invalid role")
	}

	if invitedBy == "" {
		return "", "", fmt.Errorf("invitedBy is required")
	}

	emailPattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	isValidemail, _ := regexp.MatchString(emailPattern, email)

	if !isValidemail {
		return "", "", fmt.Errorf("invalid email address")
	}

	_, alreadyRegistered := db.FindOne("users", map[string]interface{}{
		"email": email,
	})

	if alreadyRegistered {
		return "", "", fmt.Errorf("email already registered")
	}

	_, alreadyInvited := db.FindOne("invites", map[string]interface{}{
		"email": email,
	})

	if alreadyInvited {
		return "", "", fmt.Errorf("email already invited")
	}

	inviteCode, err := authentication.GenerateRandomString(96)

	hasher := sha256.New()
	hasher.Write([]byte(inviteCode))
	hashedCode := hasher.Sum(nil)

	if err != nil {
		return "", "", fmt.Errorf("error generating secure invite code: %s", err.Error())
	}

	attachUserID := uuid.New().String()

	newInvite := Invitation{
		Code:         fmt.Sprintf("%x", hashedCode),
		Email:        email,
		SentAt:       time.Now().Unix(),
		InvitedBy:    invitedBy,
		Role:         role,
		AttachUserID: attachUserID,
	}

	_, err = db.InsertOne("invites", newInvite.ToMap())

	if err != nil {
		return "", "", fmt.Errorf("unable to insert invite into database: %s", err.Error())
	}

	return inviteCode, attachUserID, nil
}

// If a user needs reinviting, use this function.
// It will return:
// (newInviteCode, error)
func ReinviteUser(db database.DatabaseAccessor, forUserID string, authUserID string, newEmail ...string) (string, error) {
	if authUserID == "" {
		return "", fmt.Errorf("authUserID required")
	}

	// If a new email is present,
	// confirm it hasn't been taken
	// already by a registered user
	// or invite.
	if len(newEmail) > 0 {
		_, alreadyRegistered := db.FindOne("users", map[string]interface{}{
			"email": newEmail[0],
		})

		if alreadyRegistered {
			return "", fmt.Errorf("email already registered")
		}

		_, alreadyInvited := db.FindOne("invites", map[string]interface{}{
			"email": newEmail[0],
		})

		if alreadyInvited {
			return "", fmt.Errorf("email already invited")
		}
	}

	_, inviteFound := db.FindOne("invites", map[string]interface{}{
		"userid": forUserID,
	})

	if !inviteFound {
		return "", fmt.Errorf("could not find invite")
	}

	inviteCode, err := authentication.GenerateRandomString(96)

	hasher := sha256.New()
	hasher.Write([]byte(inviteCode))
	hashedCode := hasher.Sum(nil)

	if err != nil {
		return "", fmt.Errorf("error generating secure invite code: %s", err.Error())
	}

	updateBody := map[string]interface{}{
		"code":    fmt.Sprintf("%x", hashedCode),
		"sentAt":  time.Now().UTC().Unix(),
		"inviter": authUserID,
	}

	if len(newEmail) > 0 {
		updateBody["email"] = newEmail[0]
	}

	_, err = db.UpdateOne("invites", map[string]interface{}{
		"userid": forUserID,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: updateBody,
	})

	if err != nil {
		return "", fmt.Errorf("received error while updating invite: %s", err)
	}

	return inviteCode, nil
}

func GetInviteFromCode(db database.DatabaseAccessor, code string) (Invitation, error) {
	if len(code) != 96 {
		return Invitation{}, fmt.Errorf("invalid token length")
	}

	hasher := sha256.New()
	hasher.Write([]byte(code))
	hashedCode := hasher.Sum(nil)

	rawInvite, inviteFound := db.FindOne("invites", map[string]interface{}{
		"code": fmt.Sprintf("%x", hashedCode),
	})

	if !inviteFound {
		return Invitation{}, fmt.Errorf("could not find token")
	}

	inv := rawInvite.(map[string]interface{})

	invite := Invitation{
		Code:         inv["code"].(string),
		Email:        inv["email"].(string),
		Role:         inv["role"].(string),
		InvitedBy:    inv["inviter"].(string),
		SentAt:       inv["sentAt"].(int64),
		AttachUserID: inv["userid"].(string),
	}

	return invite, nil
}
