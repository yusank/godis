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
	ErrNil = errors.New("redis: not found")
)
