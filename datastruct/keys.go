package datastruct

// keys file put functions for key operation

import (
	"github.com/yusank/glob"

	smap "github.com/yusank/godis/lib/shard_map"
)

var defaultCache = newMemoryCache()

// MemoryCache 总数结构
type MemoryCache struct {
	keys smap.Map
}

func newMemoryCache() *MemoryCache {
	return &MemoryCache{keys: smap.New()}
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
	defaultCache.keys.Range(func(key string, _ interface{}) bool {
		if g.Match(key) {
			keys = append(keys, key)
		}
		return true
	})

	return
}

func Del(key string) int {
	remove := defaultCache.keys.DeleteIfExists(key)

	if remove {
		return 1
	}

	return 0
}
