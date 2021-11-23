package redis

import (
	"strconv"

	"github.com/yusank/godis/datastruct"
	"github.com/yusank/godis/protocol"
)

func set(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	datastruct.Set(c.Values[0], c.Values[1], c.Values[2:]...)
	return protocol.NewResponseWithSimpleString(RespOK), nil
}

func get(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	val, err := datastruct.Get(c.Values[0])
	if err == datastruct.ErrNil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	rsp := protocol.NewResponse()
	rsp.AppendBulkStrings(val)
	return rsp, nil
}

func mget(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	values := datastruct.MGet(c.Values...)
	rsp := protocol.NewResponse()
	rsp.IsArray = true
	rsp.AppendBulkStrings(values...)

	return rsp, nil
}

func incr(c *Command) (interface{}, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	return datastruct.IncrBy(c.Values[0], 1)
}

func incrBy(c *Command) (interface{}, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	increment, err := strconv.ParseInt(c.Values[1], 10, 64)
	if err != nil {
		return nil, datastruct.ErrNotInteger
	}

	return datastruct.IncrBy(c.Values[0], increment)
}

func decr(c *Command) (interface{}, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	return datastruct.IncrBy(c.Values[0], -1)
}

func decrBy(c *Command) (interface{}, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	increment, err := strconv.ParseInt(c.Values[1], 10, 64)
	if err != nil {
		return nil, datastruct.ErrNotInteger
	}

	return datastruct.IncrBy(c.Values[0], -increment)
}

func stringAppend(c *Command) (interface{}, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	return datastruct.Append(c.Values[0], c.Values[1])
}
