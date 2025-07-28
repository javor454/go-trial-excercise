.PHONY: bench run-profile run build benchstat

bench: ## make bench version=1 solution=sequential will produce file benchmarks/bench_v1_sequential.txt
	go test -benchmem -benchtime=5s -count=6 -cpu=1,2,4,8,12 -bench=. | tee benchmarks/bench_v$(version)_$(solution).txt

run-profile: build
	PROFILE=$(profile) VERSION=$(version) SOLUTION=$(solution) ./target/main

run: build ## versions: 1, 2 solutions: sequential, worker_pool
	VERSION=$(version) SOLUTION=$(solution) ./target/main

build:
	go build -o target/main .

benchstat:
	./script/benchstat.sh $(old) $(new)