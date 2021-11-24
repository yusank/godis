package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yusank/godis/protocol"
)

var (
	debugAddr = ":7379"
)

func Test_SimpleString(t *testing.T) {
	srv := startServer(debugAddr, t)

	msg := protocol.NewResponseWithSimpleString("PING")
	err := connAndSendMsg(debugAddr, msg)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_BulkString(t *testing.T) {
	srv := startServer(debugAddr, t)

	rsp := protocol.NewResponse()
	rsp.AppendBulkInterfaces("GET")
	err := connAndSendMsg(debugAddr, rsp)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_Array(t *testing.T) {
	srv := startServer(debugAddr, t)

	rsp := protocol.NewResponse()
	rsp.AppendBulkInterfaces("MEGT", "key1")
	err := connAndSendMsg(debugAddr, rsp)
	assert.NoError(t, err)
	srv.Stop()
}
