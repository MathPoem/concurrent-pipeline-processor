package pipeline

import (
	"context"
	"errors"
	"sync"

	"concurrent-pipeline-processor/internal/processor"
	"concurrent-pipeline-processor/pkg/models"

	"golang.org/x/time/rate"
)

var (
	ErrPipelineNotStarted = errors.New("pipeline not started")
	ErrPipelineStopped    = errors.New("pipeline stopped")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
)

type pipeline struct {
	opts Options

	input     chan models.Task
	validated chan models.Task
	processed chan models.Result
	output    chan models.Result

	started bool
	stopped bool
	mu      sync.RWMutex
	wg      sync.WaitGroup
	limiter *rate.Limiter
}

// NewPipeline creates a new pipeline with the given options
func NewPipeline(opts Options) (Pipeline, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	// Create rate limiter with burst size
	limiter := rate.NewLimiter(rate.Limit(opts.TasksPerSecond), opts.BurstSize)

	return &pipeline{
		opts:      opts,
		input:     make(chan models.Task, opts.InputBufferSize),
		validated: make(chan models.Task, opts.InputBufferSize),
		processed: make(chan models.Result, opts.ResultBufferSize),
		output:    make(chan models.Result, opts.ResultBufferSize),
		limiter:   limiter,
	}, nil
}

func (p *pipeline) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.started {
		p.mu.Unlock()
		return errors.New("pipeline already started")
	}
	p.started = true
	p.mu.Unlock()

	// Start the pipeline stages
	p.wg.Add(3) // validator, processor, aggregator

	// Start validator
	go p.runValidator(ctx)

	// Start processor workers
	go p.runProcessor(ctx)

	// Start aggregator
	go p.runAggregator(ctx)

	// Monitor context cancellation
	go func() {
		<-ctx.Done()
		p.shutdown()
	}()

	return nil
}

func (p *pipeline) AddTask(task models.Task) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.started {
		return ErrPipelineNotStarted
	}
	if p.stopped {
		return ErrPipelineStopped
	}

	// Try to acquire rate limit token
	if !p.limiter.Allow() {
		return ErrRateLimitExceeded
	}

	select {
	case p.input <- task:
		return nil
	default:
		return errors.New("pipeline buffer full")
	}
}

func (p *pipeline) Results() <-chan models.Result {
	return p.output
}

func (p *pipeline) shutdown() {
	p.mu.Lock()
	if !p.stopped {
		p.stopped = true
		close(p.input)
	}
	p.mu.Unlock()

	p.wg.Wait()
	close(p.output)
}

// Stage implementations will be added in separate files
func (p *pipeline) runValidator(ctx context.Context) {
	defer func() {
		close(p.validated)
		p.wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-p.input:
			if !ok {
				return
			}
			if err := processor.ValidateTask(task); err != nil {
				p.output <- models.Result{Error: err}
				continue
			}
			select {
			case p.validated <- task:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (p *pipeline) runProcessor(ctx context.Context) {
	defer p.wg.Done()

	// Create worker pool
	var wg sync.WaitGroup
	wg.Add(p.opts.NumWorkers)

	// Start workers
	for i := 0; i < p.opts.NumWorkers; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-p.validated:
					if !ok {
						return
					}
					result := processor.ProcessTask(task)
					select {
					case p.processed <- result:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	// Wait for all workers to finish and close processed channel
	go func() {
		wg.Wait()
		close(p.processed)
	}()
}

func (p *pipeline) runAggregator(ctx context.Context) {
	defer p.wg.Done()

	agg := processor.NewAggregator(p.opts.AggregationWindow)
	defer agg.Close()

	resultChan := agg.Results()

	for {
		select {
		case <-ctx.Done():
			agg.Flush()
			return
		case result, ok := <-p.processed:
			if !ok {
				agg.Flush()
				return
			}
			agg.Add(result)
		case result := <-resultChan:
			select {
			case p.output <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}
