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

	msg := protocol.NewMessageFromSimpleStrings("PING")
	err := connAndSendMsg(debugAddr, msg)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_BulkString(t *testing.T) {
	srv := startServer(debugAddr, t)

	msg := protocol.NewMessageFromBulkStrings("GET")
	err := connAndSendMsg(debugAddr, msg)
	assert.NoError(t, err)
	srv.Stop()
}

func Test_Array(t *testing.T) {
	srv := startServer(debugAddr, t)

	msg := protocol.NewMessageFromBulkStrings(
		"MGET",
		"key1",
	)
	err := connAndSendMsg(debugAddr, msg)
	assert.NoError(t, err)
	srv.Stop()
}
