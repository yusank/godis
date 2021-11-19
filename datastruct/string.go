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

func (s *String) Set(key, value string, options []string) error {
	s.mc.Set(key, KeyTypeString)

	s.m.Store(key, value)
	return nil
}

func (s *String) Get(key string) (value string, err error) {
	v, ok := s.m.Load(key)
	if !ok {
		return "", ErrNil
	}

	return v.(string), nil
}

func (s *String) MGet(keys []string) (values []string, err error) {
	for _, key := range keys {
		v, ok := s.m.Load(key)
		if !ok {
			values = append(values, "")
		}

		values = append(values, v.(string))
	}

	return
}
