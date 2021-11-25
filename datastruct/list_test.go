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

func TestList_LRange(t *testing.T) {
	var (
		values = []string{"1", "2", "3", "4", "5", "6"}
		list   = newListByRPush(values...)
	)

	tests := []struct {
		name        string
		start, stop int
		want        []string
	}{
		{
			name:  "range_same_index",
			start: 0,
			stop:  0,
			want:  []string{values[0]},
		},
		{
			name:  "range_part_of_slice",
			start: 1,
			stop:  3,
			want:  values[1:4],
		},
		{
			name:  "range_all_list",
			start: 0,
			stop:  len(values),
			want:  values,
		},
		{
			name:  "range_out_off_range",
			start: len(values) + 1,
			stop:  len(values) * 2,
			want:  nil,
		},
		{
			name:  "range_negative_index",
			start: -len(values),
			stop:  -1, // last one
			want:  values,
		},
		{
			name:  "range_negative_single_value",
			start: -len(values) * 2,
			stop:  -len(values), // last one
			want:  []string{values[0]},
		},
		{
			name:  "range_negative_out_of_range",
			start: -len(values) * 2,
			stop:  -len(values) - 1, // last one
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSlice := list.LRange(tt.start, tt.stop)
			if !assert.Equal(t, tt.want, gotSlice) {
				return
			}
		})
	}
}
