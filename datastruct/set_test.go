package datastruct

import (
	"strconv"
	"testing"
)

func Benchmark_concurrence_map_sAdd(b *testing.B) {
	s := newSet()
	for i := 0; i < b.N; i++ {
		s.sAdd(strconv.Itoa(i))
	}

}

func Benchmark_concurrence_map_sIsMember(b *testing.B) {
	s := newSet()
	for i := 0; i < 50_000; i++ {
		s.sAdd(strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		s.sIsMember(strconv.Itoa(i))
	}
}

func Benchmark_concurrence_map_sRem(b *testing.B) {
	s := newSet()
	for i := 0; i < 50_000; i++ {
		s.sAdd(strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		s.sRem(strconv.Itoa(i))
	}
}

func Benchmark_sync_map_sAdd(b *testing.B) {
	s := newSetSyncMap()
	for i := 0; i < b.N; i++ {
		s.sAdd(strconv.Itoa(i))
	}
}

func Benchmark_sync_map_sIsMember(b *testing.B) {
	s := newSetSyncMap()
	for i := 0; i < 50_000; i++ {
		s.sAdd(strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		s.sIsMember(strconv.Itoa(i))
	}
}

func Benchmark_sync_map_sRem(b *testing.B) {
	s := newSetSyncMap()
	for i := 0; i < 50_000; i++ {
		s.sAdd(strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		s.sRem(strconv.Itoa(i))
	}
}

/*
 * goos: darwin
 * goarch: amd64
 * pkg: github.com/yusank/godis/datastruct
 * cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
 * Benchmark_concurrence_map_sAdd
 * Benchmark_concurrence_map_sAdd-12         	 2639160	       432.7 ns/op
 * Benchmark_concurrence_map_sIsMember
 * Benchmark_concurrence_map_sIsMember-12    	15346359	        75.65 ns/op
 * Benchmark_concurrence_map_sRem
 * Benchmark_concurrence_map_sRem-12         	16535848	        63.36 ns/op
 * Benchmark_sync_map_sAdd
 * Benchmark_sync_map_sAdd-12                	 1779470	       708.5 ns/op
 * Benchmark_sync_map_sIsMember
 * Benchmark_sync_map_sIsMember-12           	16335298	        71.55 ns/op
 * Benchmark_sync_map_sRem
 * Benchmark_sync_map_sRem-12                	16582731	        71.91 ns/op
 */
