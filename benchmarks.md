### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	 1000000	      1060 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	   10000	    134926 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	    1000	   1624066 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	     100	  16894472 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      50	  35175553 ns/op
    BenchmarkSimpleRTree_LoadPooled/10-4         	 3000000	       433 ns/op
    BenchmarkSimpleRTree_LoadPooled/1000-4       	   10000	    135625 ns/op
    BenchmarkSimpleRTree_LoadPooled/10000-4      	    1000	   1462810 ns/op
    BenchmarkSimpleRTree_LoadPooled/100000-4     	     100	  15882229 ns/op
    BenchmarkSimpleRTree_LoadPooled/200000-4     	      50	  33064794 ns/op

### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	 1000000	      1075 ns/op	    1836 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	   10000	    136323 ns/op	   42950 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	    1000	   1628629 ns/op	  403787 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	     100	  16912108 ns/op	 4008506 B/op	       8 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      50	  35426528 ns/op	 8006202 B/op	       8 allocs/op
    BenchmarkSimpleRTree_LoadPooled/10-4         	 3000000	       432 ns/op	     240 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/1000-4       	   10000	    135548 ns/op	     240 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/10000-4      	    1000	   1462004 ns/op	     241 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/100000-4     	     100	  17138266 ns/op	     257 B/op	       2 allocs/op
    BenchmarkSimpleRTree_LoadPooled/200000-4     	      50	  34608146 ns/op	     275 B/op	       2 allocs/op


## Benchmark Hilbert tree
    BenchmarkSimpleRTree_FindNearestPointHilbert/10-4         	10000000	       179 ns/op
    BenchmarkSimpleRTree_FindNearestPointHilbert/1000-4       	  500000	      2603 ns/op
    BenchmarkSimpleRTree_FindNearestPointHilbert/10000-4      	  500000	      3685 ns/op
    BenchmarkSimpleRTree_FindNearestPointHilbert/100000-4     	  200000	      6138 ns/op
    BenchmarkSimpleRTree_FindNearestPointHilbert/200000-4     	  200000	      6819 ns/op
    BenchmarkSimpleRTree_FindNearestPointHilbert/1000000-4    	  200000	      6654 ns/op
    BenchmarkSimpleRTree_FindNearestPointHilbert/10000000-4   	  200000	     11072 ns/op

## Benchmark Compute distances

    Benchmark_ComputeDistances-4         	100000000	        20.2 ns/op
    Benchmark_VectorComputeDistances-4   	200000000	         8.27 ns/op