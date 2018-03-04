debug:
	dlv test -- -test.run Big$
bench:
	go test -v -bench=.
bench-mem:
	go test -v -bench=. -benchmem
bench-graph:
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench . -cpuprofile benchmarks/$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg SimpleRTree.test benchmarks/$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(git rev-parse HEAD)/cpu.svg