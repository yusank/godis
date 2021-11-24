package redis

import (
	"github.com/yusank/godis/protocol"
)

// ExecuteFunc define command execute, returns slice of string as result and error if has any error occur .
type ExecuteFunc func(*Command) (*protocol.Response, error)

func defaultExecFunc(c *Command) (*protocol.Response, error) {
	return protocol.NewResponseWithSimpleString(RespOK), nil
}

var knownCommands = map[string]ExecuteFunc{
	// native
	"command": func(c *Command) (*protocol.Response, error) {
		return protocol.NewResponseWithSimpleString(RespCommand), nil
	},
	"ping": func(c *Command) (*protocol.Response, error) {
		return protocol.NewResponseWithSimpleString(RespPong), nil
	},
	"keys":   keys,
	"exists": exists,
	"type":   keyType,
	// strings
	"append": stringAppend,
	"incr":   incr,
	"incrby": incrBy,
	"decr":   decr,
	"decrby": decrBy,
	"get":    get,
	"mget":   mget,
	"set":    set,
	//... more
}
