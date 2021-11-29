package datastruct

import (
	"testing"
)

func Test_zSkipList_insert(t *testing.T) {
	zsl := newZSkipList()
	zsl.insert(10, "hello")
	zsl.insert(5, "world")
	zsl.insert(12, "golang")
	zsl.insert(20, "clang")
	zsl.insert(2, "java")
	zsl.insert(8, "javascript")
	zsl.print()
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
