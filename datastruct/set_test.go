package datastruct

//
//func Benchmark_concurrence_map_sAdd(b *testing.B) {
//	s := newSet()
//	for i := 0; i < b.N; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//}
//
//func Benchmark_concurrence_map_sIsMember(b *testing.B) {
//	s := newSet()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sIsMember(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_concurrence_map_sRem(b *testing.B) {
//	s := newSet()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sRem(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_concurrence_map_sPop(b *testing.B) {
//	s := newSet()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sPop(50)
//	}
//}
//
//func Benchmark_sync_map_sAdd(b *testing.B) {
//	s := newSetSyncMap()
//	for i := 0; i < b.N; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_sync_map_sIsMember(b *testing.B) {
//	s := newSetSyncMap()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sIsMember(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_sync_map_sRem(b *testing.B) {
//	s := newSetSyncMap()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sRem(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_non_lock_map_sAdd(b *testing.B) {
//	s := newSetNonLockMap()
//	for i := 0; i < b.N; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_non_lock_map_sIsMember(b *testing.B) {
//	s := newSetNonLockMap()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sIsMember(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_non_lock_map_sRem(b *testing.B) {
//	s := newSetNonLockMap()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sRem(strconv.Itoa(i))
//	}
//}
//
//func Benchmark_non_lock_map_sPop(b *testing.B) {
//	s := newSetNonLockMap()
//	for i := 0; i < 50_000; i++ {
//		s.sAdd(strconv.Itoa(i))
//	}
//
//	for i := 0; i < b.N; i++ {
//		s.sPop(50)
//	}
//}

/*
 * goos: darwin
 * goarch: amd64
 * pkg: github.com/yusank/godis/datastruct
 * cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
 * Benchmark_concurrence_map_sAdd
 * Benchmark_concurrence_map_sAdd-12         	 2619080	       440.8 ns/op
 * Benchmark_concurrence_map_sIsMember
 * Benchmark_concurrence_map_sIsMember-12    	13764466	        77.68 ns/op
 * Benchmark_concurrence_map_sRem
 * Benchmark_concurrence_map_sRem-12         	16740207	        65.18 ns/op
 * Benchmark_concurrence_map_sPop
 * Benchmark_concurrence_map_sPop-12         	     366	   2904074 ns/op
 * Benchmark_sync_map_sAdd
 * Benchmark_sync_map_sAdd-12                	 2101056	       765.1 ns/op
 * Benchmark_sync_map_sIsMember
 * Benchmark_sync_map_sIsMember-12           	15998791	        73.47 ns/op
 * Benchmark_sync_map_sRem
 * Benchmark_sync_map_sRem-12                	15768998	        76.62 ns/op
 * Benchmark_non_lock_map_sAdd
 * Benchmark_non_lock_map_sAdd-12            	 3233144	       359.6 ns/op
 * Benchmark_non_lock_map_sIsMember
 * Benchmark_non_lock_map_sIsMember-12       	15947702	        64.70 ns/op
 * Benchmark_non_lock_map_sRem
 * Benchmark_non_lock_map_sRem-12            	34205215	        34.53 ns/op
 * Benchmark_non_lock_map_sPop
 * Benchmark_non_lock_map_sPop-12            	 9022405	       136.1 ns/op
 * PASS
 */
