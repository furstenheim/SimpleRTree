### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4         	  300000	      6938 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4       	  100000	     22498 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4      	   50000	     37416 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4     	   50000	     32329 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4     	   50000	     32603 ns/op

### Benchmark. Memory

    BenchmarkSimpleRTree_FindNearestPoint/10-4         	  500000	      6112 ns/op	     254 B/op	       6 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4       	  100000	     26322 ns/op	     571 B/op	      17 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4      	   50000	     28899 ns/op	     763 B/op	      23 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4     	   50000	     35960 ns/op	     932 B/op	      30 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4     	   30000	     34685 ns/op	     963 B/op	      32 allocs/op
