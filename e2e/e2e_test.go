package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yusank/godis/event"
	"github.com/yusank/godis/protocol"
)

var (
	debugAddr = ":7379"
)

func Test_SimpleString(t *testing.T) {
	srv := startServer(debugAddr, t)
	event.SetGlobalEventPool(event.NewEventPool())

	msg := protocol.NewResponseWithSimpleString("PING")
	err := connAndSendMsg(debugAddr, msg)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_BulkString(t *testing.T) {
	srv := startServer(debugAddr, t)
	event.SetGlobalEventPool(event.NewEventPool())

	rsp := protocol.NewResponse().AppendBulkInterfaces("GET")
	err := connAndSendMsg(debugAddr, rsp)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_Array(t *testing.T) {
	srv := startServer(debugAddr, t)
	event.SetGlobalEventPool(event.NewEventPool())

	rsp := protocol.NewResponse().AppendBulkInterfaces("MEGT", "key1")
	err := connAndSendMsg(debugAddr, rsp)
	assert.NoError(t, err)
	srv.Stop()
}
