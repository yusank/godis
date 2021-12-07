package datastruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hashTable_set(t *testing.T) {
	type args struct {
		field string
		v     interface{}
		flag  int
	}
	tests := []struct {
		name    string
		preData []*KV // pre insert data
		args    args
		want    int
	}{
		{
			name: "no_flag",
			args: args{
				field: "f1",
				v:     "a",
				flag:  0,
			},
			want: 1,
		},
		{
			name: "no_flag_exists",
			args: args{
				field: "f1",
				v:     "b",
				flag:  0,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: "a",
				},
			},
			want: 1,
		},
		{
			name: "with_flag_no_pre_data",
			args: args{
				field: "f1",
				v:     "a",
				flag:  HSetInNx,
			},
			want: 1,
		},
		{
			name: "with_flag_pre_data",
			args: args{
				field: "f1",
				v:     "b",
				flag:  HSetInNx,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: "a",
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHashTable()
			h.mSet(tt.preData)
			got := h.set(tt.args.field, tt.args.v, tt.args.flag)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_hashTable_incrBy(t *testing.T) {
	type args struct {
		field string
		i     int64
	}
	tests := []struct {
		name    string
		args    args
		preData []*KV // pre insert data
		want    int64
		wantErr bool
	}{
		{
			name: "not exists",
			args: args{
				field: "f1",
				i:     10,
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "exists data",
			args: args{
				field: "f1",
				i:     10,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: int64(10),
				},
			},
			want:    20,
			wantErr: false,
		},
		{
			name: "incr negative number",
			args: args{
				field: "f1",
				i:     -15,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: int64(10),
				},
			},
			want:    -5,
			wantErr: false,
		},
		{
			name: "invalid string data",
			args: args{
				field: "f1",
				i:     10,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: "10",
				},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid float data",
			args: args{
				field: "f1",
				i:     10,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: float64(10),
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHashTable()
			h.mSet(tt.preData)
			got, err := h.incrBy(tt.args.field, tt.args.i)
			if tt.wantErr {
				if !assert.Error(t, err) {
					return
				}

				return
			}

			assert.Equal(t, tt.want, got)

		})
	}
}

func Test_hashTable_incrByFloat(t *testing.T) {
	type args struct {
		field string
		i     float64
	}
	tests := []struct {
		name    string
		args    args
		preData []*KV // pre insert data
		want    float64
		wantErr bool
	}{
		{
			name: "not exists",
			args: args{
				field: "f1",
				i:     1.2,
			},
			want:    1.2,
			wantErr: false,
		},
		{
			name: "exists data",
			args: args{
				field: "f1",
				i:     1.2,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: 1.3,
				},
			},
			want:    2.5,
			wantErr: false,
		},
		{
			name: "incr negative float",
			args: args{
				field: "f1",
				i:     -1.5,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: 0.9,
				},
			},
			want:    -0.6,
			wantErr: false,
		},
		{
			name: "invalid string data",
			args: args{
				field: "f1",
				i:     10,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: "10",
				},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid int data",
			args: args{
				field: "f1",
				i:     10,
			},
			preData: []*KV{
				{
					Key:   "f1",
					Value: int64(10),
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHashTable()
			h.mSet(tt.preData)
			got, err := h.incrByFloat(tt.args.field, tt.args.i)
			if tt.wantErr {
				if !assert.Error(t, err) {
					return
				}

				return
			}

			assert.Equal(t, tt.want, got)

		})
	}
}
