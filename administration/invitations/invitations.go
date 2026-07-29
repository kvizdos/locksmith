package invitations

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
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
	db.DeleteOne("invites", map[string]any{
		"code": i.Code,
	})
}

func (i Invitation) ToMap() map[string]any {
	return map[string]any{
		"code":    i.Code,
		"email":   i.Email,
		"sentAt":  i.SentAt,
		"role":    i.Role,
		"inviter": i.InvitedBy,
		"userid":  i.AttachUserID,
	}
}

func InvitationFromMap(inp any) Invitation {
	input := inp.(map[string]any)

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
	rawInvite, found := db.Find("invites", map[string]any{})

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
	email = strings.ToLower(email)

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

	_, alreadyRegistered := db.FindOne("users", map[string]any{
		"email": email,
	})

	if alreadyRegistered {
		return "", "", fmt.Errorf("email already registered")
	}

	_, alreadyInvited := db.FindOne("invites", map[string]any{
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
		_, alreadyRegistered := db.FindOne("users", map[string]any{
			"email": newEmail[0],
		})

		if alreadyRegistered {
			return "", fmt.Errorf("email already registered")
		}

		_, alreadyInvited := db.FindOne("invites", map[string]any{
			"email": newEmail[0],
		})

		if alreadyInvited {
			return "", fmt.Errorf("email already invited")
		}
	}

	_, inviteFound := db.FindOne("invites", map[string]any{
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

	updateBody := map[string]any{
		"code":    fmt.Sprintf("%x", hashedCode),
		"sentAt":  time.Now().UTC().Unix(),
		"inviter": authUserID,
	}

	if len(newEmail) > 0 {
		updateBody["email"] = newEmail[0]
	}

	_, err = db.UpdateOne("invites", map[string]any{
		"userid": forUserID,
	}, map[database.DatabaseUpdateActions]map[string]any{
		database.SET: updateBody,
	})

	if err != nil {
		return "", fmt.Errorf("received error while updating invite: %s", err)
	}

	return inviteCode, nil
}

// GetActiveInviteByEmail looks up an active (i.e. not yet consumed) invite
// for the given email, if one exists. found=false, err=nil means no invite
// exists for that email.
//
// Note: this performs an exact-match query, so callers should normalize
// casing themselves (e.g. lowercase) to match how invites are stored.
func GetActiveInviteByEmail(db database.DatabaseAccessor, email string) (Invitation, bool, error) {
	rawInvite, found := db.FindOne("invites", map[string]any{
		"email": email,
	})

	if !found {
		return Invitation{}, false, nil
	}

	inv := rawInvite.(map[string]any)

	invite := Invitation{
		Code:         inv["code"].(string),
		Email:        inv["email"].(string),
		Role:         inv["role"].(string),
		InvitedBy:    inv["inviter"].(string),
		SentAt:       inv["sentAt"].(int64),
		AttachUserID: inv["userid"].(string),
	}

	return invite, true, nil
}

// ExpireByEmail deletes the active invite for the given email, if any. It
// mirrors Invitation.Expire but doesn't require the caller to already hold
// the invite's hashed code.
func ExpireByEmail(db database.DatabaseAccessor, email string) {
	db.DeleteOne("invites", map[string]any{
		"email": email,
	})
}

// ClaimActiveInviteOnVerifiedEmail looks for an active invite matching
// userID's now-verified email and, if found, upgrades that user's role to
// the invite's role and expires the invite.
//
// This is the deferred counterpart to registering with an explicit invite
// code: a plain (non-OAuth) registration can match a pending invite by
// email, but since a client-supplied email isn't proof of anything, the
// match isn't trusted -- and the invite isn't applied -- until the user
// proves control of that address via the normal email verification flow.
// Callers should invoke this once verification succeeds (e.g. from an
// email-verification-exchange handler).
func ClaimActiveInviteOnVerifiedEmail(db database.DatabaseAccessor, userID string, email string) error {
	invite, found, err := GetActiveInviteByEmail(db, email)
	if err != nil {
		return fmt.Errorf("lookup active invite: %w", err)
	}
	if !found {
		return nil
	}

	_, err = db.UpdateOne("users", map[string]any{
		"id": userID,
	}, map[database.DatabaseUpdateActions]map[string]any{
		database.SET: {
			"role": invite.Role,
		},
	})
	if err != nil {
		return fmt.Errorf("apply invite role: %w", err)
	}

	ExpireByEmail(db, email)
	return nil
}

func GetInviteFromCode(db database.DatabaseAccessor, code string) (Invitation, error) {
	if len(code) != 96 {
		return Invitation{}, fmt.Errorf("invalid token length")
	}

	hasher := sha256.New()
	hasher.Write([]byte(code))
	hashedCode := hasher.Sum(nil)

	rawInvite, inviteFound := db.FindOne("invites", map[string]any{
		"code": fmt.Sprintf("%x", hashedCode),
	})

	if !inviteFound {
		return Invitation{}, fmt.Errorf("could not find token")
	}

	inv := rawInvite.(map[string]any)

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
