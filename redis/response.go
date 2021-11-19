package redis

import "github.com/yusank/godis/protocol"

func wrapResult(result []string) []byte {
	msg := protocol.NewMessageFromBulkStrings(result...)

	return msg.Bytes()
}

func wrapError(err error) []byte {
	msg := protocol.NewMessageFromError(err)

	return msg.Bytes()
}
