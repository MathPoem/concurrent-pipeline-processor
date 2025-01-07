package processor

import (
	"concurrent-pipeline-processor/pkg/models"
	"time"
)

// Aggregator handles the aggregation of results
type Aggregator struct {
	window    int
	buffer    []models.Result
	aggregate chan models.Result
}

// NewAggregator creates a new Aggregator with the specified window size
func NewAggregator(window int) *Aggregator {
	return &Aggregator{
		window:    window,
		buffer:    make([]models.Result, 0, window),
		aggregate: make(chan models.Result, window*2),
	}
}

// Add adds a result to the aggregator
func (a *Aggregator) Add(result models.Result) {
	if result.Error != nil {
		select {
		case a.aggregate <- result:
		case <-time.After(time.Second):
			// If we can't send after timeout, this is a serious issue
			panic("aggregator blocked for too long")
		}
		return
	}

	a.buffer = append(a.buffer, result)
	if len(a.buffer) >= a.window {
		a.flush()
	}
}

// Flush forces aggregation of any remaining results
func (a *Aggregator) Flush() {
	if len(a.buffer) > 0 {
		a.flush()
	}
}

// Results returns the channel for aggregated results
func (a *Aggregator) Results() <-chan models.Result {
	return a.aggregate
}

// Close closes the aggregator
func (a *Aggregator) Close() {
	if len(a.buffer) > 0 {
		a.flush()
	}
	close(a.aggregate)
}

func (a *Aggregator) flush() {
	if len(a.buffer) == 0 {
		return
	}

	sum := 0
	for _, r := range a.buffer {
		sum += r.Result
	}

	select {
	case a.aggregate <- models.Result{Result: sum}:
		a.buffer = a.buffer[:0]
	case <-time.After(time.Second):
		// If we can't send after timeout, this is a serious issue
		panic("aggregator blocked for too long")
	}
}
