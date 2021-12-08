package redis

import (
	"strconv"

	"github.com/yusank/godis/datastruct"
	"github.com/yusank/godis/protocol"
)

// sAdd .
func sAdd(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	cnt, err := datastruct.SAdd(c.Values[0], c.Values[1:]...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil

}

// sCard .
func sCard(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	card, err := datastruct.SCard(c.Values[0])
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(card)), nil
}

// sDiff .
func sDiff(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	if len(c.Values) == 1 {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	diffs, err := datastruct.SDiff(c.Values...)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(diffs...), nil
}

// sDiffStore .
func sDiffStore(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	if len(c.Values) == 1 {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	cnt, err := datastruct.SDiffStore(c.Values[0], c.Values[1:]...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// sInter .
func sInter(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	if len(c.Values) == 1 {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	inter, err := datastruct.SInter(c.Values...)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(inter...), nil
}

// sInterStore .
func sInterStore(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	if len(c.Values) == 1 {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	cnt, err := datastruct.SInterStore(c.Values[0], c.Values[1:]...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// sUnion .
func sUnion(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	if len(c.Values) == 1 {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	unions, err := datastruct.SUnion(c.Values...)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(unions...), nil
}

// sUnionStore .
func sUnionStore(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	if len(c.Values) == 1 {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	cnt, err := datastruct.SUnionStore(c.Values[0], c.Values[1:]...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// sIsMember .
func sIsMember(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	res, err := datastruct.SIsMember(c.Values[0], c.Values[1])
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(res)), nil
}

// sMembers .
func sMembers(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	vals, err := datastruct.SMembers(c.Values[0])
	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(vals...), nil
}

// sMove .
func sMove(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		source = c.Values[0]
		target = c.Values[1]
		value  = c.Values[2]
	)

	res, err := datastruct.SMove(source, target, value)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(res)), nil
}

// sRem .
func sRem(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	cnt, err := datastruct.SRem(c.Values[0], c.Values[1:]...)
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// SPop .
func SPop(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key = c.Values[0]
		cnt int
		err error
	)

	if len(c.Values) > 1 {
		cnt, err = strconv.Atoi(c.Values[1])
		if err != nil {
			return nil, datastruct.ErrNotInteger
		}
	}

	values, err := datastruct.SPop(key, cnt)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithEmptyArray(), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(values...), nil
}
