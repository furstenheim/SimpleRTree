## Simple RTree

Simple RTree is a blazingly fast and GC friendly RTree. It performs in 1.6 microseconds with 1 Million points for closest point queries
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

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	10000000	       176 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 2000000	       627 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	 2000000	       910 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	 1000000	      1165 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	 1000000	      1325 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000000-4 	 1000000	      1598 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000000-4         	 1000000	      2297 ns/op







