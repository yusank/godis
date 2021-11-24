package redis

import (
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

	val, err := datastruct.LPop(c.Values[0])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithBulkString(val), nil
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

	val, err := datastruct.RPop(c.Values[0])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithBulkString(val), nil
}

// lLen .
func lLen(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	ln, err := datastruct.LLen(c.Values[0])
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(ln)), nil
}
