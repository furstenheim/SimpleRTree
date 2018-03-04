### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 1000000	      1837 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  200000	      6048 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  200000	      8237 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  200000	     11066 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	     11945 ns/op


### Benchmark. Memory

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 1000000	      1825 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  200000	      6028 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  200000	      8126 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  100000	     11122 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	     12011 ns/op	       0 B/op	       0 allocs/op

