package handler

import (
	"log"

	"github.com/yusank/godis/conn"
	"github.com/yusank/godis/protocal"
)

type MessageHandler struct{}

func (MessageHandler) Handle(r conn.Reader) ([]byte, error) {
	msg, err := protocal.NewMessageFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(msg.Elements) > 1 && msg.Elements[1].Value == "ping" {
		return []byte(protocal.Pong), nil
	}

	for _, e := range msg.Elements {
		log.Println(string(e.Description), e.Value)
	}

	err = msg.Encode()
	return msg.Bytes(), err
}
