package datastruct

import (
	"sync"
)

type String struct {
	m  sync.Map
	mc *MemoryCache
}

func newString(mc *MemoryCache) *String {
	return &String{
		m:  sync.Map{},
		mc: mc,
	}
}

var (
	StringImpl = newString(defaultCache)
)
