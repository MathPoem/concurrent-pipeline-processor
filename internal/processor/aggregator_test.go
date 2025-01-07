package processor

import (
	"testing"
	"time"

	"concurrent-pipeline-processor/pkg/models"
)

func TestAggregator(t *testing.T) {
	t.Run("aggregates results in window", func(t *testing.T) {
		agg := NewAggregator(3)
		defer agg.Close()

		// Add results
		go func() {
			agg.Add(models.Result{Result: 1})
			agg.Add(models.Result{Result: 2})
			agg.Add(models.Result{Result: 3})
		}()

		// Get aggregated result
		select {
		case result := <-agg.Results():
			if result.Result != 6 { // 1 + 2 + 3
				t.Errorf("Expected sum 6, got %d", result.Result)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for aggregated result")
		}
	})

	t.Run("handles error results", func(t *testing.T) {
		agg := NewAggregator(3)
		defer agg.Close()

		errResult := models.Result{Error: ErrInvalidOperator}

		// Add error result
		go func() {
			agg.Add(errResult)
		}()

		// Error should be forwarded immediately
		select {
		case result := <-agg.Results():
			if result.Error != ErrInvalidOperator {
				t.Errorf("Expected error %v, got %v", ErrInvalidOperator, result.Error)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for error result")
		}
	})

	t.Run("flushes partial window", func(t *testing.T) {
		agg := NewAggregator(3)
		defer agg.Close()

		// Add partial window
		go func() {
			agg.Add(models.Result{Result: 1})
			agg.Add(models.Result{Result: 2})
			agg.Flush()
		}()

		// Get partial sum
		select {
		case result := <-agg.Results():
			if result.Result != 3 { // 1 + 2
				t.Errorf("Expected sum 3, got %d", result.Result)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for flushed result")
		}
	})
}
