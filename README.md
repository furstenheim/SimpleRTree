### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  500000	      3661 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  100000	     17041 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	   50000	     24730 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     33136 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     35200 ns/op



### Benchmark. Memory

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  300000	      3727 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  100000	     16920 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	   50000	     24665 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     32715 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     34629 ns/op	       0 B/op	       0 allocs/op

