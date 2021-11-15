package handler

import (
	"fmt"

	"github.com/yusank/godis/protocal"
)

type PrintHandler struct{}

func (*PrintHandler) Handle(b []byte) error {
	fmt.Println(string(b), len(b))
	msg, err := protocal.NewMessage(b)
	if err != nil {
		return err
	}

	return msg.Encode()
}
