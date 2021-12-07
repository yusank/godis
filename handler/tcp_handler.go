package handler

import (
	"context"
	"log"
	"time"

	"github.com/yusank/godis/api"
	"github.com/yusank/godis/debug"
	"github.com/yusank/godis/protocol"
	"github.com/yusank/godis/redis"
)

type TCPHandler struct{}

func (TCPHandler) Handle(r api.Reader) ([]byte, error) {
	// io data to protocol msg
	rec, err := protocol.DecodeFromReader(r)
	if err != nil {
		return nil, err
	}
	log.Println(rec)

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	// prepare cmd
	cmd := redis.NewCommandFromReceive(rec)
	rspChan := cmd.ExecuteWithContext(ctx)

	// wait for result or context timeout
	var rsp *protocol.Response
	select {
	case <-ctx.Done():
		rsp = protocol.NewResponseWithError(ctx.Err())
	case rsp = <-rspChan:
	}

	if debug.DEBUG {
		log.Println("rsp:", debug.Escape(string(rsp.Encode())))
	}

	return rsp.Encode(), err
}
