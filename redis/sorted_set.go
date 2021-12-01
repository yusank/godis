package redis

import (
	"math"
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
		minStr           = strings.ToLower(c.Values[1])
		maxStr           = strings.ToLower(c.Values[2])
		min, max         float64
		minOpen, maxOpen bool
	)

	switch minStr {
	case "-inf":
		min = math.Inf(-1)
	default:
		tmp := strings.TrimPrefix(minStr, "(")
		if tmp != minStr {
			minOpen = true
		}

		f, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			return nil, ErrValueOutOfRange
		}
		min = f
	}

	switch maxStr {
	case "+inf":
		max = math.Inf(1)
	default:
		tmp := strings.TrimPrefix(maxStr, "(")
		if tmp != maxStr {
			maxOpen = true
		}

		f, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			return nil, ErrValueOutOfRange
		}
		max = f
	}

	cnt, err := datastruct.ZCount(c.Values[0], min, minOpen, max, maxOpen)
	if err == datastruct.ErrNil {
		return protocol.NewResponseWithInteger(0), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewResponseWithInteger(int64(cnt)), nil
}
