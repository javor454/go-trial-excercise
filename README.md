# Go Trial Day Excersize
- copy of https://github.com/marek-drapal/go-trial-excersize
    - is not linked in any way to not reveal solutions

## Assignment
- change the implementation of freq function to find all occurences of loookUpColor in the Color tag in xml file
- goal is to find the most optimal algorithm
    - implement solution
    - measure its performance
    - optimize
        - apply different concurency patterns
        - use profiling / tracing 


## Solutions
### Performance comparison table (geomean)
|                      | Sequential V1        | Sequential V1 vs V2             | Seq V2 vs Worker pool V1 | Worker pool V1 vs V2     | Worker Pool V1 vs V3            | Worker pool V1 vs V4    | Worker pool V4 vs V5             | Worker pool V5 vs V6           |
|----------------------|----------------------|---------------------------------|--------------------------|--------------------------|---------------------------------|-------------------------|----------------------------------|--------------------------------|
| **Time performance** | ~1.989 s/op          | 2.082 s/op ❌ (+4.65%)          | 849.6 ms/op ✅ (-59.18%) | 857.2 ms/op ❌ (+0.89%)  | 856.3 ms/op ❌ (+0.78%)         | 845.2 ms/op ✅ (-0.52%) | 773.6 ms/op ✅ (-8.48%)          | 299.0 ms/op ✅ (-61.34%)       |
| **Memory usage**     | ~1.03 GB/op          | 917 MB/op ✅ (-12.66%)          | same                     | same                     | 942.1 MB/op ❌ (+0.12%)         | same                    | 758.5 MB/op ✅ (-8.48%)          | 167.8 MB/op ✅ (-77.87%)       |
| **Allocations**      | ~25.18 mil allocs/op | 25.13 mil allocs/op ✅ (-0.20%) | same                     | same                     | 25.22 mil allocs/op ❌ (+0.38%) | same                    | 22.22 mil allocs/op ✅ (-11.55%) | 308.6 k allocs/op ✅ (-98.61%) |

### Measurment details
- Sequential 
    - V1
        - memory bottlenecks
            - xml parsing - unmarshalling whole document
            - solutions
                - streaming parser processing one product at a time instead of loading all into memory
    - V2 with streaming parser
- Worker pool
    - Implement a worker pool with a fixed number of goroutines
    - Use channels to distribute work (file names) to workers
        - options
            - file based distribution (each worker processed entire file)
            - product based distribution (each worker processes product)
    - Learn about channel buffering and synchronization
        - Unbuffered channel - sender blocks until receiver is ready
        - Buffered channel - sender doesnt block until buffer is full
    - Explore different worker pool sizes and their impact
        - why base numWorkers on runtime.NumCPU()?
            - task is IO (file reading) with CPU work (xml parsing)
            - for IO tasks you can use more goroutines than CPU cores
        - versions
            - under utilization: numWorkers = (runtime.NumCPU() / 2)
            - optimal for CPU-bound tasks: numWorkers = runtime.NumCPU()
            - over subscription: numWorkers = (runtime.NumCPU() * 2)
    - Try to use sync package
        - sync.WaitGroup for graceful shutdown
        - sync.Pool for xml decoder reuse
    - V1
        - numWorkers equal to number of CPUs -> seems to be optimal
        - waitgroup for all documents
        - cpu profile
            - most time spend on scheduling, locking, sleeping -> goroutines waiting for work, main goroutine waiting for results
        - memory profile
            - most memory consumed by xml parsing / unmarshalling
                - reuse buffers or decoders with sync.Pool
    - V2
        - waitgroup for workers
    - V3
        - sync.Pool for Products (not for xml decoder - this requires reseting the decoder which is not supported by its api)
    - V4
        - limit buffer for jobs and results to number of workers, count found colors per each worker, then sumup results from workers in main goroutine
    - V5
        - when traversing xml tokens, search for "Color" tokens instead of "Product", then decode only color tag instead of whole product
    - V6
        - load whole file into memory and replace xml parser with regex
        - is not optimal for big xml files
    
- Fan-in / Fan-out
    - Start multiple goroutines to process files concurrently
    - Collect results using channels
    - Learn about goroutine lifecycle management
- Pipeline pattern
    - Break the work into stages: file reading → XML parsing → color counting
    - Use channels to pass data between stages
    - Learn about pipeline design and backpressure
- Circuit breaker pattern
    - A design pattern that prevents cascading failures by temporarily stopping operations that are likely to fail.
    - If your XML files are on a slow network drive, one slow file shouldn't make your entire concurrent system wait.
- Actor model
    - A concurrency model where "actors" are the universal primitives - they communicate only by sending messages to each other.
    - Each actor is isolated, making the system easier to reason about and debug.
- Memory pooling for XML parsing
    - Creating new XML decoders and buffers for each file is expensive.
    - Reduces garbage collection pressure and improves performance by reusing objects instead of creating new ones.


## TODOs
- profiling support ✅
- context and cancellation
    - add timeouts and cancelation support
    - graceful shutdown
- error handling
    - in goroutines - send through channel so it's not lost
- consider using GOMAXPROCS ✅
    - GOMAXPROCS returns max number of OS threads that can execute go code simultaneosly
    - by default is = NumCPU, or can be set lower
    - makes sense to use for container CPU limits, pure CPU-bound work or artificially limit OS thread usage
- use tracing
    - cpu profile shows only what is happening in the program
    - tracing also shows what is not happening (waiting)
    - info:
        - https://blog.gopheracademy.com/advent-2017/go-execution-tracer/ 
        - https://www.sobyte.net/post/2022-03/go-execution-tracer-by-example/
    - Goroutine tab
        - shows goroutine state: GCWaiting (waiting during GC), Running (actively executing), Runnable (ready to run but waiting for processor)
    