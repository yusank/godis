package redis

import (
	"context"
	"log"
	"strings"

	"github.com/yusank/godis/protocol"
)

type Command struct {
	Command string
	Values  []string
}

// ExecuteFunc define command execute, returns slice of string as result and error if any error occur .
type ExecuteFunc func(*Command) (*protocol.Response, error)

var implementedCommands = map[string]ExecuteFunc{}

// sadd key1 value1
// hadd key1 hkey hvalue
// zadd key1 value1 score

func NewCommandFromReceive(rec protocol.Receive) *Command {
	c := new(Command)
	for _, e := range rec {
		if c.Command == "" {
			c.Command = strings.ToLower(e)
			continue
		}

		c.Values = append(c.Values, e)
	}

	log.Println("command: ", c.Command)
	return c
}

// 1. check cmd is valid
// 2. found cmd excute func

// Execute only return rsp bytes
// if got any error when execution will transfer protocol bytes
func (c *Command) Execute(ctx context.Context) *protocol.Response {
	f, ok := implementedCommands[c.Command]
	if !ok {
		return protocol.NewResponseWithError(ErrUnknownCommand)
	}

	rsp, err := f(c)
	if err != nil {
		return protocol.NewResponseWithError(err)
	}

	return rsp
}

// PrintSupportedCmd call on debug mode
func PrintSupportedCmd() {
	log.Println("[redis/command] supported cmd count: ", len(implementedCommands))
	for c := range implementedCommands {
		log.Println("[redis/command] support: ", c)
	}
}
