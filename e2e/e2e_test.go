package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yusank/godis/protocol"
)

var (
	debugAddr = ":7379"
)

func Test_e2e(t *testing.T) {
	srv := startServer(debugAddr, t)

	e2eSimpleString(t)
	e2eBulkString(t)
	e2eArray(t)

	srv.Stop()
}

func e2eSimpleString(t *testing.T) {
	msg := protocol.NewResponseWithSimpleString("PING")
	err := connAndSendMsg(debugAddr, msg)
	assert.NoError(t, err)
}

func e2eBulkString(t *testing.T) {
	rsp := protocol.NewResponse().AppendBulkInterfaces("GET")
	err := connAndSendMsg(debugAddr, rsp)
	assert.NoError(t, err)
}

func e2eArray(t *testing.T) {
	rsp := protocol.NewResponse().AppendBulkInterfaces("MEGT", "key1")
	err := connAndSendMsg(debugAddr, rsp)
	assert.NoError(t, err)
}
