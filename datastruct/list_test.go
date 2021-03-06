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
			name:  "range_all",
			start: 0,
			stop:  -1,
			want:  values,
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

func TestList_LRemCountFromHead(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		remValue string
		remCnt   int
		wantNum  int
		wantRes  []string
	}{
		{
			name:     "rem_head",
			values:   []string{"1", "2", "3", "4", "5"},
			remValue: "1",
			remCnt:   1,
			wantNum:  1,
			wantRes:  []string{"2", "3", "4", "5"},
		},
		{
			name:     "rem_tail",
			values:   []string{"1", "2", "3", "4", "5"},
			remValue: "5",
			remCnt:   1,
			wantNum:  1,
			wantRes:  []string{"1", "2", "3", "4"},
		},
		{
			name:     "rem_from_one_element_list",
			values:   []string{"1"},
			remValue: "1",
			remCnt:   1,
			wantNum:  1,
			wantRes:  nil,
		},
		{
			name:     "rem_from_middle",
			values:   []string{"1", "2", "3", "4", "5"},
			remValue: "2",
			remCnt:   1,
			wantNum:  1,
			wantRes:  []string{"1", "3", "4", "5"},
		},
		{
			name:     "rem_from_middle_continuation",
			values:   []string{"1", "2", "2", "3", "4", "5"},
			remValue: "2",
			remCnt:   2,
			wantNum:  2,
			wantRes:  []string{"1", "3", "4", "5"},
		},
		{
			name:     "rem_from_middle_not_continuation",
			values:   []string{"1", "2", "3", "2", "3", "4", "5"},
			remValue: "2",
			remCnt:   2,
			wantNum:  2,
			wantRes:  []string{"1", "3", "3", "4", "5"},
		},
		{
			name:     "rem_from_middle_not_continuation",
			values:   []string{"1", "1", "3", "1"},
			remValue: "1",
			remCnt:   2,
			wantNum:  2,
			wantRes:  []string{"3", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := newListByRPush(tt.values...)
			cnt := list.LRemCountFromHead(tt.remValue, tt.remCnt)
			if !assert.Equal(t, tt.wantNum, cnt) {
				return
			}

			assert.Equal(t, tt.wantRes, list.LRange(0, -1))
		})
	}
}

func TestList_LRemCountFromTail(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		remValue string
		remCnt   int
		wantNum  int
		wantRes  []string
	}{
		{
			name:     "rem_tail",
			values:   []string{"1", "2", "3", "4", "5"},
			remValue: "5",
			remCnt:   1,
			wantNum:  1,
			wantRes:  []string{"1", "2", "3", "4"},
		},
		{
			name:     "rem_head",
			values:   []string{"1", "2", "3", "4", "5"},
			remValue: "1",
			remCnt:   1,
			wantNum:  1,
			wantRes:  []string{"2", "3", "4", "5"},
		},
		{
			name:     "rem_from_one_element_list",
			values:   []string{"1"},
			remValue: "1",
			remCnt:   1,
			wantNum:  1,
			wantRes:  nil,
		},
		{
			name:     "rem_from_middle",
			values:   []string{"1", "2", "3", "4", "5"},
			remValue: "2",
			remCnt:   1,
			wantNum:  1,
			wantRes:  []string{"1", "3", "4", "5"},
		},
		{
			name:     "rem_from_middle_continuation",
			values:   []string{"1", "2", "2", "3", "4", "5"},
			remValue: "2",
			remCnt:   2,
			wantNum:  2,
			wantRes:  []string{"1", "3", "4", "5"},
		},
		{
			name:     "rem_from_middle_not_continuation",
			values:   []string{"1", "2", "3", "2", "3", "4", "5"},
			remValue: "2",
			remCnt:   2,
			wantNum:  2,
			wantRes:  []string{"1", "3", "3", "4", "5"},
		},
		{
			name:     "rem_from_middle_not_continuation",
			values:   []string{"1", "1", "3", "1"},
			remValue: "1",
			remCnt:   2,
			wantNum:  2,
			wantRes:  []string{"1", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := newListByRPush(tt.values...)
			cnt := list.LRemCountFromTail(tt.remValue, tt.remCnt)
			if !assert.Equal(t, tt.wantNum, cnt) {
				return
			}

			assert.Equal(t, tt.wantRes, list.LRange(0, -1))
		})
	}
}

func TestList_LIndex(t *testing.T) {
	tests := []struct {
		name   string
		index  int
		values []string
		want   string
		want1  bool
	}{
		{
			name:   "head",
			index:  0,
			values: []string{"1", "2", "3", "4"},
			want:   "1",
			want1:  true,
		},
		{
			name:   "tail",
			index:  -1,
			values: []string{"1", "2", "3", "4"},
			want:   "4",
			want1:  true,
		},
		{
			name:   "mid",
			index:  1,
			values: []string{"1", "2", "3", "4"},
			want:   "2",
			want1:  true,
		},
		{
			name:   "invalid_1",
			index:  5,
			values: []string{"1", "2", "3", "4"},
			want:   "",
			want1:  false,
		},
		{
			name:   "invalid_2",
			index:  -5,
			values: []string{"1", "2", "3", "4"},
			want:   "",
			want1:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newListByRPush(tt.values...)
			got, got1 := l.LIndex(tt.index)
			if !assert.Equal(t, tt.want, got) {
				return
			}

			if !assert.Equal(t, tt.want1, got1) {
				return
			}
		})
	}
}

func TestList_LSet(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		index  int
		value  string
		want   bool
	}{
		{
			name:   "head",
			values: []string{"1", "2", "3"},
			index:  0,
			value:  "2",
			want:   true,
		},
		{
			name:   "tail",
			values: []string{"1", "2", "3"},
			index:  -1,
			value:  "2",
			want:   true,
		},
		{
			name:   "mid",
			values: []string{"1", "2", "3"},
			index:  1,
			value:  "1",
			want:   true,
		},
		{
			name:   "invalid",
			values: []string{"1", "2", "3"},
			index:  4,
			value:  "2",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newListByRPush(tt.values...)
			got := l.LSet(tt.index, tt.value)
			if !assert.Equal(t, tt.want, got) || !got {
				return
			}
			got1, _ := l.LIndex(tt.index)
			assert.Equal(t, tt.value, got1)
		})
	}
}

func TestList_LInsert(t *testing.T) {
	type args struct {
		target   string
		newValue string
		flag     int
	}
	tests := []struct {
		name   string
		values []string
		args   args
		want   bool
	}{
		{
			name:   "head",
			values: []string{"2", "2", "3", "4"},
			args: args{
				target:   "2",
				newValue: "1",
				flag:     0,
			},
			want: true,
		},
		{
			name:   "tail",
			values: []string{"1", "2", "3", "4"},
			args: args{
				target:   "4",
				newValue: "5",
				flag:     0,
			},
			want: true,
		},
		{
			name:   "after",
			values: []string{"1", "3", "3", "4"},
			args: args{
				target:   "1",
				newValue: "2",
				flag:     1,
			},
			want: true,
		},
		{
			name:   "before",
			values: []string{"1", "2", "2", "4"},
			args: args{
				target:   "4",
				newValue: "3",
				flag:     -1,
			},
			want: true,
		},
		{
			name:   "not_found",
			values: []string{"1", "2", "3", "4"},
			args: args{
				target:   "-1",
				newValue: "1",
				flag:     -1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newListByRPush(tt.values...)
			got := l.LInsert(tt.args.target, tt.args.newValue, tt.args.flag)
			assert.Equal(t, tt.want, got)
			l.print()
		})
	}
}
