package commands

import "custom-redis/store"

func ProcessGet(request CommandRequest) (string, error) {
	entry := store.GetEntry(request.Key)
	if entry.Value == "" {
		return "(nil)", nil
	}
	return entry.Value, nil
}
