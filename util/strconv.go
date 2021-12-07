package util

import (
	"strconv"
)

// ConvertToValidValue convert s either int64 or float64 otherwise return s itself.
func ConvertToValidValue(s string) interface{} {
	// try int64
	i64, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i64
	}

	f64, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f64
	}

	return s
}
