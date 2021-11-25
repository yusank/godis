package redis

import "errors"

var (
	ErrUnknownCommand       = errors.New("unknown command key")
	ErrCommandArgsNotEnough = errors.New("command args not enough")
	ErrValueOutOfRange      = errors.New(" ERR value is out of range, must be positive")
)

var (
	RespOK      = "OK"
	RespNone    = "none"
	RespPong    = "PONG"
	RespCommand = "COMMAND"
)
