package commands

import (
	"custom-redis/store"
	"fmt"
	"time"
)

// exists
// 1. No expiration: return -1
// 2. expiration: return remaining time in seconds

// dont exists
// 1. return -2
func ProcessTTL(request CommandRequest) (string, error) {
	entry := store.GetEntry(request.Key)
	if entry.Value == "" {
		return "-2", nil
	}

	if entry.ExpiresAt.IsZero() {
		return "-1", nil
	}

	remaining := int(time.Until(entry.ExpiresAt).Seconds())
	if remaining <= 0 {
		return "-2", nil
	}
	return fmt.Sprintf("%d", remaining), nil
}
