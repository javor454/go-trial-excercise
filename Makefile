.PHONY: bench run-profile run build benchstat

bench: ## make bench version=1 name=my_test will produce file benchmarks/bench_v1_my_test.txt
	go test -benchmem -benchtime=5s -count=6 -cpu=1,2,4,8,12 -bench=. | tee benchmarks/bench_v$(version)_$(name).txt

run-profile: build
	PROFILE=$(profile) VERSION=$(version) NAME=$(name) ./target/main

run: build
	./target/main

build:
	go build -o target/main .

benchstat:
	./script/benchstat.sh $(old) $(new)