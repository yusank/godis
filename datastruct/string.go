package datastruct

import (
	"strconv"
)

/*
 * String Command
 * 底层均存 字符串,如果有整数操作在内内存处理 写入时还是按字符串
 */

func Set(key string, value string, options ...string) {
	defaultCache.keys.Store(key, &KeyInfo{Type: KeyTypeString, Value: value})
}

func Append(key, value string) (ln int, err error) {
	v, ok := defaultCache.keys.Load(key)
	if !ok {
		Set(key, value)
		return len(value), nil
	}

	info := v.(*KeyInfo)
	if info.Type != KeyTypeString {
		return 0, ErrKeyAndCommandNotMatch
	}

	ln = len(info.Value.(string)) + len(value)
	info.Value = info.Value.(string) + value

	return
}

func IncrBy(key string, incrBy int64) (value int64, err error) {
	v, ok := defaultCache.keys.Load(key)
	if !ok {
		Set(key, strconv.FormatInt(incrBy, 10))
		return incrBy, nil
	}

	info := v.(*KeyInfo)
	if info.Type != KeyTypeString {
		return 0, ErrKeyAndCommandNotMatch
	}

	i64, err1 := strconv.ParseInt(info.Value.(string), 10, 64)
	if err1 != nil {
		return 0, ErrNotInteger
	}
	value = i64 + incrBy
	info.Value = strconv.FormatInt(value, 10)

	return
}

func Get(key string) (string, error) {
	v, ok := defaultCache.keys.Load(key)
	if !ok {
		return "", ErrNil
	}

	info := v.(*KeyInfo)
	if info.Type != KeyTypeString {
		return "", ErrKeyAndCommandNotMatch
	}

	return info.Value.(string), nil
}

func MGet(keys ...string) []interface{} {
	// TODO: make result as slice of interface with len(keys)
	var result []interface{}

	for _, key := range keys {
		val, err := Get(key)
		if err == ErrNil {
			result = append(result, nil)
			continue
		}

		if err != nil {
			result = append(result, err)
			continue
		}

		result = append(result, val)
	}

	return result
}
