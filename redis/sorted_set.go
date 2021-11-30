package redis

import (
	"strconv"
	"strings"

	"github.com/yusank/godis/datastruct"
	"github.com/yusank/godis/protocol"
)

// zAdd .
func zAdd(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key     = c.Values[0]
		members = make([]*datastruct.ZSetMember, 0)
		flag    int
	)

	for i := 1; i < len(c.Values)-1; {
		switch strings.ToLower(c.Values[i]) {
		case "xx":
			flag |= datastruct.ZAddInXx
			i++
		case "nx":
			flag |= datastruct.ZAddInNx
			i++
		case "incr":
			flag |= datastruct.ZAddInIncr
			i++
		// add more
		default:
			score, err := strconv.ParseFloat(c.Values[i], 64)
			if err != nil {
				return nil, ErrValueOutOfRange
			}

			members = append(members, &datastruct.ZSetMember{
				Score: score,
				Value: c.Values[i+1],
			})
			i += 2
		}
	}

	cnt, err := datastruct.ZAdd(key, members, flag)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// zScore .
func zScore(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	f, err := datastruct.ZScore(c.Values[0], c.Values[1])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithBulkString(strconv.FormatFloat(f, 'g', -1, 64)), nil
}
