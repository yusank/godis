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

// 1. check cmd is valid
// 2. found cmd excute func

// Execute only return rsp bytes
// if got any error when execution will transfer protocol bytes
func (c *Command) Execute(ctx context.Context) (rsp []byte) {
	result, err := exec(c)
	if err != nil {
		return wrapError(err)
	}
	return wrapResult(result)
}
