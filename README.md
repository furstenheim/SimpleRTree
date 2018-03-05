### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 1000000	      1085 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  300000	      4834 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  200000	      7446 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  200000	     10594 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	     11179 ns/op


