package commands

import "custom-redis/store"

func ProcessExists(request CommandRequest) (string, error) {
	if store.ExistsEntry(request.Key) {
		return "1", nil
	}
	return "0", nil
}
