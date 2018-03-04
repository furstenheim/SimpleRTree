### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  500000	      3644 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  100000	     17649 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  100000	     27525 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     33100 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     32779 ns/op


### Benchmark. Memory

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  500000	      3840 ns/op	      13 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  100000	     16977 ns/op	     144 B/op	       9 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	   50000	     30359 ns/op	     223 B/op	      13 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     31566 ns/op	     312 B/op	      19 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     32787 ns/op	     343 B/op	      21 allocs/op
