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
	Type KeyType
}

func (m *MemoryCache) Exists(key string) bool {
	_, ok := m.keys.Load(key)
	return ok
}

func (m *MemoryCache) Set(key string, kt KeyType) {
	m.keys.Store(key, &KeyInfo{Type: kt})
}

func (m *MemoryCache) Type(key string) (kt KeyType, found bool) {
	v, ok := m.keys.Load(key)
	if !ok {
		return "", false
	}

	return v.(*KeyInfo).Type, true
}
