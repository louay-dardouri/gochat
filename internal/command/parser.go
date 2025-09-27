package command

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

const (
	Prefix  = "/"
	CmdNick = "nick"
	CmdSend = "send"
	CmdView = "view"
)

func Parse(msg string) (*Command, error) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return nil, nil
	}

	if !strings.HasPrefix(msg, Prefix) {
		return &Command{
			Name: CmdSend,
			Args: strings.Fields(msg),
		}, nil
	}

	msg = msg[1:]
	parts := strings.Fields(msg)
	commandName := strings.ToLower(parts[0])
	args := parts[1:]

	switch commandName {
	case CmdNick, CmdSend, CmdView:
		return &Command{
			Name: commandName,
			Args: args,
		}, nil
	default:
		return nil, fmt.Errorf("unknown command %s%s", Prefix, commandName)
	}
}
