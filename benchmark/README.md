# Benchmark

``` shell
 * goos: darwin
 * goarch: amd64
 * pkg: github.com/yusank/godis/datastruct
 * cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

$ redis-benchmark -p 7379 -n 10000 -q --csv
```

- 2021/12/08 first benchmark [result](benchmark.20211208.csv).