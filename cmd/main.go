package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"concurrent-pipeline-processor/internal/config"
	"concurrent-pipeline-processor/internal/logger"
	"concurrent-pipeline-processor/internal/pipeline"
	"concurrent-pipeline-processor/pkg/models"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "", "path to config file")
	flag.Parse()

	// Load configuration
	cfg := config.DefaultConfig()

	// Load from environment variables
	cfg.LoadFromEnv()

	// Load from config file if specified
	if *configFile != "" {
		if err := cfg.LoadFromFile(*configFile); err != nil {
			fmt.Printf("Failed to load config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Initialize(logger.Config{
		Level:      cfg.Service.LogLevel,
		Debug:      cfg.Service.Debug,
		TimeFormat: cfg.Service.TimeFormat,
		Pretty:     cfg.Service.PrettyLog,
	})

	log := logger.GetLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create pipeline with configuration
	p, err := pipeline.NewPipeline(pipeline.Options{
		NumWorkers:        cfg.Pipeline.NumWorkers,
		AggregationWindow: cfg.Pipeline.AggregationWindow,
		TasksPerSecond:    cfg.Pipeline.TasksPerSecond,
		BurstSize:         cfg.Pipeline.BurstSize,
		InputBufferSize:   cfg.BufferSizes.InputChannel,
		ResultBufferSize:  cfg.BufferSizes.ResultChannel,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pipeline")
	}

	// Start pipeline
	if err := p.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start pipeline")
	}

	// Log configuration
	log.Info().
		Int("workers", cfg.Pipeline.NumWorkers).
		Int("aggregation_window", cfg.Pipeline.AggregationWindow).
		Int("tasks_per_second", cfg.Pipeline.TasksPerSecond).
		Int("burst_size", cfg.Pipeline.BurstSize).
		Int("input_buffer", cfg.BufferSizes.InputChannel).
		Int("result_buffer", cfg.BufferSizes.ResultChannel).
		Bool("debug", cfg.Service.Debug).
		Msg("Starting pipeline with configuration")

	// Handle shutdown signals
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		sig := <-signals
		log.Info().Str("signal", sig.String()).Msg("Shutdown signal received, stopping pipeline...")
		cancel()
	}()

	// Add tasks
	go func() {
		defer cancel() // Cancel context when done adding tasks

		for i := 0; i < 1000; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				task := generateTask()
				for {
					err := p.AddTask(task)
					if err == nil {
						log.Debug().
							Int("task_value", task.Value).
							Int("operations", len(task.Operations)).
							Msg("Task added successfully")
						break
					}
					if err == pipeline.ErrRateLimitExceeded {
						log.Debug().Msg("Rate limit exceeded, waiting to retry")
						select {
						case <-ctx.Done():
							return
						case <-time.After(time.Millisecond * 10):
							continue
						}
					} else if err.Error() == "pipeline buffer full" {
						log.Debug().Msg("Pipeline buffer full, waiting to retry")
						select {
						case <-ctx.Done():
							return
						case <-time.After(time.Millisecond * 10):
							continue
						}
					} else {
						log.Error().Err(err).Msg("Failed to add task")
						return
					}
				}
			}
		}
	}()

	// Collect results with timestamps
	startTime := time.Now()
	taskCount := 0
	for result := range p.Results() {
		taskCount++
		elapsed := time.Since(startTime)
		rate := float64(taskCount) / elapsed.Seconds()

		if result.Error != nil {
			log.Error().
				Err(result.Error).
				Int("task_count", taskCount).
				Float64("current_rate", rate).
				Msg("Error processing task")
			continue
		}

		log.Info().
			Int("result", result.Result).
			Int("task_count", taskCount).
			Float64("current_rate", rate).
			Dur("elapsed", elapsed).
			Msg("Task processed successfully")
	}
}

func generateTask() models.Task {
	value := rand.Intn(100)
	numberOfOperations := rand.Intn(5)
	operations := make([]models.Operation, numberOfOperations)

	for j := 0; j < numberOfOperations; j++ {
		operations[j] = models.Operation{
			Operator: models.Operator(rand.Intn(int(models.OperatorTotalAmount))),
			Value:    rand.Intn(10) + 1, // Avoid zero for division
		}
	}

	return models.Task{
		Value:      value,
		Operations: operations,
	}
}
