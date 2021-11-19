package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readBulkOrArrayLength(t *testing.T) {
	var (
		got_1 = 12
		got_2 = 3
	)

	assert.Equal(t, got_1, readBulkOrArrayLength([]byte("$12\r\n")))
	assert.Equal(t, got_2, readBulkOrArrayLength([]byte("*3\r\n")))
}

func Test_readBulkStrings(t *testing.T) {
	var (
		want = []byte("abcdefg")
		data = append(want, []byte(CRLF)...)
	)
	buf := bytes.NewBuffer(data)
	r := bufio.NewReader(buf)

	gotVal, err := readBulkStrings(r, len(want))
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, want, gotVal)
}

func Test_readArray(t *testing.T) {
	var (
		str  = "$2\r\nok\r\n+ping\r\n:12\r\n"
		data = []byte(str)
	)
	buf := bytes.NewBuffer(data)
	r := bufio.NewReader(buf)

	gotVal, err := readArray(r, 3)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Len(t, gotVal, 3) {
		return
	}

	assert.Equal(t, gotVal[0].Value, "ok")
	assert.Equal(t, gotVal[0].Description, DescriptionBulkStrings)
	assert.Equal(t, gotVal[1].Value, "ping")
	assert.Equal(t, gotVal[1].Description, DescriptionSimpleStrings)
	assert.Equal(t, gotVal[2].Value, "12")
	assert.Equal(t, gotVal[2].Description, DescriptionIntegers)
}
