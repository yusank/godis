package handler

import (
	"context"
	"log"

	"github.com/yusank/godis/conn"
	"github.com/yusank/godis/protocol"
	"github.com/yusank/godis/redis"
)

type MessageHandler struct{}

func (MessageHandler) Handle(r conn.Reader) ([]byte, error) {
	// io data to protocol msg
	msg, err := protocol.NewMessageFromReader(r)
	if err != nil {
		return nil, err
	}
	log.Println(msg)

	c := redis.NewCommandFromMsg(msg)
	rsp, err := c.Excute(context.Background())
	if err != nil {
		return nil, err
	}

	return rsp, err
}
