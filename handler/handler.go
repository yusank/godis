package handler

import (
	"github.com/yusank/godis/conn"
)

type Handler interface {
	// Handle return reply and error
	Handle(r conn.Reader) ([]byte, error)
}
