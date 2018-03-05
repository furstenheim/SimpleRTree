### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 2000000	       754 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  500000	      3116 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  300000	      4743 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  200000	      6988 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	      7828 ns/op


### Benchmark Load CPU
    BenchmarkSimpleRTree_Load/10-4      	  200000	     10233 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    2000	    739989 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     200	   6100647 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  52645807 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 112139103 ns/op

### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  300000	      9708 ns/op	    2408 B/op	      46 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    2000	    733970 ns/op	  126040 B/op	    1464 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   5847426 ns/op	 1262392 B/op	   14309 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  59375040 ns/op	12629512 B/op	  143292 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 115049926 ns/op	25256472 B/op	  286519 allocs/op
