package services

import (
	"errors"
	"fmt"
	"strings"
)

type Command string

const (
	CommandPing   Command = "PING"
	CommandGet    Command = "GET"
	CommandSet    Command = "SET"
	CommandDel    Command = "DEL"
	CommandExists Command = "EXISTS"
)

type CommandRequest struct {
	Command Command
	Key     string
	Value   string
}

func parseCommand(message string) (CommandRequest, error) {
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return CommandRequest{}, errors.New("empty message")
	}

	command := Command(strings.ToUpper(parts[0]))

	switch command {
	case CommandPing:
		if len(parts) != 1 {
			return CommandRequest{}, errors.New("PING takes no arguments")
		}
		return CommandRequest{Command: command}, nil

	case CommandGet, CommandDel, CommandExists:
		if len(parts) != 2 {
			return CommandRequest{}, fmt.Errorf("%s requires one key", command)
		}
		return CommandRequest{Command: command, Key: parts[1]}, nil

	case CommandSet:
		if len(parts) != 3 {
			return CommandRequest{}, errors.New("SET requires a key and a value")
		}
		return CommandRequest{
			Command: command,
			Key:     parts[1],
			Value:   parts[2],
		}, nil

	default:
		return CommandRequest{}, fmt.Errorf("unknown command: %s", command)
	}
}
