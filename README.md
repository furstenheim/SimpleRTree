## Simple RTree

Simple RTree is a blazingly fast and GC friendly RTree. It performs in 1.6 microseconds with 1 Million points for closest point queries
(measured in a i5-2450M CPU @ 2.50GHz with 4Gb of RAM). It is GC friendly, queries require 0 allocations.
Building the index requires exactly 8 allocations.

To achieve this speed, the index has three restrictions. It is static, once built it cannot be changed.
It only accepts points, no bboxes or lines. And it only accepts (for now) one query, closest point to a given coordinate.

Beware, to achieve top performance one of the hot functions has been rewritten in assembly.
Library works in x86 but it probably won't work in other architectures. PRs are welcome to fix this deficiency.

![Simple Recursive Layout](./example.png?raw=true "Simple Recursive Layout")

### Basic Usage

The format of the points is a single array where each too coordinates represent a point


    import "SimpleRTree"
    points := []float64{0.0, 0.0, 1.0, 1.0} // array of two points 0, 0 and 1, 1

    fp := SimpleRTree.FlatPoints(points)
    r := SimpleRTree.New().Load(fp)
    closestX, closestY, distanceSquared := r.FindNearestPoint(1.0, 3.0)
    // 1.0, 1.0, 4.0


## Documentation
To access the whole documentation you can access the following [link](https://godoc.org/github.com/furstenheim/SimpleRTree#SimpleRTree).

### Benchmark. CPU

These are the benchmarks for finding the nearest point once the index has been built.

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	10000000	       173 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 3000000	       635 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	 2000000	       863 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	 1000000	      1127 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	 1000000	      1357 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000000-4 	 1000000	      1593 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000000-4         	 1000000	      1992 ns/op









