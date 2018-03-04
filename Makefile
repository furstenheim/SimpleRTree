debug:
	dlv test -- -test.run Big$
bench:
	go test -v -bench=.
## Show allocs per test
bench-mem:
	go test -v -bench=. -benchmem
## Show allocs trace
bench-mem-trace:
	go test -c
	GODEBUG=allocfreetrace=1 ./SimpleRTree.test -test.run=none -test.benchtime=10ms -test.bench=BenchmarkSimpleRTree_FindNearestPointMemory 2>trace.log
bench-graph:
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench . -cpuprofile benchmarks/$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg SimpleRTree.test benchmarks/$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(git rev-parse HEAD)/cpu.svg