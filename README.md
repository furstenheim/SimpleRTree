### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  300000	      5582 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  100000	     19974 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	   50000	     35650 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     31729 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     33127 ns/op

### Benchmark. Memory

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  500000	      6201 ns/op	     203 B/op	       3 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	   50000	     24056 ns/op	     479 B/op	      11 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	   50000	     33384 ns/op	     655 B/op	      16 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     31372 ns/op	     808 B/op	      22 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     33221 ns/op	     840 B/op	      24 allocs/op
