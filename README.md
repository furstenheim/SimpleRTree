### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 1000000	      1112 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  300000	      4846 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  200000	      7415 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  200000	     10258 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	     11079 ns/op



