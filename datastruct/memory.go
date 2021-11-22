package datastruct

import (
	"sync"
)

var defaultCache = newMemoryCache()

// MemoryCache 总数结构
type MemoryCache struct {
	keys sync.Map
}

func newMemoryCache() *MemoryCache {
	return &MemoryCache{keys: sync.Map{}}
}

type KeyInfo struct {
	Type  KeyType
	Value interface{}
}

/*
 * Common Command
 */

func Exists(key string) bool {
	_, ok := defaultCache.keys.Load(key)
	return ok
}
func Type(key string) (kt KeyType, found bool) {
	v, ok := defaultCache.keys.Load(key)
	if !ok {
		return "", false
	}

	return v.(*KeyInfo).Type, true
}

/*
 * String Command
 */

func Set(key, value string, options ...string) {
	defaultCache.keys.Store(key, &KeyInfo{Type: KeyTypeString, Value: value})
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
