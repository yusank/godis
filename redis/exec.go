package redis

// exec 实现 Command 到 真正执行底层数据操作的过程
func exec(c *Command) (reply []interface{}, err error) {
	f, ok := KnownCommands[c.Command]
	if !ok {
		return nil, ErrUnknownCommand
	}

	result, err := f(c)
	if err != nil {
		return nil, err
	}

	return result, nil
}
