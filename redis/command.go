package redis

import (
	"context"
	"log"
	"strings"

	"github.com/yusank/godis/protocol"
)

type Command struct {
	Ctx     context.Context
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

	return c
}

// 1. check cmd is valid
// 2. found cmd excute func

// ExecuteWithContext is async method which return a result channel and run command ina new go routine.
// if got any error when execution will transfer protocol bytes
func (c *Command) ExecuteWithContext(ctx context.Context) chan *protocol.Response {
	var (
		rspChan = make(chan *protocol.Response, 1)
	)

	c.Ctx = ctx
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}

	go func() {
		defer func() {
			if recover() != nil {
				log.Println(c.Command, c.Values)
			}
		}()
		f, ok := implementedCommands[c.Command]
		if !ok {
			log.Println(c.Command, c.Values)
			c.putRspToChan(rspChan, protocol.NewResponseWithError(ErrUnknownCommand))
			return
		}

		rsp, err := f(c)
		if err != nil {
			log.Println(c.Command)
			c.putRspToChan(rspChan, protocol.NewResponseWithError(err))
			return
		}

		c.putRspToChan(rspChan, rsp)
	}()

	return rspChan
}

func (c *Command) putRspToChan(ch chan *protocol.Response, rsp *protocol.Response) {
	// if ctx has error won't put
	if c.Ctx != nil && c.Ctx.Err() != nil {
		return
	}

	ch <- rsp
}

// PrintSupportedCmd call on debug mode
func PrintSupportedCmd() {
	log.Println("[redis/command] supported cmd count: ", len(implementedCommands))
	for c := range implementedCommands {
		log.Println("[redis/command] support: ", c)
	}
}
