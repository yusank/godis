package util

import (
	"fmt"
	"testing"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/yusank/godis/util
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStringConcat
BenchmarkStringConcat-12    	14886952	        70.91 ns/op
PASS
*/
func BenchmarkStringConcat(b *testing.B) {
	var slice = []string{
		"$",
		"10",
		"\r\n",
		"12345123456",
		"\r\n",
	}

	for i := 0; i < b.N; i++ {
		StringConcat(16, slice...)
	}
}

/*
goos: darwin
goarch: amd64
pkg: github.com/yusank/godis/util
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStringConcat2
BenchmarkStringConcat2-12    	11279284	       107.4 ns/op
PASS
*/
func BenchmarkStringConcat2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("$%d\r\n%s\r\n", 10, "1234512345")
	}
}
