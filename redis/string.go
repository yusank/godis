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
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponse().AppendBulkInterfaces(val), nil
}

func mget(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	values := datastruct.MGet(c.Values...)

	return protocol.NewResponse(true).AppendBulkInterfaces(values...), nil
}

func incr(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	i64, err := datastruct.IncrBy(c.Values[0], 1)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(i64), nil
}

func incrBy(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	increment, err := strconv.ParseInt(c.Values[1], 10, 64)
	if err != nil {
		return nil, datastruct.ErrNotInteger
	}

	i64, err := datastruct.IncrBy(c.Values[0], increment)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(i64), nil
}

func decr(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	i64, err := datastruct.IncrBy(c.Values[0], -1)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(i64), nil
}

func decrBy(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	decrement, err := strconv.ParseInt(c.Values[1], 10, 64)
	if err != nil {
		return nil, datastruct.ErrNotInteger
	}

	i64, err := datastruct.IncrBy(c.Values[0], -decrement)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(i64), nil
}

// append
func stringAppend(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	i, err := datastruct.Append(c.Values[0], c.Values[1])
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(i)), nil
}
