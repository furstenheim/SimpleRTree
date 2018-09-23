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

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 3000000	       549 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 1000000	      1812 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  500000	      2909 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  500000	      3983 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  300000	      4250 ns/op




### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  200000	      9070 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    3000	    711563 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     200	   6888222 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  55542970 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 105043234 ns/op


### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  200000	      9470 ns/op	    3040 B/op	      57 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    3000	    715728 ns/op	  149760 B/op	     377 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     200	   7487503 ns/op	 1467229 B/op	    3083 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  62603439 ns/op	14587572 B/op	   30357 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 111886749 ns/op	29163014 B/op	   60590 allocs/op


