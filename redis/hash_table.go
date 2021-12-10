package redis

import (
	"strconv"

	"github.com/yusank/godis/datastruct"
	"github.com/yusank/godis/protocol"
	"github.com/yusank/godis/util"
)

// hDel .
func hDel(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	cnt, err := datastruct.HDel(c.Values[0], c.Values[1])
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// hExists .
func hExists(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	cnt, err := datastruct.HExists(c.Values[0], c.Values[1])
	if err != nil && err != datastruct.ErrNil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil

}

// hGet .
func hGet(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	v, err := datastruct.HGet(c.Values[0], c.Values[1])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}
	var res string
	switch vv := v.(type) {
	case string:
		res = vv
	case int64:
		res = strconv.FormatInt(vv, 10)
	case float64:
		res = strconv.FormatFloat(vv, 'g', -1, 64)
	}

	return protocol.NewResponseWithBulkString(res), nil
}

// hGetAll .
func hGetAll(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	m, err := datastruct.HGetAll(c.Values[0])
	if err == datastruct.ErrNil || len(m) == 0 {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	var (
		values = make([]interface{}, len(m)*2)
		i      int
	)

	for k, v := range m {
		values[i] = k
		values[i+1] = v
		i += 2
	}

	return protocol.NewResponse(true).AppendBulkInterfaces(values...), nil
}

// hIncrBy .
func hIncrBy(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key     = c.Values[0]
		field   = c.Values[1]
		incrStr = c.Values[2]
	)

	i64, err := strconv.ParseInt(incrStr, 10, 64)
	if err != nil {
		return nil, datastruct.ErrNotInteger
	}

	res, err := datastruct.HIncrBy(key, field, i64)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(res), nil
}

// hIncrByFloat .
func hIncrByFloat(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key     = c.Values[0]
		field   = c.Values[1]
		incrStr = c.Values[2]
	)

	f64, err := strconv.ParseFloat(incrStr, 64)
	if err != nil {
		return nil, datastruct.ErrNotInteger
	}

	res, err := datastruct.HIncrByFloat(key, field, f64)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithBulkString(strconv.FormatFloat(res, 'g', -1, 64)), nil
}

// hKeys .
func hKeys(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	values, err := datastruct.HKeys(c.Values[0])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(values...), nil
}

// hLen .
func hLen(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	l, err := datastruct.HLen(c.Values[0])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithInteger(0), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(l)), nil
}

// hMGet .
func hMGet(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	values, err := datastruct.HMGet(c.Values[0], c.Values[1:]...)
	if err == datastruct.ErrNil {
		values = make([]interface{}, len(c.Values[1:]))
		err = nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkInterfaces(values...), nil
}

// hMSet .
func hMSet(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key = c.Values[0]
		kvs = make([]*datastruct.KV, 0)
	)

	for i := 1; i < len(c.Values)-1; i += 2 {
		kvs = append(kvs, &datastruct.KV{
			Key:   c.Values[i],
			Value: util.ConvertToValidValue(c.Values[i+1]),
		})
	}

	_, err := datastruct.HSet(key, datastruct.HSetInNone, kvs...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithSimpleString(RespOK), nil
}

// hSet .
func hSet(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key = c.Values[0]
		kvs = make([]*datastruct.KV, 0)
	)

	for i := 1; i < len(c.Values)-1; i += 2 {
		kvs = append(kvs, &datastruct.KV{
			Key:   c.Values[i],
			Value: util.ConvertToValidValue(c.Values[i+1]),
		})
	}

	cnt, err := datastruct.HSet(key, datastruct.HSetInNone, kvs...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// hSetNx .
func hSetNx(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key = c.Values[0]
		kvs = make([]*datastruct.KV, 1)
	)

	kvs[0] = &datastruct.KV{
		Key:   c.Values[1],
		Value: util.ConvertToValidValue(c.Values[2]),
	}

	cnt, err := datastruct.HSet(key, datastruct.HSetInNx, kvs...)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// hVals .
func hVals(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	values, err := datastruct.HVals(c.Values[0])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithEmptyArray(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(values...), nil
}
