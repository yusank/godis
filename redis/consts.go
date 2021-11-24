package redis

import "errors"

var (
	ErrUnknownCommand       = errors.New("unknown command key")
	ErrCommandArgsNotEnough = errors.New("command args not enough")
)

var (
	RespOK      = "OK"
	RespNone    = "none"
	RespPong    = "PONG"
	RespCommand = "COMMAND"
)
