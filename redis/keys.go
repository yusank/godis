package redis

import (
	"time"

	"github.com/yusank/godis/datastruct"
	"github.com/yusank/godis/protocol"
)

// global commands like `keys`, `exists`

func keys(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	values := datastruct.Keys(c.Values[0])
	return protocol.NewResponse(true).AppendBulkStrings(values...), nil
}

func exists(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	found := datastruct.Exists(c.Values[0])
	rsp := protocol.NewResponseWithInteger(0)
	if found {
		rsp = protocol.NewResponseWithInteger(1)
	}

	return rsp, nil
}

// del .
func del(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	cnt := datastruct.Del(c.Values[0])
	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// type
func keyType(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	kt, ok := datastruct.Type(c.Values[0])
	if !ok {
		rsp := protocol.NewResponseWithSimpleString(RespNone)
		return rsp, nil
	}

	rsp := protocol.NewResponseWithSimpleString(string(kt))
	return rsp, nil
}

func ping(c *Command) (*protocol.Response, error) {
	time.Sleep(time.Second)
	return protocol.NewResponseWithSimpleString(RespPong), nil
}

func command(c *Command) (*protocol.Response, error) {
	return protocol.NewResponseWithSimpleString(RespCommand), nil
}
