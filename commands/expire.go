package commands

import (
	"custom-redis/store"
	"strconv"
)

//exists
//1. seconds>0: set the expiration time and return 1
//2. seconds<=0: delete the key and return 1

// dont exists
// 1. return "0"
func ProcessExpire(request CommandRequest) (string, error) {
	entry := store.GetEntry(request.Key)
	if entry.Value == "" {
		return "0", nil
	}

	seconds, err := strconv.Atoi(request.Value)
	if err != nil {
		return "", err
	}

	if seconds <= 0 {
		store.DeleteEntry(request.Key)
	} else {
		store.SetExpiration(request.Key, seconds)
	}
	return "1", nil
}
