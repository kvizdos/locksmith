package magic

import (
	"testing"
	"time"
)

func TestFilterActiveNoManualSelection(t *testing.T) {
	active := make(chan MagicAuthentications)

	magicDoExpire, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"xyz"},
		TTL:                time.Hour * -1,
	})

	magicShouldNotExpire1, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"abc"},
		TTL:                time.Hour * 1,
	})

	magicShouldNotExpire2, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"abc"},
		TTL:                time.Hour * 1,
	})

	magics := MagicAuthentications{
		magicShouldNotExpire1,
		magicDoExpire,
		magicShouldNotExpire2,
	}

	go FilterActive(active, magics)

	keep := <-active

	if len(keep) != 2 {
		t.Errorf("got incorrect number of magics back: %d", len(keep))
	}

	if keep[0].AllowedPermissions[0] != "abc" ||
		keep[1].AllowedPermissions[0] != "abc" {
		t.Errorf("got incorrect permissions.")
	}
}

func TestFilterActiveManuallyExpire(t *testing.T) {
	active := make(chan MagicAuthentications)

	magicDoExpire, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"xyz"},
		TTL:                time.Hour * -1,
	})

	magicShouldNotExpire1, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"abc"},
		TTL:                time.Hour * 1,
	})

	magicShouldNotExpire2, _, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"abc"},
		TTL:                time.Hour * 1,
	})

	magicShouldManuallyExpire, manualID, _ := CreateMagicAuthentication(MagicAuthenticationVariables{
		ForUserID:          "test-uid",
		AllowedPermissions: []string{"xyz"},
		TTL:                time.Hour * 1,
	})

	magics := MagicAuthentications{
		magicShouldNotExpire1,
		magicDoExpire,
		magicShouldManuallyExpire,
		magicShouldNotExpire2,
	}

	go FilterActive(active, magics, manualID)

	keep := <-active

	if len(keep) != 2 {
		t.Errorf("got incorrect number of magics back: %d", len(keep))
	}

	if keep[0].AllowedPermissions[0] != "abc" ||
		keep[1].AllowedPermissions[0] != "abc" {
		t.Errorf("got incorrect permissions.")
	}
}
