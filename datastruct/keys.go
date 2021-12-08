package datastruct

// keys file put functions for key operation

import (
	cm "github.com/yusank/concurrent-map"
	"github.com/yusank/glob"
)

var defaultCache = newMemoryCache()

// MemoryCache 总数结构
type MemoryCache struct {
	keys cm.ConcurrentMap
}

func newMemoryCache() *MemoryCache {
	return &MemoryCache{keys: cm.New()}
}

type KeyInfo struct {
	Type  KeyType
	Value interface{}
}

/*
 * util functions
 */

func loadKeyInfo(key string, tp KeyType) (info *KeyInfo, err error) {
	v, ok := defaultCache.keys.Get(key)
	if !ok {
		return nil, ErrNil
	}

	info = v.(*KeyInfo)
	if info.Type != tp {
		return nil, ErrKeyAndCommandNotMatch
	}

	return info, nil
}

/*
 * Common Command
 */

func Exists(key string) bool {
	return defaultCache.keys.Has(key)
}

func Type(key string) (kt KeyType, found bool) {
	v, ok := defaultCache.keys.Get(key)
	if !ok {
		return "", false
	}

	return v.(*KeyInfo).Type, true
}

func Keys(pattern string) (keys []string) {
	g := glob.MustCompile(pattern)
	defaultCache.keys.RangeCb(func(key string, _ interface{}) bool {
		if g.Match(key) {
			keys = append(keys, key)
		}
		return true
	})

	return
}

func Del(key string) int {
	remove := defaultCache.keys.RemoveCb(key, func(key string, v interface{}, exists bool) bool {
		return exists
	})

	if remove {
		return 1
	}

	return 0
}
