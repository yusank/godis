package redis

// ExecuteFunc define command execute, returns slice of string as result and error if has any error occur .
type ExecuteFunc func(*Command) ([]interface{}, error)

func defaultExecFunc(c *Command) ([]interface{}, error) {
	return []interface{}{RespOK}, nil
}

var KnownCommands = map[string]ExecuteFunc{
	// native
	"command": func(c *Command) ([]interface{}, error) {
		return []interface{}{RespCommand}, nil
	},
	"ping": func(c *Command) ([]interface{}, error) {
		return []interface{}{RespPong}, nil
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
