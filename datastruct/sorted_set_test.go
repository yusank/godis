package datastruct

import (
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

func Benchmark_zslRandomLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		zslRandomLevel()
	}
	/*
	* goos: darwin
	* goarch: amd64
	* pkg: github.com/yusank/godis/datastruct
	* cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
	* Benchmark_zslRandomLevel
	* Benchmark_zslRandomLevel-12    	43659864	        27.17 ns/op
	 */
}

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
