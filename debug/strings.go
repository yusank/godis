package debug

import (
	"strings"
)

func Escape(str string) string {
	return strings.Replace(str, "\r\n", "\\r\\n", -1)
}
