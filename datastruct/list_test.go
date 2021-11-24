package datastruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_print(t *testing.T) {
	var (
		values = []string{"1", "2", "3", "4"}
		list   = newListByRPush(values...)
	)

	list.print()
}

func TestList_LPop(t *testing.T) {
	var (
		values = []string{"1", "2", "3", "4"}
		list   = newListByRPush(values...)
	)

	for i, value := range values {
		val, ok := list.LPop()
		if !ok {
			assert.FailNow(t, "LPop return false", i)
			return
		}

		assert.Equal(t, value, val)
	}

	_, ok := list.LPop()
	assert.Equal(t, false, ok)
	_, ok = list.RPop()
	assert.Equal(t, false, ok)
}

func TestList_RPop(t *testing.T) {
	var (
		values = []string{"1", "2", "3", "4"}
		list   = newListByLPush(values...)
	)

	for i, value := range values {
		val, ok := list.RPop()
		if !ok {
			assert.FailNow(t, "LPop return false", i)
			return
		}

		assert.Equal(t, value, val)
	}

	_, ok := list.RPop()
	assert.Equal(t, false, ok)
	_, ok = list.LPop()
	assert.Equal(t, false, ok)
}
