## Simple RTree

Simple RTree is a blazingly fast and GC friendly RTree. It performs in 2.36 microseconds with 1 Million points for closest point queries
(measured in a i5-2450M CPU @ 2.50GHz with 4Gb of RAM). It is GC friendly, queries require 0 allocations.
Building the index requires exactly 8 allocations.

To achieve this speed, the index has three restrictions. It is static, once built it cannot be changed.
It only accepts points, no bboxes or lines. It only accepts (for now) closest point queries.

Beware, to achieve top performance one of the hot functions has been rewritten in assembly.
Library works in x86 but it probably won't work in other architectures. PRs are welcome to fix this deficiency.

![Simple Recursive Layout](./example.png?raw=true "Simple Recursive Layout")

### Usage

The format of the points is a single array where each too coordinates represent a point

    points := []float64{0.0, 0.0, 1.0, 1.0} // array of two points 0, 0 and 1, 1

The library exposes only two methods. One to load and one to find nearest point

    import "SimpleRTree"

    fp := SimpleRTree.FlatPoints(points)
    r := SimpleRTree.New().Load(fp)
    closestX, closestY, distanceSquared, found := r.FindNearestPoint(1.0, 3.0)
    // 1.0, 1.0, 4.0, true

Additionally the tree can be built using options:

    r := SimpleRTree.NewWithOptions(Options{...})

For example:

    NewWithOptions(Options{UnsafeConcurrencyMode:true})

Is a slightly more efficient rtree, but which cannot be accessed from more than one goroutine at the same time.

### Benchmark. CPU

These are the benchmarks for finding the nearest point once the index has been built.

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	 5000000	       252 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 2000000	       881 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	 1000000	      1298 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	 1000000	      1828 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	 1000000	      1998 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000000-4 	  500000	      2364 ns/op

### Benchmark Load CPU

These are the benchmarks for the initial load

    BenchmarkSimpleRTree_Load/10-4      	  500000	      2206 ns/op
    BenchmarkSimpleRTree_Load/1000-4    	   10000	    262154 ns/op
    BenchmarkSimpleRTree_Load/10000-4   	     500	   3847292 ns/op
    BenchmarkSimpleRTree_Load/100000-4  	      30	  36191977 ns/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  72849948 ns/op

### Benchmark Load mem

    BenchmarkSimpleRTree_Load/10-4      	 1000000	      2205 ns/op	    2048 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/1000-4    	   10000	    262363 ns/op	   83350 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/10000-4   	     500	   3784520 ns/op	  804430 B/op	       7 allocs/op
    BenchmarkSimpleRTree_Load/100000-4  	      30	  41201134 ns/op	 8005321 B/op	       8 allocs/op
    BenchmarkSimpleRTree_Load/200000-4  	      20	  75697528 ns/op	16008905 B/op	       8 allocs/op

## Benchmark Compute distances

    Benchmark_ComputeDistances-4         	50000000	        23.9 ns/op
    Benchmark_VectorComputeDistances-4   	100000000	        12.8 ns/op




