### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  500000	      2206 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	   10000	    262154 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     500	   3847292 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      30	  36191977 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  72849948 ns/op
    BenchmarkSimpleRTree_LoadPooled/10-4         	 1000000	      1444 ns/op
    BenchmarkSimpleRTree_LoadPooled/1000-4       	   10000	    207618 ns/op
    BenchmarkSimpleRTree_LoadPooled/10000-4      	    1000	   2094395 ns/op
    BenchmarkSimpleRTree_LoadPooled/100000-4     	     100	  23940990 ns/op
    BenchmarkSimpleRTree_LoadPooled/200000-4     	      30	  49239794 ns/op

### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	 1000000	      2454 ns/op	    1853 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	   10000	    385690 ns/op	   42966 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   5489529 ns/op	  403803 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  55128015 ns/op	 4008522 B/op	       8 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  94483116 ns/op	 8006218 B/op	       8 allocs/op
    BenchmarkSimpleRTree_LoadPooled/10-4         	 1000000	      1379 ns/op	     240 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/1000-4       	   10000	    172062 ns/op	     240 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/10000-4      	     500	   2078616 ns/op	     243 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/100000-4     	      50	  23456011 ns/op	     275 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/200000-4     	      30	  48702382 ns/op	     299 B/op	       2 allocs/op

## Benchmark Compute distances

    Benchmark_ComputeDistances-4         	50000000	        23.9 ns/op
    Benchmark_VectorComputeDistances-4   	100000000	        12.8 ns/op
