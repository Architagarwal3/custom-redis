package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Command string

const (
	CommandPing   Command = "PING"
	CommandGet    Command = "GET"
	CommandSet    Command = "SET"
	CommandDel    Command = "DEL"
	CommandExists Command = "EXISTS"
	CommandIncr   Command = "INCR"
	CommandDecr   Command = "DECR"
	CommandExpire Command = "EXPIRE"
	CommandTTL    Command = "TTL"
)

type CommandRequest struct {
	Command Command
	Key     string
	Value   string
}

func ParseCommand(message string) (CommandRequest, error) {
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

	case CommandGet, CommandDel, CommandExists, CommandIncr, CommandDecr, CommandTTL:
		if len(parts) != 2 {
			return CommandRequest{}, fmt.Errorf("%s requires one key", command)
		}
		return CommandRequest{Command: command, Key: parts[1]}, nil

	case CommandSet:
		if len(parts) != 3 {
			return CommandRequest{}, fmt.Errorf("%s requires a key and a value", command)
		}
		return CommandRequest{
			Command: command,
			Key:     parts[1],
			Value:   parts[2],
		}, nil

	case CommandExpire:
		if len(parts) != 3 {
			return CommandRequest{}, fmt.Errorf("%s requires a key and a value", command)
		}

		_, err := strconv.Atoi(parts[2])
		if err != nil {
			return CommandRequest{}, fmt.Errorf("invalid expiration time: %s", parts[2])
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
