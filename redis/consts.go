package redis

import "errors"

var (
	ErrUnknownCommand = errors.New("unknown command key")
)

var (
	RespOK      = "OK"
	RespPong    = "PONG"
	RespCommand = "COMMAND"
)
