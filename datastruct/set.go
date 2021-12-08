package datastruct

import (
	"sync"

	cm "github.com/yusank/concurrent-map"
)

type set struct {
	// m is a concurrentMap use as hash map and store set member as key
	m      cm.ConcurrentMap
	length int
}

func newSet() *set {
	return &set{
		m: cm.New(),
	}
}

func (s *set) sAdd(key string) int {
	if s.m.SetIfAbsent(key, 0) {
		s.length++
		return 1
	}

	return 0
}

func (s *set) sIsMember(key string) bool {
	return s.m.Has(key)
}

func (s *set) sRem(key string) int {
	remove := s.m.RemoveCb(key, func(_ string, _ interface{}, exists bool) bool {
		return exists
	})

	if remove {
		s.length--
		return 1
	}

	return 0
}

func sDiff(s1, s2 *set) *set {
	var result = newSet()
	s1.m.IterCb(func(key string, _ interface{}) {
		if !s2.m.Has(key) {
			result.sAdd(key)
		}
	})

	return result
}

func sInter(s1, s2 *set) *set {
	var result = newSet()
	s1.m.IterCb(func(key string, _ interface{}) {
		if s2.m.Has(key) {
			result.sAdd(key)
		}
	})

	return result
}

func sUnion(sets ...*set) *set {
	var result = newSet()
	for _, s := range sets {
		s.m.IterCb(func(key string, _ interface{}) {
			result.sAdd(key)
		})
	}

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
		defaultCache.keys.Set(key, &KeyInfo{
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

func SIsMember(key, value string) (int, error) {
	s, err := loadAndCheckSet(key, true)
	if err != nil {
		return 0, err
	}

	ok := s.m.Has(value)
	if ok {
		return 1, nil
	}

	return 0, nil
}

func SMembers(key string) ([]string, error) {
	s, err := loadAndCheckSet(key, true)
	if err != nil {
		return nil, err
	}

	return s.m.Keys(), nil
}

// SMove move value from source to target
func SMove(source, target, value string) (int, error) {
	sset, err := loadAndCheckSet(source, true)
	if err != nil {
		return 0, err
	}

	loaded := sset.m.RemoveCb(value, func(_ string, _ interface{}, exists bool) bool {
		return exists
	})
	if !loaded {
		// not exist
		return 0, ErrNil
	}

	return SAdd(target, value)
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

	if result == nil || result.length == 0 {
		return nil, ErrNil
	}

	return result.m.Keys(), nil
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

	defaultCache.keys.Set(storeKey, &KeyInfo{
		Type:  KeyTypeSet,
		Value: result,
	})

	return result.length, nil
}

// SInter 交集
func SInter(keys ...string) ([]string, error) {
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

		result = sInter(result, s)
	}

	if result == nil || result.length == 0 {
		return nil, ErrNil
	}

	return result.m.Keys(), nil
}

func SInterStore(storeKey string, keys ...string) (int, error) {
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

		result = sInter(result, s)
	}

	if result == nil {
		result = newSet()
	}

	defaultCache.keys.Set(storeKey, &KeyInfo{
		Type:  KeyTypeSet,
		Value: result,
	})

	return result.length, nil
}

// SUnion 并集
func SUnion(keys ...string) ([]string, error) {
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

		result = sUnion(result, s)
	}

	if result == nil || result.length == 0 {
		return nil, ErrNil
	}

	return result.m.Keys(), nil
}

func SUnionStore(storeKey string, keys ...string) (int, error) {
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

		result = sUnion(result, s)
	}

	if result == nil {
		result = newSet()
	}

	defaultCache.keys.Set(storeKey, &KeyInfo{
		Type:  KeyTypeSet,
		Value: result,
	})

	return result.length, nil
}

/*
 * compare with sync map
 * just for compare, not using in any data struct
 */

type setSyncMap struct {
	m      sync.Map
	length int
}

func newSetSyncMap() *setSyncMap {
	return &setSyncMap{
		m: sync.Map{},
	}
}

func (s *setSyncMap) sAdd(key string) int {
	_, loaded := s.m.LoadOrStore(key, 0)
	if loaded {
		return 0
	}

	s.length++
	return 1
}

func (s *setSyncMap) sIsMember(key string) bool {
	_, found := s.m.Load(key)
	return found
}

func (s *setSyncMap) sRem(key string) int {
	_, loaded := s.m.LoadAndDelete(key)
	if loaded {
		s.length--
		return 1
	}

	return 0
}
