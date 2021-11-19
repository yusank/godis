package protocol

import (
	"fmt"
	"strings"
)

/*
 * put all messages and elements encode funtions
 */

func EncodeDataWithError(err error) string {
	return fmt.Sprintf("-%s\r\n", err.Error())
}

func EncodeDataWithBulkString(str string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
}

func EncodeDataWithNilString(_ string) string {
	return "$-1\r\n"
}

func EncodeDataWithSimpleString(str string) string {
	return "+" + str + "\r\n"
}

func EncodeDataWithInteger(str string) string {
	return fmt.Sprintf(":%s\r\n", str)
}

func EncodeDataWithArray(encodeData ...string) string {
	if len(encodeData) == 0 {
		return ""
	}

	sb := new(strings.Builder)
	_, _ = sb.WriteString(fmt.Sprintf("*%d\r\n", len(encodeData)))
	for _, str := range encodeData {
		_, _ = sb.WriteString(str)
	}

	return sb.String()
}
