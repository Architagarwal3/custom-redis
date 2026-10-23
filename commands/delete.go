package commands

import "custom-redis/store"

func ProcessDelete(request CommandRequest) (string, error) {
	if !store.ExistsEntry(request.Key) {
		return "0", nil
	}
	store.DeleteEntry(request.Key)
	return "1", nil
}
