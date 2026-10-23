package main

import (
	"custom-redis/commands"
	"errors"
	"sync"
)

var commandMu sync.Mutex

func HandleMessage(message string) (string, error) {
	request, err := commands.ParseCommand(message)
	if err != nil {
		return "", err
	}

	commandMu.Lock()
	defer commandMu.Unlock()

	switch request.Command {
	case commands.CommandPing:
		return "PONG", nil

	case commands.CommandGet:
		return commands.ProcessGet(request)

	case commands.CommandSet:
		return commands.ProcessSet(request)

	case commands.CommandDel:
		return commands.ProcessDelete(request)

	case commands.CommandExists:
		return commands.ProcessExists(request)

	case commands.CommandIncr:
		return commands.ProcessIncrement(request)

	case commands.CommandDecr:
		return commands.ProcessDecrement(request)

	case commands.CommandExpire:
		return commands.ProcessExpire(request)

	case commands.CommandTTL:
		return commands.ProcessTTL(request)

	default:
		return "", errors.New("Unknown command: " + string(request.Command))
	}
}
