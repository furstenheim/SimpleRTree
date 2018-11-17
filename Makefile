debug:
	dlv test -- -test.run Big$
bench:
	go test -v -bench=Find
## Show allocs per test
bench-mem:
	go test -v -bench=Find -benchmem
## Show allocs trace
bench-mem-trace:
	go test -c
	GODEBUG=allocfreetrace=1 ./SimpleRTree.test -test.run=none -test.benchtime=10ms -test.bench=BenchmarkSimpleRTree_FindNearestPointMemory 2>trace.log
bench-graph:
	mkdir -p benchmarks/$$(date +%F)$$(git rev-parse HEAD)
	go test -run=XXX -bench Find -cpuprofile benchmarks/$$(date +%F)$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg SimpleRTree.test benchmarks/$$(date +%F)$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(date +%F)$$(git rev-parse HEAD)/cpu.svg
bench-load:
	go test -v -bench=Load
## Show allocs per test
bench-mem-load:
	go test -v -bench=Load -benchmem
## Show allocs trace
bench-mem-trace-load:
	go test -c
	GODEBUG=allocfreetrace=1 ./SimpleRTree.test -test.run=none -test.benchtime=10ms -test.bench=BenchmarkSimpleRTree_LoadMemory 2>trace.log
bench-graph-load:
	mkdir -p benchmarks/$$(date +%F)$$(git rev-parse HEAD)
	go test -run=XXX -bench Load -cpuprofile benchmarks/$$(date +%F)$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg SimpleRTree.test benchmarks/$$(date +%F)$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(date +%F)$$(git rev-parse HEAD)/cpu.svg
	echo "File at benchmarks/$$(date +%F)$$(git rev-parse HEAD)/cpu.svg"
bench-graph-load-mem:
	mkdir -p benchmarks/$$(date +%F)$$(git rev-parse HEAD)
	go test -run=XXX -bench Load -memprofile benchmarks/$$(date +%F)$$(git rev-parse HEAD)/heap.prof
	go tool pprof -lines -sample_index=alloc_objects -svg SimpleRTree.test benchmarks/$$(date +%F)$$(git rev-parse HEAD)/heap.prof > benchmarks/$$(date +%F)$$(git rev-parse HEAD)/heap.svg
	echo "File at benchmarks/$$(date +%F)$$(git rev-parse HEAD)/heap.svg"

bench-compute-distances:
	go test -run=Compute -bench Compute

compile-to-assembly:
	go build -gcflags -S . 2> a