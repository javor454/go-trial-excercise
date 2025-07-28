#!/bin/bash

# Check if required arguments are provided
if [ $# -lt 2 ]; then
    echo "Usage: $0 <old_benchmark_file> <new_benchmark_file>"
    echo "Example: $0 benchmarks/bench_v1_baseline.txt benchmarks/bench_v2_optimized.txt"
    exit 1
fi

OLD_BENCHMARK="$1"
NEW_BENCHMARK="$2"

# Check if benchmark files exist
if [ ! -f "$OLD_BENCHMARK" ]; then
    echo "Error: Old benchmark file '$OLD_BENCHMARK' not found"
    exit 1
fi

if [ ! -f "$NEW_BENCHMARK" ]; then
    echo "Error: New benchmark file '$NEW_BENCHMARK' not found"
    exit 1
fi

# Run benchstat using Docker
echo "Comparing benchmarks using Dockerized benchstat..."
echo "Old: $OLD_BENCHMARK"
echo "New: $NEW_BENCHMARK"
echo ""

docker run --rm \
    -v "$(pwd):/work" \
    -w /work \
    golang:latest \
    sh -c "go install golang.org/x/perf/cmd/benchstat@latest && benchstat \"$OLD_BENCHMARK\" \"$NEW_BENCHMARK\""
