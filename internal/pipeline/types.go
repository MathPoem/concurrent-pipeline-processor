package pipeline

import (
	"context"
	"errors"

	"concurrent-pipeline-processor/pkg/models"
)

var (
	// ErrInvalidNumWorkers is returned when the number of workers is invalid
	ErrInvalidNumWorkers = errors.New("number of workers must be greater than 0")
	// ErrInvalidAggregationWindow is returned when the aggregation window size is invalid
	ErrInvalidAggregationWindow = errors.New("aggregation window must be greater than 0")
	// ErrInvalidRateLimit is returned when the rate limit is invalid
	ErrInvalidRateLimit = errors.New("rate limit must be greater than 0")
)

// Pipeline represents the main interface for the concurrent pipeline processor
type Pipeline interface {
	// Start initializes and starts the pipeline
	Start(ctx context.Context) error
	// AddTask adds a new task to the pipeline
	AddTask(task models.Task) error
	// Results returns a channel for receiving processed results
	Results() <-chan models.Result
}

// Options contains configuration options for the pipeline
type Options struct {
	// NumWorkers specifies the number of concurrent workers for processing tasks
	NumWorkers int
	// AggregationWindow specifies the size of the window for aggregating results
	AggregationWindow int
	// TasksPerSecond specifies the maximum number of tasks that can be processed per second
	TasksPerSecond int
	// BurstSize specifies the maximum number of tasks that can be processed in a burst
	BurstSize int
	// InputBufferSize specifies the size of the input channel buffer
	InputBufferSize int
	// ResultBufferSize specifies the size of the result channel buffer
	ResultBufferSize int
}

// Validate checks if the options are valid
func (o Options) Validate() error {
	if o.NumWorkers <= 0 {
		return ErrInvalidNumWorkers
	}
	if o.AggregationWindow <= 0 {
		return ErrInvalidAggregationWindow
	}
	if o.TasksPerSecond <= 0 {
		return ErrInvalidRateLimit
	}
	if o.BurstSize < o.TasksPerSecond {
		o.BurstSize = o.TasksPerSecond
	}
	if o.InputBufferSize <= 0 {
		o.InputBufferSize = o.NumWorkers * o.AggregationWindow * 2
	}
	if o.ResultBufferSize <= 0 {
		o.ResultBufferSize = o.NumWorkers * o.AggregationWindow * 2
	}
	return nil
}
