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
		log.Println(string(e.Description), e.Len, e.Value)
		if e.Description == protocol.DescriptionArray {
			continue
		}

		if c.Command == "" {
			c.Command = strings.ToLower(e.Value)
			continue
		}

		c.Values = append(c.Values, e.Value)
	}

	return c
}

func (c *Command) Excute(ctx context.Context) (rsp []byte, err error) {
	f, ok := KnownCommands[c.Command]
	if !ok {
		return protocol.BuildErrorResponseBystes(UnknownCommand.Error()), nil
	}

	return f(c)
}
