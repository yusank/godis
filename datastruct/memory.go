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
