package redis

// ExecuteFunc define command execute, returns slice of string as result and error if has any error occur .
type ExecuteFunc func(*Command) ([]string, error)

func defaultExecFunc(c *Command) ([]string, error) {
	return []string{RespOK}, nil
}

var KnownCommands = map[string]ExecuteFunc{
	// native
	"command": func(c *Command) ([]string, error) {
		return []string{RespCommand}, nil
	},
	"ping": func(c *Command) ([]string, error) {
		return []string{RespPong}, nil
	},
	// strings
	"append": defaultExecFunc,
	"incr":   defaultExecFunc,
	"decr":   defaultExecFunc,
	"get":    defaultExecFunc,
	"mget": func(c *Command) ([]string, error) {
		return c.Values, nil
	},
	"set": defaultExecFunc,
	//... more
}
