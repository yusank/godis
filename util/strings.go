package util

import (
	"strings"
)

func StringConcat(n int, strs ...string) string {
	sb := new(strings.Builder)
	sb.Grow(n)
	for _, str := range strs {
		sb.WriteString(str)
	}

	return sb.String()
}
