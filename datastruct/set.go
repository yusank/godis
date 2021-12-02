package datastruct

import (
	"sync"
)

type set struct {
	m sync.Map
}

func newSet() *set {
	return &set{m: sync.Map{}}
}

func loadSet(key string) (*set, error) {
	info, err := loadKeyInfo(key, KeyTypeSortedSet)
	if err != nil {
		return nil, err
	}

	s := info.Value.(*set)
	return s, nil
}

func SAdd(key string, values ...string) (int, error) {
	s, err := loadSet(key)
	if err == ErrNil {
		s = newSet()
		err = nil
	}

	if err != nil {
		return 0, err
	}

	var cnt int
	for _, value := range values {
		_, loaded := s.m.LoadOrStore(value, 0)
		if !loaded {
			cnt++
		}
	}

	return cnt, nil
}
