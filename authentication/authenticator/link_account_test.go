package authenticator

import (
	"context"
	"testing"
	"time"
)

func TestLinkAccount_InsertsIntoAuthLinks(t *testing.T) {
	t.Parallel()
	db := newTestDB(map[string]map[string]any{
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	err := a.LinkAccount(context.Background(), "user-abc", "google", "https://accounts.google.com", "sub-xyz")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// The TestDatabase.InsertOne always succeeds and adds to the table.
	links := db.Tables["auth_links"]
	if len(links) == 0 {
		t.Fatal("expected an auth_link to be inserted")
	}

	// Find the inserted link (there's only one).
	var inserted map[string]any
	for _, v := range links {
		inserted = v.(map[string]any)
	}

	if inserted["provider"] != "google" {
		t.Errorf("expected provider 'google', got %v", inserted["provider"])
	}
	if inserted["subject"] != "sub-xyz" {
		t.Errorf("expected subject 'sub-xyz', got %v", inserted["subject"])
	}
	if inserted["user_id"] != "user-abc" {
		t.Errorf("expected user_id 'user-abc', got %v", inserted["user_id"])
	}
	if inserted["issuer"] != "https://accounts.google.com" {
		t.Errorf("expected issuer 'https://accounts.google.com', got %v", inserted["issuer"])
	}

	// linked_at should be a recent unix timestamp (int64).
	linkedAt, ok := inserted["linked_at"].(int64)
	if !ok {
		t.Fatalf("expected linked_at to be int64, got %T", inserted["linked_at"])
	}
	now := time.Now().UTC().Unix()
	if linkedAt > now || linkedAt < now-10 {
		t.Errorf("linked_at timestamp %d is too far from now %d", linkedAt, now)
	}
}

func TestLinkAccount_MissingTableStillSucceeds(t *testing.T) {
	t.Parallel()
	// TestDatabase auto-creates tables on InsertOne.
	db := newTestDB(nil)
	a := newTestAuthorizer(db)

	err := a.LinkAccount(context.Background(), "user-1", "github", "https://github.com", "gh-sub-1")
	if err != nil {
		t.Fatalf("expected no error even when table absent, got: %v", err)
	}

	// The table should now exist.
	if _, ok := db.Tables["auth_links"]; !ok {
		t.Fatal("expected auth_links table to be created automatically")
	}
}

func TestLinkAccount_ExpectedKeys(t *testing.T) {
	t.Parallel()
	db := newTestDB(map[string]map[string]any{
		"auth_links": {},
	})
	a := newTestAuthorizer(db)

	_ = a.LinkAccount(context.Background(), "u1", "provider", "https://issuer", "subj")

	var inserted map[string]any
	for _, v := range db.Tables["auth_links"] {
		inserted = v.(map[string]any)
	}

	requiredKeys := []string{"provider", "subject", "user_id", "issuer", "linked_at"}
	for _, key := range requiredKeys {
		if _, found := inserted[key]; !found {
			t.Errorf("expected key %q in inserted auth_link, but it was missing", key)
		}
	}
}
