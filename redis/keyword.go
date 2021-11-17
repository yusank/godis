package redis

import "github.com/yusank/godis/protocol"

type ExecuteFunc func(*Command) ([]byte, error)

func defaultExecFunc(c *Command) ([]byte, error) {
	return protocol.OK, nil
}

var KnownCommands = map[string]ExecuteFunc{
	// native
	"command": func(c *Command) ([]byte, error) {
		return protocol.Command, nil
	},
	"ping": func(c *Command) ([]byte, error) {
		return protocol.Pong, nil
	},
	// strings
	"append": defaultExecFunc,
	"incr":   defaultExecFunc,
	"decr":   defaultExecFunc,
	"get":    defaultExecFunc,
	"mget":   defaultExecFunc,
	"set":    defaultExecFunc,
	//... more
}
