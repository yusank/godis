package datastruct

import (
	"strconv"

	smap "github.com/yusank/godis/lib/shard_map"
)

type hashTable struct {
	table  smap.Map
	length int
}

// KV is an alias for key-value struct use for hash map multi set
type KV struct {
	Key   string
	Value interface{}
}

func newHashTable() *hashTable {
	return &hashTable{
		table: smap.New(),
	}
}

func (h *hashTable) del(field string) int {
	remove := !h.table.DeleteIfExists(field)
	if remove {
		h.length--
		return 1
	}

	return 0
}

func (h *hashTable) set(field string, v interface{}, flag int) int {
	// set only not exists
	if flag&HSetInNx != 0 {
		if h.table.SetIfAbsent(field, v) {
			h.length++
			return 1
		}

		return 0
	}

	if !h.table.Has(field) {
		h.length++
	}

	h.table.Set(field, v)
	return 1
}

func (h *hashTable) exists(field string) bool {
	return h.table.Has(field)
}

func (h *hashTable) get(field string) (interface{}, bool) {
	return h.table.Get(field)
}

func (h *hashTable) getAll() map[string]interface{} {
	return h.table.Items()
}

func (h *hashTable) incrBy(field string, i int64) (int64, error) {
	res, abort := h.table.Upsert(field, i, func(exist bool, valueInMap interface{}, newValue interface{}) (interface{}, bool) {
		if !exist {
			return i, false
		}

		i64, ok := valueInMap.(int64)
		if !ok {
			return ErrNotInteger, true
		}

		return i64 + newValue.(int64), false
	})

	if abort {
		// return true only return error
		return 0, res.(error)
	}

	return res.(int64), nil
}

func (h *hashTable) incrByFloat(field string, i float64) (float64, error) {
	res, abort := h.table.Upsert(field, i, func(exist bool, valueInMap interface{}, newValue interface{}) (interface{}, bool) {
		if !exist {
			return i, false
		}

		i64, ok := valueInMap.(float64)
		if !ok {
			return ErrNotFloat, true
		}

		return i64 + newValue.(float64), false
	})

	if abort {
		// return true only return error
		return 0, res.(error)
	}

	return res.(float64), nil
}

func (h *hashTable) keys() []string {
	return h.table.Keys()
}

func (h *hashTable) mGet(fields ...string) []interface{} {
	var result = make([]interface{}, len(fields))

	for i, field := range fields {
		value, _ := h.get(field)
		if value == nil {
			result[i] = value
			continue
		}
		switch v := value.(type) {
		case string:
			result[i] = v
		case int64:
			result[i] = strconv.FormatInt(v, 10)
		case float64:
			result[i] = strconv.FormatFloat(v, 'g', -1, 64)
		}
	}

	return result
}

func (h *hashTable) mSet(kvs []*KV) {
	for _, kv := range kvs {
		h.set(kv.Key, kv.Value, HSetInNone)
	}
}

func (h *hashTable) values() []string {
	var result []string
	h.table.Range(func(_ string, value interface{}) bool {
		var vStr string
		switch v := value.(type) {
		case string:
			vStr = v
		case int64:
			vStr = strconv.FormatInt(v, 10)
		case float64:
			vStr = strconv.FormatFloat(v, 'g', -1, 64)
		}
		result = append(result, vStr)

		return true
	})
	
	return result
}

/*
 * commands
 */

func loadAndCheckHashTable(key string, checkLen bool) (*hashTable, error) {
	info, err := loadKeyInfo(key, KeyTypeHashTable)
	if err != nil {
		return nil, err
	}

	h := info.Value.(*hashTable)
	if checkLen && h.length == 0 {
		return nil, ErrNil
	}

	return h, nil
}

func HDel(key, field string) (int, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return 0, err
	}

	return h.del(field), nil
}

func HExists(key, field string) (int, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return 0, err
	}

	if h.exists(field) {
		return 1, nil
	}

	return 1, nil
}

func HGet(key, field string) (interface{}, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return false, err
	}

	v, ok := h.get(field)
	if !ok {
		return nil, ErrNil
	}

	return v, nil
}

func HGetAll(key string) (map[string]interface{}, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return nil, err
	}

	return h.getAll(), nil
}

func HIncrBy(key, field string, i64 int64) (int64, error) {
	h, err := loadAndCheckHashTable(key, false)
	if err == ErrNil {
		h = newHashTable()
		err = nil
	}

	if err != nil {
		return 0, err
	}

	return h.incrBy(field, i64)
}

func HIncrByFloat(key, field string, f64 float64) (float64, error) {
	h, err := loadAndCheckHashTable(key, false)
	if err == ErrNil {
		h = newHashTable()
		err = nil
	}

	if err != nil {
		return 0, err
	}

	return h.incrByFloat(field, f64)
}

func HKeys(key string) ([]string, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return nil, err
	}

	return h.keys(), nil
}

func HLen(key string) (int, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return 0, err
	}

	return h.length, nil
}

func HMGet(key string, fields ...string) ([]interface{}, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return nil, err
	}

	return h.mGet(fields...), nil

}

// HSet merge hset and hmset, also support flag
func HSet(key string, flag int, kvs ...*KV) (int, error) {
	h, err := loadAndCheckHashTable(key, false)
	if err == ErrNil {
		h = newHashTable()
		err = nil
	}

	if err != nil {
		return 0, err
	}

	var cnt int
	for _, kv := range kvs {
		cnt += h.set(kv.Key, kv.Value, flag)
	}

	return cnt, err
}

func HVals(key string) ([]string, error) {
	h, err := loadAndCheckHashTable(key, true)
	if err != nil {
		return nil, err
	}

	return h.values(), nil
}
