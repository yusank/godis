package redis

import (
	"github.com/yusank/godis/datastruct"
)

// ExecuteFunc define command execute, returns slice of string as result and error if has any error occur .
type ExecuteFunc func(*Command) (interface{}, error)

func defaultExecFunc(c *Command) (interface{}, error) {
	return []interface{}{RespOK}, nil
}

var KnownCommands = map[string]ExecuteFunc{
	// native
	"command": func(c *Command) (interface{}, error) {
		return RespCommand, nil
	},
	"ping": func(c *Command) (interface{}, error) {
		return RespPong, nil
	},
	// strings
	"append": defaultExecFunc,
	"incr":   defaultExecFunc,
	"decr":   defaultExecFunc,
	"get": func(command *Command) (interface{}, error) {
		if len(command.Values) < 1 {
			return nil, ErrCommandArgsNotEnough
		}

		val, err := datastruct.Get(command.Values[0])
		if err == datastruct.ErrNil {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		return val, nil
	},
	"mget": func(command *Command) (interface{}, error) {
		if len(command.Values) < 1 {
			return nil, ErrCommandArgsNotEnough
		}

		return datastruct.MGet(command.Values...), nil

	},
	"set": func(command *Command) (interface{}, error) {
		if len(command.Values) < 2 {
			return nil, ErrCommandArgsNotEnough
		}

		datastruct.Set(command.Values[0], command.Values[1], command.Values[2:]...)
		return RespOK, nil
	},
	//... more
}
