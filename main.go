package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

type (
	product struct {
		XMLName xml.Name `xml:"Product"`
		Color   string   `xml:"Color"`
	}

	document struct {
		XMLName  xml.Name  `xml:"Products"`
		Products []product `xml:"Product"`
	}

	Args struct {
		Version string
		Profile string
		Name    string
	}
)

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
		f, _ := os.Create(fmt.Sprintf("profiles/cpu/prof_v%s_%s.prof", args.Version, args.Name))
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		return solutionFn(color, docs), nil
	case "mem":
		f, _ := os.Create(fmt.Sprintf("profiles/mem/prof_v%s_%s.prof", args.Version, args.Name))
		defer f.Close()

		runtime.GC() // ensure a clean memory state
		result := solutionFn(color, docs)

		pprof.WriteHeapProfile(f)

		return result, nil
	case "block":
		f, _ := os.Create(fmt.Sprintf("profiles/block/prof_v%s_%s.prof", args.Version, args.Name))
		defer f.Close()

		blockProfile := pprof.Lookup("block")
		blockProfile.WriteTo(f, 0)

		return solutionFn(color, docs), nil
	case "go":
		f, _ := os.Create(fmt.Sprintf("profiles/go/prof_v%s_%s.prof", args.Version, args.Name))
		defer f.Close()

		pprof.Lookup("goroutine").WriteTo(f, 0)

		return solutionFn(color, docs), nil
	case "mut":
		f, _ := os.Create(fmt.Sprintf("profiles/mut/prof_v%s_%s.prof", args.Version, args.Name))
		defer f.Close()

		mutexProfile := pprof.Lookup("mutex")
		mutexProfile.WriteTo(f, 0)
		return solutionFn(color, docs), nil
	default:
		return 0, fmt.Errorf("invalid profile: %s", args.Profile)
	}
}

func freq(color string, docs []string) int {
	args := parseArgs()
	var solutionFn func(color string, docs []string) int

	switch args.Version {
	case "1":
		solutionFn = freqSequentialV1
	case "2":
		solutionFn = freqSequentialV2
	default:
		log.Fatalf("invalid version: %s", args.Version)
	}

	if args.Profile != "" && args.Name != "" {
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
		Version: os.Getenv("version"),
		Profile: os.Getenv("profile"),
		Name:    os.Getenv("name"),
	}
}
