package handler

import (
	"fmt"
)

type PrintHandler struct{}

func (*PrintHandler) Handle(b []byte) error {
	fmt.Println(string(b))

	return nil
}
