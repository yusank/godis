package datastruct

import (
	"errors"
)

type KeyType string

const (
	KeyTypeString    KeyType = "string"
	KeyTypeList      KeyType = "list"
	KeyTypeSet       KeyType = "setIfNotExists"
	KeyTypeSortedSet KeyType = "sortedSet"
	KeyTypeHashTable KeyType = "hashTable"
)

var (
	ErrNil                   = errors.New("redis: not found")
	ErrKeyAndCommandNotMatch = errors.New("key type and command not match")
	ErrNotInteger            = errors.New("value is not an integer or out of range")
	ErrNotFloat              = errors.New("ERR value is not a valid float")
)

// zAdd flags
const (
	ZAddInNone = 0
	ZAddInIncr = 1 << (iota - 1)
	ZAddInNx
	ZAddInXx
)

const (
	ZRangeInNone       = 0
	ZRangeInWithScores = 1 << (iota - 1)
	ZRangeInByScore
	ZRangeInByLex
)

const (
	HSetInNone = 0
	HSetInNx   = 1 << (iota - 1)
)
