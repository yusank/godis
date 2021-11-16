package datastruct

import (
	"sync"
)

// MemoryCache 总数结构
type MemoryCache struct {
	Keys sync.Map
}
