## Simple RTree

Simple RTree is a blazingly fast and GC friendly RTree. It performs under 4 microseconds with 1 Million points for closest point queries
(measured in a i5-2450M CPU @ 2.50GHz with 4Gb of RAM). It is GC friendly, queries require 0 allocation.
Building the index requires 19 allocations independently of the number of points.

To achieve this speed, the index has three restrictions. It is static, once built it cannot be changed.
It only accepts points, no bboxes or lines. It only accepts (for now) closest point queries.


### Usage

The format of the points is a single array where each too coordinates represent a point

    points := []float64{0.0, 0.0, 1.0, 1.0} // array of two points 0, 0 and 1, 1

The library exposes only two methods. One to load and one to find nearest point

    import "SimpleRTree"

    fp := SimpleRTree.FlatPoints(points)
    r := SimpleRTree.New().Load(fp)
    closestX, closestY, distance, found := r.FindNearestPoint(1.0, 2.0)
    // 1.0, 1.0, 1.0, true



### Benchmark. CPU

These are the benchmarks for finding the nearest point once the index has been built. There are 0 allocations in the heap

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 3000000	       466 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 1000000	      1333 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	 1000000	      1881 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	  500000	      2624 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	  500000	      2850 ns/op

### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  300000	      3591 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	   10000	    404584 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     500	   6526718 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      20	  58900581 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  77797865 ns/op


### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	  500000	      3160 ns/op	    2744 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	    5000	    390040 ns/op	  101176 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     500	   3635069 ns/op	  970296 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      30	  40846081 ns/op	 9605176 B/op	      19 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  80569543 ns/op	19206200 B/op	      19 allocs/op


## Benchmark Compute distances

    Benchmark_ComputeDistances-4         	100000000	        24.0 ns/op
    Benchmark_VectorComputeDistances-4   	100000000	        19.0 ns/op



