package server

import (
	"context"
	"log"
	"time"

	"github.com/yusank/godis/debug"
	"github.com/yusank/godis/protocol"
	"github.com/yusank/godis/redis"
)

// handleRequest receive protocol.Receive as params and return response
func handleRequest(rec protocol.Receive) *protocol.Response {
	// io data to protocol msg
	if debug.DEBUG {
		log.Println(rec)
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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

	return rsp
}
