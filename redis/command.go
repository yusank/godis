package redis

import (
	"context"
	"log"
	"strings"

	"github.com/yusank/godis/protocol"
)

type Command struct {
	Command string
	Keys    []string
	Values  []string
	Options []string
}

// sadd key1 value1
// hadd key1 hkey hvalue
// zadd key1 value1 score

func NewCommandFromMsg(msg *protocol.Message) *Command {
	c := new(Command)
	for _, e := range msg.Elements {
		if c.Command == "" {
			c.Command = strings.ToLower(e.Value)
			continue
		}

		c.Values = append(c.Values, e.Value)
	}

	log.Println("command: ", c.Command)
	return c
}

func (c *Command) Excute(ctx context.Context) (rsp []byte, err error) {
	f, ok := KnownCommands[c.Command]
	if !ok {
		return []byte(protocol.EncodeDataWithError(ErrUnknownCommand)), nil
	}

	return f(c)
}
