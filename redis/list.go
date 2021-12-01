package redis

import (
	"strconv"
	"strings"

	"github.com/yusank/godis/datastruct"
	"github.com/yusank/godis/protocol"
)

// lPush .
func lPush(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	ln, err := datastruct.LPush(c.Values[0], c.Values[1:]...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(ln)), nil
}

// lPop .
func lPop(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	var count = 1
	if len(c.Values) > 1 {
		cnt, err := strconv.Atoi(c.Values[1])
		if err != nil || cnt <= 0 {
			return nil, ErrValueOutOfRange
		}

		count = cnt
	}

	values, err := datastruct.LPop(c.Values[0], count)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	rsp := protocol.NewResponse()
	rsp.AppendBulkStrings(values...)
	if len(values) > 0 {
		rsp.SetIsArray()
	}

	return rsp, nil
}

// rPush .
func rPush(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	ln, err := datastruct.RPush(c.Values[0], c.Values[1:]...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(ln)), nil
}

// rPop .
func rPop(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	var count = 1
	if len(c.Values) > 1 {
		cnt, err := strconv.Atoi(c.Values[1])
		if err != nil || cnt <= 0 {
			return nil, ErrValueOutOfRange
		}

		count = cnt
	}

	values, err := datastruct.RPop(c.Values[0], count)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	rsp := protocol.NewResponse()
	rsp.AppendBulkStrings(values...)
	if len(values) > 0 {
		rsp.SetIsArray()
	}
	return rsp, nil
}

// lLen .
func lLen(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	ln, err := datastruct.LLen(c.Values[0])
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(ln)), nil
}

// lRange .
func lRange(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	start, err := strconv.Atoi(c.Values[1])
	if err != nil {
		return nil, ErrValueOutOfRange
	}

	stop, err := strconv.Atoi(c.Values[2])
	if err != nil {
		return nil, ErrValueOutOfRange
	}

	values, err := datastruct.LRange(c.Values[0], start, stop)
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	if len(values) == 0 {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	rsp := protocol.NewResponse(true)
	rsp.AppendBulkStrings(values...)
	return rsp, nil
}

// lRem .
func lRem(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	count, err := strconv.Atoi(c.Values[1])
	if err != nil {
		return nil, ErrValueOutOfRange
	}

	n, err := datastruct.LRem(c.Values[0], count, c.Values[2])
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(n)), nil
}

// lIndex .
func lIndex(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	index, err := strconv.Atoi(c.Values[1])
	if err != nil {
		return nil, ErrValueOutOfRange
	}

	val, err := datastruct.LIndex(c.Values[1], index)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithBulkString(val), nil
}

// lSet .
func lSet(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	index, err := strconv.Atoi(c.Values[1])
	if err != nil {
		return nil, ErrValueOutOfRange
	}

	err = datastruct.LSet(c.Values[0], index, c.Values[2])
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithSimpleString(RespOK), nil
}

// lInsert .
func lInsert(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 4 {
		return nil, ErrCommandArgsNotEnough
	}

	var flag = 0
	switch strings.ToLower(c.Values[1]) {
	case "after":
		flag = 1
	case "before":
		flag = -1
	default:
		return nil, ErrCommandArgsNotEnough
	}

	n, err := datastruct.LInsert(c.Values[0], c.Values[2], c.Values[3], flag)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(n)), nil
}
