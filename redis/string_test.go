package redis

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseInt(t *testing.T) {
	v, err := strconv.ParseInt("12", 10, 64)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, int64(12), v)
}
