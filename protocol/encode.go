package protocol

import (
	"fmt"
	"strings"
)

/*
 * put all messages and elements encode funtions
 */

func encodeError(err error) string {
	return fmt.Sprintf("-%s\r\n", err.Error())
}

func encodeBulkString(str string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
}

func encodeNilString() string {
	return "$-1\r\n"
}

func encodeSimpleString(str string) string {
	return "+" + str + "\r\n"
}

func encodeInteger(str string) string {
	return fmt.Sprintf(":%s\r\n", str)
}

func encodeBulkStrings(strSlice ...string) string {
	if len(strSlice) == 0 {
		return encodeNilString()
	}

	sb := new(strings.Builder)
	_, _ = sb.WriteString(fmt.Sprintf("*%d\r\n", len(strSlice)))
	for _, str := range strSlice {
		_, _ = sb.WriteString(encodeBulkString(str))
	}

	return sb.String()
}
