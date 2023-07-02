package invitations

import (
	"fmt"
	"regexp"
	"time"

	"kv.codes/locksmith/authentication"
	"kv.codes/locksmith/database"
	"kv.codes/locksmith/roles"
)

type Invitation struct {
	Code      string `json:"code"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	InvitedBy string `json:"inviter"`
	SentAt    int64  `json:"sentAt"` // time that invite was sent
}

func (i Invitation) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    i.Code,
		"email":   i.Email,
		"sentAt":  i.SentAt,
		"role":    i.Role,
		"inviter": i.InvitedBy,
	}
}

// InviteUser() is a handler that allows applications to directly
// import users (think through migration, importing, etc). It returns
// a string and an error, where the string is the "invite code" used
// to register an account.
// InvitedBy is the UID of the user who invited this email.
func InviteUser(db database.DatabaseAccessor, email string, role string, invitedBy string) (string, error) {
	if !roles.RoleExists(role) {
		return "", fmt.Errorf("invalid role")
	}

	if invitedBy == "" {
		return "", fmt.Errorf("invitedBy is required")
	}

	emailPattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	isValidemail, _ := regexp.MatchString(emailPattern, email)

	if !isValidemail {
		return "", fmt.Errorf("invalid email address")
	}

	_, alreadyRegistered := db.FindOne("users", map[string]interface{}{
		"email": email,
	})

	if alreadyRegistered {
		return "", fmt.Errorf("email already registered")
	}

	_, alreadyInvited := db.FindOne("invites", map[string]interface{}{
		"email": email,
	})

	if alreadyInvited {
		return "", fmt.Errorf("email already invited")
	}

	inviteCode, err := authentication.GenerateRandomString(96)

	if err != nil {
		return "", fmt.Errorf("error generating secure invite code: %s", err.Error())
	}

	newInvite := Invitation{
		Code:      inviteCode,
		Email:     email,
		SentAt:    time.Now().Unix(),
		InvitedBy: invitedBy,
		Role:      role,
	}

	_, err = db.InsertOne("invites", newInvite.ToMap())

	if err != nil {
		return "", fmt.Errorf("unable to insert invite into database: %s", err.Error())
	}

	return inviteCode, nil
}

func GetInviteFromCode(db database.DatabaseAccessor, code string) (Invitation, error) {
	if len(code) != 96 {
		return Invitation{}, fmt.Errorf("invalid token length")
	}

	rawInvite, inviteFound := db.FindOne("invites", map[string]interface{}{
		"code": code,
	})

	if !inviteFound {
		return Invitation{}, fmt.Errorf("could not find token")
	}

	inv := rawInvite.(map[string]interface{})

	invite := Invitation{
		Code:      inv["code"].(string),
		Email:     inv["email"].(string),
		Role:      inv["role"].(string),
		InvitedBy: inv["inviter"].(string),
		SentAt:    inv["sentAt"].(int64),
	}

	return invite, nil
}
