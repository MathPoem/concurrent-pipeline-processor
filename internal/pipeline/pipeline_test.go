package pipeline

import (
	"context"
	"testing"
	"time"

	"concurrent-pipeline-processor/pkg/models"
)

func TestPipeline(t *testing.T) {
	t.Run("processes tasks successfully", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p, err := NewPipeline(Options{
			NumWorkers:        2,
			AggregationWindow: 2,
			TasksPerSecond:    100,
			BurstSize:         200,
			InputBufferSize:   100,
			ResultBufferSize:  100,
		})
		if err != nil {
			t.Fatalf("Failed to create pipeline: %v", err)
		}

		if err := p.Start(ctx); err != nil {
			t.Fatalf("Failed to start pipeline: %v", err)
		}

		// Add tasks that will sum to 10 (5 + 5)
		tasks := []models.Task{
			{
				Value: 2,
				Operations: []models.Operation{
					{Operator: models.OperatorPlus, Value: 3},
				},
			},
			{
				Value: 2,
				Operations: []models.Operation{
					{Operator: models.OperatorPlus, Value: 3},
				},
			},
		}

		for _, task := range tasks {
			if err := p.AddTask(task); err != nil {
				t.Fatalf("Failed to add task: %v", err)
			}
		}

		// Get aggregated result
		select {
		case result := <-p.Results():
			if result.Error != nil {
				t.Errorf("Unexpected error: %v", result.Error)
			}
			if result.Result != 10 { // (2+3) + (2+3)
				t.Errorf("Expected sum 10, got %d", result.Result)
			}
		case <-time.After(2 * time.Second):
			t.Error("Timeout waiting for results")
		}
	})

	t.Run("handles validation errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p, err := NewPipeline(Options{
			NumWorkers:        2,
			AggregationWindow: 2,
			TasksPerSecond:    100,
			BurstSize:         200,
			InputBufferSize:   100,
			ResultBufferSize:  100,
		})
		if err != nil {
			t.Fatalf("Failed to create pipeline: %v", err)
		}

		if err := p.Start(ctx); err != nil {
			t.Fatalf("Failed to start pipeline: %v", err)
		}

		// Add invalid task (nil operations)
		task := models.Task{
			Value:      5,
			Operations: nil,
		}

		if err := p.AddTask(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		// Should receive error
		select {
		case result := <-p.Results():
			if result.Error == nil {
				t.Error("Expected validation error, got nil")
			}
		case <-time.After(2 * time.Second):
			t.Error("Timeout waiting for error result")
		}
	})

	t.Run("handles graceful shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		p, err := NewPipeline(Options{
			NumWorkers:        2,
			AggregationWindow: 2,
			TasksPerSecond:    100,
			BurstSize:         200,
			InputBufferSize:   100,
			ResultBufferSize:  100,
		})
		if err != nil {
			t.Fatalf("Failed to create pipeline: %v", err)
		}

		if err := p.Start(ctx); err != nil {
			t.Fatalf("Failed to start pipeline: %v", err)
		}

		// Add a task
		task := models.Task{
			Value: 2,
			Operations: []models.Operation{
				{Operator: models.OperatorPlus, Value: 3},
			},
		}

		if err := p.AddTask(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		// Cancel context and wait a bit for shutdown to complete
		cancel()
		time.Sleep(100 * time.Millisecond)

		// Try to add task after shutdown
		err = p.AddTask(task)
		if err != ErrPipelineStopped {
			t.Errorf("Expected ErrPipelineStopped, got %v", err)
		}
	})

	t.Run("handles rate limiting", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p, err := NewPipeline(Options{
			NumWorkers:        2,
			AggregationWindow: 2,
			TasksPerSecond:    1, // Very low rate limit for testing
			BurstSize:         1,
			InputBufferSize:   100,
			ResultBufferSize:  100,
		})
		if err != nil {
			t.Fatalf("Failed to create pipeline: %v", err)
		}

		if err := p.Start(ctx); err != nil {
			t.Fatalf("Failed to start pipeline: %v", err)
		}

		task := models.Task{
			Value: 2,
			Operations: []models.Operation{
				{Operator: models.OperatorPlus, Value: 3},
			},
		}

		// First task should succeed
		if err := p.AddTask(task); err != nil {
			t.Fatalf("Failed to add first task: %v", err)
		}

		// Second task should be rate limited
		if err := p.AddTask(task); err != ErrRateLimitExceeded {
			t.Errorf("Expected ErrRateLimitExceeded, got %v", err)
		}
	})
}
