# Golang Engineer Test Task: Concurrent Pipeline Processor

## Task Description
Create a concurrent pipeline processor that handles tasks through multiple stages with fan-out and fan-in patterns, implementing proper error handling and graceful shutdown.

### Requirements

1. Create a pipeline that processes tasks through the following stages:
   - Input generator - generates or receives tasks
   - Validator - validates task data, like if []Operation is nil
   - Processor - processes tasks in parallel with NumWorkers. It should take Value from each task and perform operations one by one (check error like divide by zero)
   - Aggregator - summ all result in a group of tasks with size of AggregationWindow. If AggregationWindow = 50 so you need to sum every 50 tasks' results and return as final result of pipeline

2. Implementation requirements:
   - Use channels for communication between stages
   - Implement fan-out pattern for the processing stage (multiple goroutines processing tasks)
   - Implement fan-in pattern to collect results from processing goroutines
   - Implement proper error handling that stops the entire pipeline if any goroutine encounters an error
   - Ensure graceful shutdown of all goroutines (no goroutine leaks)
   - If error occured during Validation stage or Processing stage you need to return result with error and stop whole pipeline
   - Use context for cancellation
   - Implement proper resource cleanup

3. The task should simulate processing of the following data structure:
```go
type Operator int

const (
	OperatorPlus Operator = iota
	OperatorMinus
	OperatorDivide
	OperatorMultiply
	OperatorTotalAmount
)

type Operation struct {
	Operator Operator
	Value    int
}

type Task struct {
	Value      int
	Operations []Operation
}

type Result struct {
	Result int
	Error  error
}
```

### Technical Requirements

1. Code organization:
   - Unit tests for critical components
   - Proper error handling patterns

2. Must use:
   - Go channels
   - Goroutines
   - Context package
   - sync package (if needed)
   - Standard Go project layout

### Example Interface

Your solution should implement something similar to this interface:

```go
type Pipeline interface {
	Start(ctx context.Context) error
	AddTask(task Task) error
	Results() <-chan Result
}

type Options struct {
	NumWorkers        int
	AggregationWindow int
}

func NewPipeline(opts Options) Pipeline
```

### Evaluation Criteria

1. Concurrency patterns implementation
   - Proper use of channels
   - Correct fan-out/fan-in implementation
   - Error handling across goroutines
   - Graceful shutdown implementation

2. Code quality
   - Clean and maintainable code
   - Error handling
   - Testing coverage

3. Performance considerations
   - Efficient resource usage
   - No goroutine leaks
   - Proper context cancellation


### Sample Usage

The implementation should be usable in a way similar to this:

```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := NewPipeline(Options{
		NumWorkers:        10,
		AggregationWindow: 50,
	})

	// Start pipeline
	if err := p.Start(ctx); err != nil {
		log.Fatal(err)
	}
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
		fmt.Printf("SIGINT or SIGTERM received, shutdown\n")
	}()

	// Add tasks
	go func() {
		for i := 0; i < 1000; i++ {
			value := rand.Int()
			numberOfOperations := rand.Intn(1000)
			operations := make([]Operation, numberOfOperations)
			for j := 0; j < numberOfOperations; j++ {
				operations[j].Operator = Operator(rand.Intn(int(OperatorTotalAmount)))
				operations[j].Value = rand.Intn(100) + 1
			}
			task := Task{
				Value:      value,
				Operations: operations,
			}
			if err := p.AddTask(task); err != nil {
				log.Printf("Failed to add task: %v", err)
			}
		}
	}()

	// Collect results
	for result := range p.Results() {
		if result.Error != nil {
			log.Printf("Result error: %v", result.Error)
		}
		log.Printf("Result: %d", result.Result)
	}
}
```