package redis

import "github.com/yusank/godis/protocol"

func wrapError(err error) []byte {
	msg := protocol.NewMessageFromError(err)

	return msg.Bytes()
}

func wrapResult(result interface{}) []byte {
	var results []interface{}
	switch v := result.(type) {
	case []interface{}:
		results = v
	default:
		results = append(results, v)
	}

	msg := protocol.NewMessageFromResults(results)

	return msg.Bytes()
}
