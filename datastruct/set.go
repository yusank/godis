package datastruct

import (
	smap "github.com/yusank/godis/lib/shard_map"
)

type set struct {
	m      smap.Map
	length int
}

func newSet() *set {
	return &set{
		m: smap.New(),
	}
}

func (s *set) sAdd(key string) int {
	if s.m.SetIfAbsent(key, 0) {
		s.length++
		return 1
	}

	return 0
}

func (s *set) sRem(key string) int {
	exists := s.m.DeleteIfExists(key)

	if exists {
		s.length--
		return 1
	}

	return 0
}

func (s *set) sPop(cnt int) []string {
	var (
		result []string
	)
	s.m.RangeAndDelete(func(key string, value interface{}) (del, rng bool) {
		if cnt <= 0 {
			return false, false
		}

		result = append(result, key)
		cnt--
		return true, true
	})

	return result
}

func sDiff(s1, s2 *set) *set {
	var result = newSet()
	s1.m.Range(func(key string, _ interface{}) bool {
		if !s2.m.Has(key) {
			result.sAdd(key)
		}

		return true
	})

	return result
}

func sInter(s1, s2 *set) *set {
	var result = newSet()
	s1.m.Range(func(key string, _ interface{}) bool {
		if s2.m.Has(key) {
			result.sAdd(key)
		}

		return true
	})

	return result
}

func sUnion(sets ...*set) *set {
	var result = newSet()
	for _, s := range sets {
		s.m.Range(func(key string, _ interface{}) bool {
			result.sAdd(key)
			return true
		})

	}

	return result
}

/*
 * Commands
 */

func loadAndCheckSet(key string, check bool) (*set, error) {
	info, err := loadKeyInfo(key, KeyTypeSet)
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

	exists := sset.m.DeleteIfExists(value)
	if !exists {
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

// SInter ??????
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

// SUnion ??????
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

func SPop(key string, cnt int) ([]string, error) {
	s, err := loadAndCheckSet(key, true)
	if err != nil {
		return nil, err
	}

	return s.sPop(cnt), nil
}

/*
 * compare with sync map
 * just for compare, not using in any data struct
 */
//
//type setSyncMap struct {
//	m      sync.Map
//	length int
//}
//
//func newSetSyncMap() *setSyncMap {
//	return &setSyncMap{
//		m: sync.Map{},
//	}
//}
//
//func (s *setSyncMap) sAdd(key string) int {
//	_, loaded := s.m.LoadOrStore(key, 0)
//	if loaded {
//		return 0
//	}
//
//	s.length++
//	return 1
//}
//
//func (s *setSyncMap) sIsMember(key string) bool {
//	_, found := s.m.Load(key)
//	return found
//}
//
//func (s *setSyncMap) sRem(key string) int {
//	_, loaded := s.m.LoadAndDelete(key)
//	if loaded {
//		s.length--
//		return 1
//	}
//
//	return 0
//}

// use shardMap as example
//
//type setNonLockMap struct {
//	m      Map
//	length int
//}
//
//func newSetNonLockMap() *setNonLockMap {
//	return &setNonLockMap{
//		m: NewMap(),
//	}
//}
//
//func (s *setNonLockMap) sAdd(key string) int {
//	absent := s.m.SetIfAbsent(key, 0)
//	if !absent {
//		return 0
//	}
//
//	s.length++
//	return 1
//}
//
//func (s *setNonLockMap) sIsMember(key string) bool {
//	_, found := s.m.Get(key)
//	return found
//}
//
//func (s *setNonLockMap) sRem(key string) int {
//	exists := s.m.DeleteIfExists(key)
//	if exists {
//		s.length--
//		return 1
//	}
//
//	return 0
//}
//
//func (s *setNonLockMap) sPop(cnt int) []string {
//	var (
//		result []string
//	)
//	s.m.RangeAndDelete(func(key string, value interface{}) (del, rng bool) {
//		if cnt <= 0 {
//			return false, false
//		}
//
//		result = append(result, key)
//		cnt--
//		return true, true
//	})
//
//	return result
//}
