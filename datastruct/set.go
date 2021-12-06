package datastruct

import (
	"sync"
)

type set struct {
	m      sync.Map
	length int
}

func newSet() *set {
	return &set{m: sync.Map{}}
}

func (s *set) toSlice() []string {
	var (
		result = make([]string, s.length)
		i      int
	)
	s.m.Range(func(key, _ interface{}) bool {
		result[i] = key.(string)
		i++
		return true
	})

	return result
}

func (s *set) sAdd(key string) int {
	_, loaded := s.m.LoadOrStore(key, 0)
	if loaded {
		return 0
	}
	s.length++
	return 1
}

func (s *set) sRem(key string) int {
	_, loaded := s.m.LoadAndDelete(key)
	if !loaded {
		return 0
	}

	s.length--
	return 1
}

func sDiff(s1, s2 *set) *set {
	var result = newSet()
	s1.m.Range(func(key, _ interface{}) bool {
		if _, ok := s2.m.Load(key); !ok {
			result.sAdd(key.(string))
		}

		return true
	})

	return result
}

/*
 * Commands
 */

func loadAndCheckSet(key string, check bool) (*set, error) {
	info, err := loadKeyInfo(key, KeyTypeSortedSet)
	if err != nil {
		return nil, err
	}

	s := info.Value.(*set)
	if check && s.length == 0 {
		return nil, ErrNil
	}

	return s, nil
}

func SAdd(key string, values ...string) (int, error) {
	s, err := loadAndCheckSet(key, false)
	if err == ErrNil {
		s = newSet()
		defaultCache.keys.Store(key, &KeyInfo{
			Type:  KeyTypeSet,
			Value: s,
		})
		err = nil
	}

	if err != nil {
		return 0, err
	}

	var cnt int
	for _, value := range values {
		cnt += s.sAdd(value)
	}

	return cnt, nil
}

func SCard(key string) (int, error) {
	s, err := loadAndCheckSet(key, true)
	if err != nil {
		return 0, err
	}

	return s.length, nil
}

func SRem(key string, values ...string) (int, error) {
	s, err := loadAndCheckSet(key, true)
	if err != nil {
		return 0, err
	}

	var cnt int
	for _, value := range values {
		cnt += s.sRem(value)
	}

	return cnt, nil
}

func SDiff(keys ...string) ([]string, error) {
	var result *set

	for i, key := range keys {
		s, err := loadAndCheckSet(key, true)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = s
			continue
		}

		result = sDiff(result, s)
	}

	if result == nil {
		return nil, ErrNil
	}

	return result.toSlice(), nil
}

func SDiffStore(storeKey string, keys ...string) (int, error) {
	var result *set

	for i, key := range keys {
		s, err := loadAndCheckSet(key, true)
		if err != nil {
			return 0, err
		}

		if i == 0 {
			result = s
			continue
		}

		result = sDiff(result, s)
	}

	if result == nil {
		result = newSet()
	}

	defaultCache.keys.Store(storeKey, &KeyInfo{
		Type:  KeyTypeSet,
		Value: result,
	})

	return result.length, nil
}
