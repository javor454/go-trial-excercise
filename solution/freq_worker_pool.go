package solution

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

// FreqWorkerPoolV1 waitgroup is used to wait for all documents to be processed
func FreqWorkerPoolV1(color string, docs []string) int {
	numWorkers := min(runtime.NumCPU(), len(docs)) // avoid unnecessary goroutines

	jobs := make(chan string, len(docs))
	results := make(chan int, len(docs))
	var wg sync.WaitGroup

	wg.Add(len(docs))

	// start workers
	for range numWorkers {
		go func() {
			for job := range jobs {
				results <- FreqSequentialV2(color, []string{job})
				wg.Done()
			}
		}()
	}

	// send jobs
	go func() {
		for _, doc := range docs {
			jobs <- doc
		}
		close(jobs) // signal no more jobs
	}()

	// all results should be collected
	go func() {
		wg.Wait()
		close(results)
	}()

	// range should exit whenever results channel is closed and it ranges over all results
	// TODO only one goroutine counts the total
	total := 0
	for result := range results {
		total += result
	}

	return total
}

// FreqWorkerPoolV2 waitgroup is used to wait for all workers to finish and close results channel
func FreqWorkerPoolV2(color string, docs []string) int {
	numWorkers := min(runtime.NumCPU(), len(docs))
	// TODO runtime.GOMAXPROCS() // avoid unnecessary goroutines
	// TODO traces - how much threads im using

	jobs := make(chan string, len(docs)) // TODO buffer equal to number of workers
	results := make(chan int, len(docs)) // TODO buffer equal to number of workers - cant send more results than
	// TODO use more goroutines for results
	// TODO count total in each worker
	var wg sync.WaitGroup

	wg.Add(numWorkers)

	// wait for all workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// start workers
	for range numWorkers {
		go func() {
			// range until jobs channel is closed
			for job := range jobs {
				results <- FreqSequentialV2(color, []string{job})
			}
			wg.Done()
		}()
	}

	// send jobs
	go func() {
		for _, doc := range docs {
			jobs <- doc
		}
		close(jobs) // signal no more jobs
	}()

	// range should exit whenever results channel is closed and it ranges over all results
	total := 0
	for result := range results {
		total += result
	}

	return total
}

func FreqWorkerPoolV3(color string, docs []string) int {
	numWorkers := min(runtime.NumCPU(), len(docs)) // avoid unnecessary goroutines

	jobs := make(chan string, len(docs))
	results := make(chan int, len(docs))
	var wg sync.WaitGroup

	wg.Add(len(docs))

	var productPool = &sync.Pool{
		New: func() any {
			return new(Product)
		},
	}

	// start workers
	for range numWorkers {
		go func() {
			for job := range jobs {
				results <- freqStream(color, job, productPool)
				wg.Done()
			}
		}()
	}

	// send jobs
	go func() {
		for _, doc := range docs {
			jobs <- doc
		}
		close(jobs) // signal no more jobs
	}()

	// all results should be collected
	go func() {
		wg.Wait()
		close(results)
	}()

	// range should exit whenever results channel is closed and it ranges over all results
	total := 0
	for result := range results {
		total += result
	}

	return total
}

func freqStream(color string, doc string, productPool *sync.Pool) int {
	var found int

	fileName := fmt.Sprintf("%s.xml", doc[:4])

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatal(err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "Product" {
				// get product from pool
				product := productPool.Get().(*Product)
				// reset fields
				*product = Product{}

				err = decoder.DecodeElement(&product, &t)
				if err != nil {
					log.Fatal(err)
				}

				if strings.Contains(product.Color, color) {
					found++
				}
			}
		}
	}

	return found
}
