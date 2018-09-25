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

    BenchmarkSimpleRTree_Load/10-4      	  300000	      8263 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    3000	    679020 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     300	   6251160 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  54439483 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  90571491 ns/op


### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  200000	      6546 ns/op	    2664 B/op	      53 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    5000	    422483 ns/op	  106152 B/op	     199 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     200	   5872441 ns/op	 1024456 B/op	    1239 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  51067505 ns/op	10156584 B/op	   11659 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	 101268226 ns/op	20309224 B/op	   23195 allocs/op




