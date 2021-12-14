// package smap provide shard_map

package smap

const shardCount = 32

// Map uses for non-lock single thread situation.
type Map []*Shard

// Shard of Map
// 分片
type Shard struct {
	items map[string]interface{}
}

func New() Map {
	m := make(Map, shardCount)
	for i := 0; i < shardCount; i++ {
		m[i] = &Shard{items: make(map[string]interface{})}
	}

	return m
}

func (m Map) Get(key string) (interface{}, bool) {
	shard := m.GetShard(key)

	v, ok := shard.items[key]
	return v, ok
}

func (m Map) Has(key string) bool {
	shard := m.GetShard(key)
	_, ok := shard.items[key]
	return ok
}

func (m Map) Set(key string, value interface{}) {
	shard := m.GetShard(key)
	shard.items[key] = value
}

func (m Map) SetIfAbsent(key string, value interface{}) bool {
	shard := m.GetShard(key)
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
		return true
	}

	return false
}

func (m Map) DeleteIfExists(key string) bool {
	shard := m.GetShard(key)
	_, ok := shard.items[key]
	if !ok {
		return false
	}

	delete(shard.items, key)
	return true
}

func (m Map) LoadAndDelete(key string) (v interface{}, loaded bool) {
	shard := m.GetShard(key)
	v, loaded = shard.items[key]
	if !loaded {
		return nil, false
	}

	delete(shard.items, key)
	return v, loaded
}

func (m Map) Delete(key string) {
	shard := m.GetShard(key)
	delete(shard.items, key)
}

func (m Map) Range(f func(key string, value interface{}) bool) {
	for i := range m {
		shard := (m)[i]
		for s, v := range shard.items {
			if !f(s, v) {
				return
			}
		}
	}
}

func (m Map) RangeAndDelete(f func(key string, value interface{}) (del, rng bool)) {
	for i := range m {
		shard := (m)[i]
		for s, v := range shard.items {
			del, rng := f(s, v)
			if del {
				delete(shard.items, s)
			}

			if !rng {
				return
			}
		}
	}
}

// UpsertFunc callback for upsert
// if after found oldValue and want to stop the upsert op, you can return result and true for it
type UpsertFunc func(exist bool, valueInMap, newValue interface{}) (result interface{}, abort bool)

// Upsert - update or insert value and support abort operation after callback
func (m Map) Upsert(key string, value interface{}, f UpsertFunc) (res interface{}, abort bool) {
	shard := m.GetShard(key)
	old, ok := shard.items[key]
	res, abort = f(ok, old, value)
	if abort {
		return
	}
	shard.items[key] = res
	return
}

func (m Map) Keys() []string {
	keys := make([]string, 0)
	for i := range m {
		shard := (m)[i]
		for key := range shard.items {
			keys = append(keys, key)
		}
	}

	return keys
}

// Items copy all k-v from Map
func (m Map) Items() map[string]interface{} {
	tmp := make(map[string]interface{})
	for i := range m {
		shard := (m)[i]
		for k, v := range shard.items {
			tmp[k] = v
		}
	}

	return tmp
}

// GetShard returns shard under given key
func (m Map) GetShard(key string) *Shard {
	return m[uint(fnv32(key))%uint(shardCount)]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
