package handler

import (
	"context"
	"log"

	"github.com/yusank/godis/conn"
	"github.com/yusank/godis/debug"
	"github.com/yusank/godis/protocol"
	"github.com/yusank/godis/redis"
)

type MessageHandler struct{}

func (MessageHandler) Handle(r conn.Reader) ([]byte, error) {
	// io data to protocol msg
	rec, err := protocol.DecodeFromReader(r)
	if err != nil {
		return nil, err
	}
	log.Println(rec)

	rsp := redis.NewCommandFromReceive(rec).Execute(context.Background())
	log.Println("rsp:", debug.Escape(string(rsp.Encode())))
	return rsp.Encode(), err
}
