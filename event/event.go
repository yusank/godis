package event

import (
	"github.com/yusank/godis/protocol"
	"github.com/yusank/godis/redis"
)

type Event struct {
	Cmd  *redis.Command
	Rsp  *protocol.Response
	done chan struct{}
}

func NewEvent(cmd *redis.Command) *Event {
	return &Event{
		Cmd:  cmd,
		done: make(chan struct{}),
	}
}

func (e *Event) Done() <-chan struct{} {
	return e.done
}
