package protocol

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDataWithInteger(t *testing.T) {
	var want = ":12\r\n"
	assert.Equal(t, want, encodeInteger("12"))
}

func TestEncodeDataWithError(t *testing.T) {
	var want = "-EOF\r\n"
	assert.Equal(t, want, encodeError(io.EOF))
}

func TestEncodeDataWithBulkString(t *testing.T) {
	var want = "$2\r\nOK\r\n"
	assert.Equal(t, want, encodeBulkString("OK"))
}

func TestEncodeDataWithNilString(t *testing.T) {
	var want = "$-1\r\n"
	assert.Equal(t, want, encodeNilString(""))
}

func TestEncodeDataWithSimpleString(t *testing.T) {
	var want = "+OK\r\n"
	assert.Equal(t, want, encodeSimpleString("OK"))
}

func TestEncodeDataWithArray(t *testing.T) {
	var want = "*2\r\n$6\r\nHello \r\n$5\r\nWorld\r\n"
	assert.Equal(t, want, encodeBulkStrings(
		"Hello ",
		"World",
	))
}
