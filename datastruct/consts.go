package datastruct

import (
	"errors"
)

type KeyType string

const (
	KeyTypeString    KeyType = "string"
	KeyTypeList      KeyType = "list"
	KeyTypeSet       KeyType = "set"
	KeyTypeSortedSet KeyType = "sortedSet"
	KeyTypeHashMap   KeyType = "hashMap"
)

var (
	ErrNil                   = errors.New("redis: not found")
	ErrKeyAndCommandNotMatch = errors.New("key type and command not match")
	ErrNotInteger            = errors.New("value is not an integer or out of range")
)
