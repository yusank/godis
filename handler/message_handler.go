package handler

import (
	"log"

	"github.com/yusank/godis/conn"
	"github.com/yusank/godis/protocol"
)

type MessageHandler struct{}

func (MessageHandler) Handle(r conn.Reader) ([]byte, error) {
	// io data to protocol msg
	msg, err := protocol.NewMessageFromReader(r)
	if err != nil {
		return nil, err
	}

	if len(msg.Elements) > 1 && msg.Elements[1].Value == "ping" {
		return []byte(protocol.Pong), nil
	}

	for _, e := range msg.Elements {
		log.Println(string(e.Description), e.Value)
	}

	err = msg.Encode()
	return msg.Bytes(), err
}
