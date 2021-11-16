package protocal

type RedisCommand struct {
	Command string
	Key     string
	Values  []string
	Options []string
}

// sadd key1 value1
// hadd key1 hkey hvalue
// zadd key1 value1 score
