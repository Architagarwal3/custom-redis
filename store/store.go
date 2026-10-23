package store

import "time"

type Entry struct {
	Value     string
	ExpiresAt time.Time // zero means no expiration
}

var store = make(map[string]Entry)

func SetEntry(key string, entry Entry) {
	store[key] = entry
}

// lazy expiration check: if the entry has expired, treat it as non-existent
func getEntry(key string) (Entry, bool) {
	entry, exists := store[key]
	if exists && !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		DeleteEntry(key)
		return Entry{}, false
	}
	return entry, exists
}

func GetEntry(key string) Entry {
	entry, _ := getEntry(key)
	return entry
}

func ExistsEntry(key string) bool {
	_, exists := getEntry(key)
	return exists
}

// expiration is in seconds
func SetExpiration(key string, expiration int) {
	if !ExistsEntry(key) {
		return
	}

	entry, _ := getEntry(key)
	entry.ExpiresAt = time.Now().Add(time.Duration(expiration) * time.Second)
	SetEntry(key, entry)
}

func DeleteEntry(key string) {
	delete(store, key)
}
