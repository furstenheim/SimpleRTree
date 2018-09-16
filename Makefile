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
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench Find -cpuprofile benchmarks/$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg SimpleRTree.test benchmarks/$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(git rev-parse HEAD)/cpu.svg
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
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench Load -cpuprofile benchmarks/$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg SimpleRTree.test benchmarks/$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(git rev-parse HEAD)/cpu.svg
	echo "File at benchmarks/$$(git rev-parse HEAD)/cpu.svg"