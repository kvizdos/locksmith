package magic

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/kvizdos/locksmith/authentication/signing"
)

func TestMain(m *testing.M) {
	sp, _ := signing.CreateSigningPackage()

	MagicSigningPackage = &sp
	m.Run()
}

func TestMagicCreateSuccess(t *testing.T) {
	mac, token, err := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"permission"},
		TTL:                time.Hour,
	})

	if err != nil {
		t.Errorf("got a bad error code: %s", err.Error())
		return
	}

	if len(token) == 0 {
		t.Error("expected to receive a token!")
		return
	}

	if mac.Code == token {
		t.Error("code is not hashed!")
	}

	if len(mac.AllowedPermissions) != 1 {
		t.Errorf("got incorrect number of allowed permissions: %d", len(mac.AllowedPermissions))
		return
	}
	if mac.AllowedPermissions[0] != "permission" {
		t.Errorf("got incorrect allowed permission: %s", mac.AllowedPermissions[0])
		return
	}

	decodedToken, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		t.Errorf("got an error while decoding base64: %s", err)
		return
	}

	tokenSections := strings.Split(string(decodedToken), ":")

	if !MagicSigningPackage.Validate(tokenSections[0]+":"+tokenSections[1]+":"+tokenSections[2], tokenSections[3]) {
		t.Error("should have valid signature!")
	}
}
