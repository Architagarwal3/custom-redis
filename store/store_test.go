package store

import (
	"testing"
	"time"
)

func TestSetEntryAndGetEntry(t *testing.T) {
	const key = "store-set-get-test-key"
	cleanupTestEntry(t, key)

	want := Entry{Value: "test value"}
	SetEntry(key, want)

	if got := GetEntry(key); got != want {
		t.Errorf("GetEntry(%q) = %+v, want %+v", key, got, want)
	}
}

func TestGetEntryMissingKey(t *testing.T) {
	const key = "store-get-missing-test-key"
	DeleteEntry(key)

	if got := GetEntry(key); got != (Entry{}) {
		t.Errorf("GetEntry(%q) = %+v, want zero-value Entry", key, got)
	}
}

func TestGetEntryDeletesExpiredEntry(t *testing.T) {
	const key = "store-get-expired-test-key"
	cleanupTestEntry(t, key)
	SetEntry(key, Entry{
		Value:     "expired value",
		ExpiresAt: time.Now().Add(-time.Second),
	})

	if got := GetEntry(key); got != (Entry{}) {
		t.Errorf("GetEntry(%q) = %+v, want zero-value Entry for expired key", key, got)
	}
	if ExistsEntry(key) {
		t.Errorf("expired key %q still exists after GetEntry", key)
	}
}

func TestExistsEntry(t *testing.T) {
	const key = "store-exists-test-key"
	cleanupTestEntry(t, key)

	if ExistsEntry(key) {
		t.Fatalf("ExistsEntry(%q) = true for missing key, want false", key)
	}

	SetEntry(key, Entry{Value: "value"})
	if !ExistsEntry(key) {
		t.Errorf("ExistsEntry(%q) = false for existing key, want true", key)
	}
}

func TestExistsEntryDeletesExpiredEntry(t *testing.T) {
	const key = "store-exists-expired-test-key"
	cleanupTestEntry(t, key)
	SetEntry(key, Entry{
		Value:     "expired value",
		ExpiresAt: time.Now().Add(-time.Second),
	})

	if ExistsEntry(key) {
		t.Errorf("ExistsEntry(%q) = true for expired key, want false", key)
	}
}

func TestDeleteEntry(t *testing.T) {
	const key = "store-delete-test-key"
	SetEntry(key, Entry{Value: "value"})

	DeleteEntry(key)

	if ExistsEntry(key) {
		t.Errorf("key %q still exists after DeleteEntry", key)
	}
}

func TestSetExpiration(t *testing.T) {
	const key = "store-set-expiration-test-key"
	cleanupTestEntry(t, key)
	SetEntry(key, Entry{Value: "value"})

	start := time.Now()
	SetExpiration(key, 60)
	entry := GetEntry(key)

	if entry.Value != "value" {
		t.Errorf("value after SetExpiration = %q, want %q", entry.Value, "value")
	}
	if entry.ExpiresAt.Before(start.Add(60*time.Second)) ||
		entry.ExpiresAt.After(time.Now().Add(60*time.Second)) {
		t.Errorf("expiration = %v, want approximately 60 seconds from now", entry.ExpiresAt)
	}
}

func TestSetExpirationForMissingKey(t *testing.T) {
	const key = "store-set-expiration-missing-test-key"
	DeleteEntry(key)

	SetExpiration(key, 60)

	if ExistsEntry(key) {
		t.Errorf("SetExpiration created missing key %q", key)
	}
}

func cleanupTestEntry(t *testing.T, key string) {
	t.Helper()
	DeleteEntry(key)
	t.Cleanup(func() {
		DeleteEntry(key)
	})
}
