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

// zRank .
func zRank(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	rank, err := datastruct.ZRank(c.Values[0], c.Values[1])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}

	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(rank)), nil
}

// zRem .
func zRem(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 2 {
		return nil, ErrCommandArgsNotEnough
	}

	cnt, err := datastruct.ZRem(c.Values[0], c.Values[1:]...)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithInteger(0), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// zCard .
func zCard(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 1 {
		return nil, ErrCommandArgsNotEnough
	}

	cnt, err := datastruct.ZCard(c.Values[0])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithInteger(0), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// zCount .
func zCount(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key    = c.Values[0]
		minStr = strings.ToLower(c.Values[1])
		maxStr = strings.ToLower(c.Values[2])
	)

	cnt, err := datastruct.ZCount(key, minStr, maxStr)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithInteger(0), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}

// zIncr .
func zIncr(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	score, err := strconv.ParseFloat(c.Values[1], 64)
	if err != nil {
		return nil, ErrValueOutOfRange
	}

	curScore, err := datastruct.ZIncr(c.Values[0], score, c.Values[2])
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithInteger(0), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithBulkString(strconv.FormatFloat(curScore, 'g', -1, 64)), nil
}

// zRange .
func zRange(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key    = c.Values[0]
		minStr = c.Values[1]
		maxStr = c.Values[2]
		flag   = datastruct.ZRangeInNone
	)

	for _, s := range c.Values[2:] {
		switch strings.ToLower(s) {
		case "withscores":
			flag |= datastruct.ZRangeInWithScores
		case "byscore":
			flag |= datastruct.ZRangeInByScore

		}
	}

	values, err := datastruct.ZRange(key, minStr, maxStr, flag)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(values...), nil
}

// zRangeByScore .
func zRangeByScore(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	values, err := datastruct.ZRange(c.Values[0], c.Values[1], c.Values[2], datastruct.ZRangeInByScore)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(values...), nil
}

// zRevRange .
func zRevRange(c *Command) (*protocol.Response, error) {
	if len(c.Values) < 3 {
		return nil, ErrCommandArgsNotEnough
	}

	var (
		key    = c.Values[0]
		minStr = c.Values[1]
		maxStr = c.Values[2]
		flag   = datastruct.ZRangeInNone
	)

	for _, s := range c.Values[2:] {
		switch strings.ToLower(s) {
		case "withscores":
			flag |= datastruct.ZRangeInWithScores
		}
	}

	values, err := datastruct.ZRevRange(key, minStr, maxStr, flag)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithNilBulk(), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponse(true).AppendBulkStrings(values...), nil

}
