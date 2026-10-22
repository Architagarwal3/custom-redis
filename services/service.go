package services

import (
	"errors"
)

var store = make(map[string]string)

func HandleMessage(message string) (string, error) {
	request, err := parseCommand(message)
	if err != nil {
		return "", err
	}

	switch request.Command {
	case CommandPing:
		return "PONG", nil

	case CommandGet:
		if value, exists := store[request.Key]; exists {
			return value, nil
		}
		return "", errors.New("Key not found")

	case CommandSet:
		store[request.Key] = request.Value
		return "OK", nil

	case CommandDel:
		delete(store, request.Key)
		return "OK", nil

	case CommandExists:
		if _, exists := store[request.Key]; exists {
			return "1", nil
		}
		return "0", nil

	default:
		return "", errors.New("Unknown command: " + string(request.Command))
	}
}
