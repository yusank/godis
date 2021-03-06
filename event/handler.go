package event

import (
	"context"
	"log"
	"time"

	"github.com/yusank/godis/debug"
	"github.com/yusank/godis/persistence"
	"github.com/yusank/godis/protocol"
	"github.com/yusank/godis/redis"
)

// HandleRequest receive protocol.Receive as params and return response
func HandleRequest(rec *protocol.Receive) []byte {
	// io data to protocol msg
	if debug.DEBUG {
		log.Println(rec)
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// todo: check if open persistence
	persistence.AOF.WriteCommand(rec)

	// prepare cmd
	cmd := redis.NewCommandFromReceive(ctx, rec)
	if cmd == nil {
		return nil
	}

	// event handler handle event one by one
	e := NewEvent(cmd)
	AddEvent(e)

	// wait for result or context timeout
	var rsp *protocol.Response
	select {
	case <-ctx.Done():
		rsp = protocol.NewResponseWithError(ctx.Err())
	case <-e.Done():
		rsp = e.Rsp
	}

	if debug.DEBUG {
		log.Println("rsp:", debug.Escape(string(rsp.Encode())))
	}

	return rsp.Encode()
}
