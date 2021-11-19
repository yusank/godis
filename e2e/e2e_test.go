package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yusank/godis/protocol"
)

var (
	_debug_addr = ":7379"
)

func Test_SimpleString(t *testing.T) {
	srv := startServer(_debug_addr, t)

	msg := protocol.NewMessageFromEncodeData(protocol.EncodeDataWithSimpleString("PING"))
	err := connAndSendMsg(_debug_addr, msg)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_BulkString(t *testing.T) {
	srv := startServer(_debug_addr, t)

	msg := protocol.NewMessageFromEncodeData(protocol.EncodeDataWithBulkString("GET"))
	err := connAndSendMsg(_debug_addr, msg)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_Array(t *testing.T) {
	srv := startServer(_debug_addr, t)

	msg := protocol.NewMessageFromEncodeData(protocol.EncodeDataWithArray(
		protocol.EncodeDataWithBulkString("MGET"),
		protocol.EncodeDataWithBulkString("key1"),
	))
	err := connAndSendMsg(_debug_addr, msg)
	assert.NoError(t, err)
	srv.Stop()
}
