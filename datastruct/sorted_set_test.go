package datastruct

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_zSkipList_insert(t *testing.T) {
	zs := prepareZSetForTest()
	zs.zsl.print()
}

func prepareZSetForTest() *zSet {
	zs := newZSet()
	zs.zAdd(10, "hello", ZAddInNone)
	zs.zAdd(5, "world", ZAddInNone)
	zs.zAdd(12, "golang", ZAddInNone)
	zs.zAdd(20, "clang", ZAddInNone)
	zs.zAdd(2, "java", ZAddInNone)
	zs.zAdd(8, "javascript", ZAddInNone)
	zs.zAdd(1, "clang", ZAddInNone)

	return zs
}

/*
func Test_zslRandomLevel(t *testing.T) {
	var (
		static = make(map[int]int)
		cnt    = 100_000
	)

	for i := 0; i < cnt; i++ {
		static[zslRandomLevel()]++
	}

	for i := 0; i < ZSkipListMaxLevel; i++ {
		t.Log(i, static[i])
	}
}
*/

//func Benchmark_zslRandomLevel(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		zslRandomLevel()
//	}
//	/*
//	* goos: darwin
//	* goarch: amd64
//	* pkg: github.com/yusank/godis/datastruct
//	* cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
//	* Benchmark_zslRandomLevel
//	* Benchmark_zslRandomLevel-12    	43659864	        27.17 ns/op
//	 */
//}

func Test_zSkipList_rank(t *testing.T) {
	type args struct {
		score float64
		value string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "c",
			args: args{
				score: 1,
				value: "clang",
			},
			want: 1,
		},
		{
			name: "java",
			args: args{
				score: 2,
				value: "java",
			},
			want: 2,
		},
		{
			name: "w",
			args: args{
				score: 5,
				value: "world",
			},
			want: 3,
		},
		{
			name: "js",
			args: args{
				score: 8,
				value: "javascript",
			},
			want: 4,
		},
		{
			name: "h",
			args: args{
				score: 10,
				value: "hello",
			},
			want: 5,
		},
		{
			name: "go",
			args: args{
				score: 12,
				value: "golang",
			},
			want: 6,
		},
	}
	zs := prepareZSetForTest()
	zs.zsl.print()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := zs.zsl.rank(tt.args.score, tt.args.value)
			assert.Equal(t, tt.want, int(got))
		})
	}
}

func Test_zSkipList_count(t *testing.T) {

	type args struct {
		min     float64
		minOpen bool
		max     float64
		maxOpen bool
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "both_close",
			args: args{
				min:     0,
				minOpen: false,
				max:     1,
				maxOpen: false,
			},
			want: 1,
		},
		{
			name: "single_1",
			args: args{
				min:     1,
				minOpen: true,
				max:     2,
				maxOpen: false,
			},
			want: 1,
		},
		{
			name: "single_1",
			args: args{
				min:     1,
				minOpen: false,
				max:     2,
				maxOpen: true,
			},
			want: 1,
		},
		{
			name: "left_open",
			args: args{
				min:     1,
				minOpen: true,
				max:     10,
				maxOpen: false,
			},
			want: 4,
		},
		{
			name: "right_open",
			args: args{
				min:     1,
				minOpen: false,
				max:     7,
				maxOpen: true,
			},
			want: 3,
		},
		{
			name: "left_inf",
			args: args{
				min:     math.Inf(-1),
				minOpen: false,
				max:     10,
				maxOpen: false,
			},
			want: 5,
		},
		{
			name: "right_inf",
			args: args{
				min:     5,
				minOpen: false,
				max:     math.Inf(1),
				maxOpen: false,
			},
			want: 4,
		},
		{
			name: "both_inf",
			args: args{
				min:     math.Inf(-1),
				minOpen: false,
				max:     math.Inf(1),
				maxOpen: false,
			},
			want: 6,
		},
	}

	zs := prepareZSetForTest()
	zs.zsl.print()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := zs.zsl.count(tt.args.min, tt.args.minOpen, tt.args.max, tt.args.maxOpen)
			assert.Equal(t, int(tt.want), int(got))
		})
	}
}

func Test_zSkipList_zRange(t *testing.T) {
	type args struct {
		start      int
		stop       int
		withScores bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "without_score",
			args: args{
				start:      0,
				stop:       1,
				withScores: false,
			},
			want: []string{"clang", "java"},
		},
		{
			name: "with_score",
			args: args{
				start:      0,
				stop:       1,
				withScores: true,
			},
			want: []string{"clang", "1", "java", "2"},
		},
		{
			name: "head_to_tail",
			args: args{
				start:      0,
				stop:       -1,
				withScores: false,
			},
			want: []string{"clang", "java", "world", "javascript", "hello", "golang"},
		},
		{
			name: "reverse",
			args: args{
				start:      -2,
				stop:       -1,
				withScores: false,
			},
			want: []string{"hello", "golang"},
		},
		{
			name: "out_of_range",
			args: args{
				start:      7,
				stop:       10,
				withScores: false,
			},
			want: nil,
		},
		{
			name: "reverse_out_of_range",
			args: args{
				start:      -10,
				stop:       -9,
				withScores: false,
			},
			want: nil,
		},
	}

	zs := prepareZSetForTest()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := zs.zsl.zRange(tt.args.start, tt.args.stop, tt.args.withScores)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_zSkipList_zRangeByScore(t *testing.T) {

	type args struct {
		min        float64
		minOpen    bool
		max        float64
		maxOpen    bool
		withScores bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "both_close",
			args: args{
				min:        0,
				minOpen:    false,
				max:        1,
				maxOpen:    false,
				withScores: false,
			},
			want: []string{"clang"},
		},
		{
			name: "left_close",
			args: args{
				min:        0,
				minOpen:    false,
				max:        5,
				maxOpen:    true,
				withScores: false,
			},
			want: []string{"clang", "java"},
		},

		{
			name: "right_close",
			args: args{
				min:        1,
				minOpen:    true,
				max:        5,
				maxOpen:    false,
				withScores: false,
			},
			want: []string{"java", "world"},
		},
		{
			name: "both_open",
			args: args{
				min:        1,
				minOpen:    true,
				max:        5,
				maxOpen:    true,
				withScores: false,
			},
			want: []string{"java"},
		},
		{
			name: "left_inf",
			args: args{
				min:        math.Inf(-1),
				minOpen:    false,
				max:        5,
				maxOpen:    false,
				withScores: false,
			},
			want: []string{"clang", "java", "world"},
		},
		{
			name: "right_inf",
			args: args{
				min:        5,
				minOpen:    false,
				max:        math.Inf(1),
				maxOpen:    false,
				withScores: false,
			},
			want: []string{"world", "javascript", "hello", "golang"},
		},
		{
			name: "both_inf",
			args: args{
				min:        math.Inf(-1),
				minOpen:    false,
				max:        math.Inf(1),
				maxOpen:    false,
				withScores: false,
			},
			want: []string{"clang", "java", "world", "javascript", "hello", "golang"},
		},
	}
	zs := prepareZSetForTest()
	zs.zsl.print()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := zs.zsl.zRangeByScore(tt.args.min, tt.args.minOpen, tt.args.max, tt.args.maxOpen, tt.args.withScores)
			assert.Equal(t, tt.want, got)
		})
	}
}
