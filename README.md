### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 1000000	      1085 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  300000	      4353 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  200000	      6586 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  200000	      9607 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	     10289 ns/op



### Benchmark Load CPU
    BenchmarkSimpleRTree_Load/10-4      	  300000	      9547 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    3000	    865619 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   6539675 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  69242877 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 101854519 ns/op
