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

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 3000000	       571 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 1000000	      1894 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  500000	      2723 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  300000	      3796 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  300000	      3943 ns/op




### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  300000	      6020 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	    3000	    442769 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     500	   3809178 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  58292293 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  82287931 ns/op




### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  200000	      8695 ns/op	    2792 B/op	      53 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    3000	    748523 ns/op	  122536 B/op	     199 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     200	   6377959 ns/op	 1180104 B/op	    1239 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      30	  79146203 ns/op	11762216 B/op	   11659 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      10	 123593821 ns/op	23512296 B/op	   23195 allocs/op



