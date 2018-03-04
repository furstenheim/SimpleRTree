### Simple RTree

This is a static RTree that only accepts points as input


### Benchmark. CPU

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  500000	      4104 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  100000	     17384 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	  100000	     29522 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     31716 ns/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     33006 ns/op

### Benchmark. Memory

    BenchmarkSimpleRTree_FindNearestPoint/10-4      	  300000	      4798 ns/op	      59 B/op	       2 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	  100000	     18148 ns/op	     189 B/op	      10 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	   50000	     29156 ns/op	     269 B/op	      15 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	   50000	     31629 ns/op	     359 B/op	      21 allocs/op
    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	   50000	     33116 ns/op	     392 B/op	      23 allocs/op

