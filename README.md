### Simple RTree

This is a static RTree that only accepts points as input. The format of the points is a single array where each too coordinates represent a point

    points := []float64{0.0, 0.0, 1.0, 1.0} // array of two points 0, 0 and 1, 1

The library exposes only two methods. One to load and one to find nearest point

    fp := FlatPoints(points)
    r := New().Load(fp)
    closestX, closestY, distance, found := r.FindNearestPoint(1.0, 2.0)
    // 1.0, 1.0, 1.0, true



### Benchmark. CPU

These are the benchmarks nearest point once the index is tree. There are 0 allocations in the heap

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 2000000	       702 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  500000	      3005 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  300000	      4794 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  200000	      7572 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	      7735 ns/op



### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  300000	      4965 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    5000	    400279 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   4047283 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  52048036 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  90448290 ns/op


### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  300000	      4950 ns/op	    3472 B/op	      56 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    5000	    391634 ns/op	  198832 B/op	     376 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   4138240 ns/op	 1942276 B/op	    3082 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  51294535 ns/op	19388006 B/op	   30356 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  92538965 ns/op	38764060 B/op	   60589 allocs/op

