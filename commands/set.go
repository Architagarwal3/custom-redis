package commands

import (
	"custom-redis/store"
	"time"
)

func ProcessSet(request CommandRequest) (string, error) {
	store.SetEntry(request.Key, store.Entry{
		Value:     request.Value,
		ExpiresAt: time.Time{},
	})
	return "OK", nil
}
