package datastruct

/*
 * String Command
 */

func Set(key string, value interface{}, options ...string) {
	defaultCache.keys.Store(key, &KeyInfo{Type: KeyTypeString, Value: value})
}

func Incr(key string, incrBy float64) (value float64, err error) {
	v, ok := defaultCache.keys.Load(key)
	if !ok {
		Set(key, incrBy)
		return incrBy, nil
	}

	info := v.(*KeyInfo)
	if info.Type != KeyTypeString {
		return 0, ErrKeyAndCommandNotMatch
	}

	info.Value = info.Value.(float64) + incrBy
	return info.Value.(float64), nil
}

func Get(key string) (interface{}, error) {
	v, ok := defaultCache.keys.Load(key)
	if !ok {
		return "", ErrNil
	}

	info := v.(*KeyInfo)
	if info.Type != KeyTypeString {
		return "", ErrKeyAndCommandNotMatch
	}

	return info.Value, nil
}

func MGet(keys ...string) []interface{} {
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
