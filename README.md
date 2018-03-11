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

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 2000000	       689 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  500000	      3054 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  300000	      4717 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  200000	      7033 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  200000	      7560 ns/op


### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  200000	     10233 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    2000	    739989 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     200	   6100647 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  52645807 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 112139103 ns/op

### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  300000	      9708 ns/op	    2408 B/op	      46 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    2000	    733970 ns/op	  126040 B/op	    1464 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   5847426 ns/op	 1262392 B/op	   14309 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  59375040 ns/op	12629512 B/op	  143292 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 115049926 ns/op	25256472 B/op	  286519 allocs/op
