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

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 3000000	       516 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 1000000	      1805 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  500000	      2677 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  500000	      3561 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  500000	      3803 ns/op



### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  500000	      3528 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    5000	    399319 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   3934939 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      30	  37848676 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  67023164 ns/op



### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  500000	      3160 ns/op	    2744 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    5000	    390040 ns/op	  101176 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     500	   3635069 ns/op	  970296 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      30	  40846081 ns/op	 9605176 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  80569543 ns/op	19206200 B/op	      19 allocs/op




