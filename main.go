package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"

	"trialday/solution"
)

type Args struct {
	Version  string
	Profile  string
	Solution string
	Trace    bool
}

const (
	docsNumber  = 5000
	lookUpColor = "White"
)

func main() {
	docs := make([]string, docsNumber)
	for i := range docs {
		docs[i] = fmt.Sprintf("data-%.4d.xml", i)
	}

	n := freq(lookUpColor, docs)

	log.Printf("Searching through %d files with products. Found products with color: %s %d times.", len(docs), lookUpColor, n)
}

func profilingWrapper(color string, docs []string, args Args, solutionFn func(color string, docs []string) int) (int, error) {
	switch args.Profile {
	case "cpu":
		f, _ := os.Create(fmt.Sprintf("profiles/cpu/prof_v%s_%s.prof", args.Version, args.Solution))
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		return solutionFn(color, docs), nil
	case "mem":
		f, _ := os.Create(fmt.Sprintf("profiles/mem/prof_v%s_%s.prof", args.Version, args.Solution))
		defer f.Close()

		runtime.GC() // ensure a clean memory state
		result := solutionFn(color, docs)

		pprof.WriteHeapProfile(f)

		return result, nil
	case "block":
		f, _ := os.Create(fmt.Sprintf("profiles/block/prof_v%s_%s.prof", args.Version, args.Solution))
		defer f.Close()

		blockProfile := pprof.Lookup("block")
		blockProfile.WriteTo(f, 0)

		return solutionFn(color, docs), nil
	case "go":
		f, _ := os.Create(fmt.Sprintf("profiles/go/prof_v%s_%s.prof", args.Version, args.Solution))
		defer f.Close()

		pprof.Lookup("goroutine").WriteTo(f, 0)

		return solutionFn(color, docs), nil
	case "mut":
		f, _ := os.Create(fmt.Sprintf("profiles/mut/prof_v%s_%s.prof", args.Version, args.Solution))
		defer f.Close()

		mutexProfile := pprof.Lookup("mutex")
		mutexProfile.WriteTo(f, 0)
		return solutionFn(color, docs), nil
	default:
		return 0, fmt.Errorf("invalid profile: %s", args.Profile)
	}
}

// TODO add trace, compare worker pool solutions

func freq(color string, docs []string) int {
	args := parseArgs()
	var solutionFn func(color string, docs []string) int

	switch args.Solution {
	case "sequential":
		switch args.Version {
		case "1":
			solutionFn = solution.FreqSequentialV1
		case "2":
			solutionFn = solution.FreqSequentialV2
		default:
			log.Fatalf("invalid version: %s", args.Version)
		}
	case "worker_pool":
		switch args.Version {
		case "1":
			solutionFn = solution.FreqWorkerPoolV1
		case "2":
			solutionFn = solution.FreqWorkerPoolV2
		case "3":
			solutionFn = solution.FreqWorkerPoolV3
		case "4":
			solutionFn = solution.FreqWorkerPoolV4
		case "5":
			solutionFn = solution.FreqWorkerPoolV5
		case "6":
			solutionFn = solution.FreqWorkerPoolV6
		default:
			log.Fatalf("invalid version: %s", args.Version)
		}
	}

	if args.Trace {
		traceFile, err := os.OpenFile(
			fmt.Sprintf("traces/trace_%s_%s.out", args.Solution, args.Version),
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			0644,
		)
		if err != nil {
			log.Fatalf("failed to open trace file: %v", err)
		}
		defer traceFile.Close()
		if err := trace.Start(traceFile); err != nil {
			log.Fatalf("failed to start trace: %v", err)
		}
		defer trace.Stop()
	}

	if args.Profile != "" {
		result, err := profilingWrapper(color, docs, args, solutionFn)
		if err != nil {
			log.Fatal(err)
		}

		return result
	}

	return solutionFn(color, docs)
}

func parseArgs() Args {
	return Args{
		Version:  os.Getenv("version"),
		Profile:  os.Getenv("profile"),
		Solution: os.Getenv("solution"),
		Trace:    os.Getenv("trace") == "1",
	}
}
