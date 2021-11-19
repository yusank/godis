package redis

import "errors"

type KeyType string

const (
	KeyTypeString    KeyType = "string"
	KeyTypeList      KeyType = "list"
	KeyTypeSet       KeyType = "set"
	KeyTypeSortedSet KeyType = "sortedSet"
	KeyTypeHashMap   KeyType = "hashMap"
)

var (
	ErrUnknownCommand = errors.New("unknown command key")
)

var (
	RespOK      = "OK"
	RespPong    = "PONG"
	RespCommand = "COMMAND"
)
