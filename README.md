# Go Trial Day Excersize
- change the implementation of freq function to find all occurences of loookUpColor in the Color tag in xml file
- goal is to find the most optimal algorithm
    - implement solution
    - measure its performance
    - optimize
        - apply different concurency patterns
        - use profiling / tracing 


## Solutions
- Sequential 
    - V1
        - memory bottlenecks
            - xml parsing - unmarshalling whole document
            - solutions
                - streaming parser processing one product at a time instead of loading all into memory
    - V2 with streaming parser
        - compared to V1 on average
            - Time performance: +4.65%  (slower)
            - Memory usage:     -12.66% (less memory)
            - Allocations:      -0.20%  (fewer allocations)
- Worker pool
    - Implement a worker pool with a fixed number of goroutines
    - Use channels to distribute work (file names) to workers
    - Learn about channel buffering and synchronization
    - Explore different worker pool sizes and their impact
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