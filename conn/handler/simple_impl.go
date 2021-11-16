package handler

import (
	"log"

	"github.com/yusank/godis/protocal"
)

type PrintHandler struct{}

func (PrintHandler) Handle(b []byte) ([]byte, error) {
	log.Println(string(b), len(b))
	msg, err := protocal.NewMessageFromBytes(b)
	if err != nil {
		return nil, err
	}

	err = msg.Encode()
	return msg.OriginalData, err
}
