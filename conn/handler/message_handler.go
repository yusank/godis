package handler

import (
	"github.com/yusank/godis/conn"
	"github.com/yusank/godis/protocal"
)

type MessageHandler struct{}

func (MessageHandler) Handle(r conn.Reader) ([]byte, error) {
	msg, err := protocal.NewMessageFromReader(r)
	if err != nil {
		return nil, err
	}

	err = msg.Encode()
	return msg.Bytes(), err
}
