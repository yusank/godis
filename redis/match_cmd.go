package redis

import (
	"github.com/yusank/godis/protocol"
)

// ExecuteFunc define command execute, returns slice of string as result and error if has any error occur .
type ExecuteFunc func(*Command) (*protocol.Response, error)

var knownCommands = map[string]ExecuteFunc{
	// native
	"command": func(c *Command) (*protocol.Response, error) {
		return protocol.NewResponseWithSimpleString(RespCommand), nil
	},
	"ping": func(c *Command) (*protocol.Response, error) {
		return protocol.NewResponseWithSimpleString(RespPong), nil
	},
	"keys":   keys,
	"exists": exists,
	"type":   keyType,
	// strings
	"append": stringAppend,
	"incr":   incr,
	"incrby": incrBy,
	"decr":   decr,
	"decrby": decrBy,
	"get":    get,
	"mget":   mget,
	"set":    set,
	// list
	"lpush":   lPush,
	"lpop":    lPop,
	"llen":    lLen,
	"rpush":   rPush,
	"rpop":    rPop,
	"lrange":  lRange,
	"lrem":    lRem,
	"lindex":  lInsert,
	"lset":    lSet,
	"linsert": lInsert,
	//... more
	// zset
	"zadd":   zAdd,
	"zscore": zScore,
	"zrank":  zRank,
	"zrem":   zRem,
	"zcard":  zCard,
	"zcount": zCount,
	"zincr":  zIncr,
	"zrange": zRange,
}
