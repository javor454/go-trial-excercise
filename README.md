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
|| Sequential V1 | Sequential V1 vs V2 | Sequential V2 vs Worker pool V1 | Worker pool V1 vs V2 | Worker Pool V1 vs V3 |
|-|-|-|-|-|-|
| **Time performance** | ~1.989 s/op          | 2.082 s/op ❌ (+4.65%)          | 849.6 ms/op ✅ (-59.18%) | 857.2 ms/op ❌ (+0.89%) | 856.3 ms/op ❌ (+0.78%)         |
| **Memory usage**     | ~1.03 GB/op          | 917 MB/op ✅ (-12.66%)          | same                     | same                    | 942.1 MB/op ❌ (+0.12%)         |
| **Allocations**      | ~25.18 mil allocs/op | 25.13 mil allocs/op ✅ (-0.20%) | same                     | same                    | 25.22 mil allocs/op ❌ (+0.38%) |


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