package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the service
type Config struct {
	// Pipeline configuration
	Pipeline struct {
		NumWorkers        int `json:"num_workers"`
		AggregationWindow int `json:"aggregation_window"`
		TasksPerSecond    int `json:"tasks_per_second"`
		BurstSize         int `json:"burst_size"`
	} `json:"pipeline"`

	// Service configuration
	Service struct {
		LogLevel   string `json:"log_level"`
		Debug      bool   `json:"debug"`
		TimeFormat string `json:"time_format"`
		PrettyLog  bool   `json:"pretty_log"`
	} `json:"service"`

	// Buffer sizes
	BufferSizes struct {
		InputChannel  int `json:"input_channel"`
		ResultChannel int `json:"result_channel"`
	} `json:"buffer_sizes"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	cfg := &Config{}

	// Pipeline defaults
	cfg.Pipeline.NumWorkers = 10
	cfg.Pipeline.AggregationWindow = 50
	cfg.Pipeline.TasksPerSecond = 100
	cfg.Pipeline.BurstSize = 200

	// Service defaults
	cfg.Service.LogLevel = "info"
	cfg.Service.Debug = false
	cfg.Service.TimeFormat = "2006-01-02T15:04:05.000Z07:00"
	cfg.Service.PrettyLog = true

	// Buffer defaults
	cfg.BufferSizes.InputChannel = 1000
	cfg.BufferSizes.ResultChannel = 1000

	return cfg
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	// Pipeline config
	if v := os.Getenv("PIPELINE_NUM_WORKERS"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.Pipeline.NumWorkers = i
		}
	}
	if v := os.Getenv("PIPELINE_AGGREGATION_WINDOW"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.Pipeline.AggregationWindow = i
		}
	}
	if v := os.Getenv("PIPELINE_TASKS_PER_SECOND"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.Pipeline.TasksPerSecond = i
		}
	}
	if v := os.Getenv("PIPELINE_BURST_SIZE"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.Pipeline.BurstSize = i
		}
	}

	// Service config
	if v := os.Getenv("SERVICE_LOG_LEVEL"); v != "" {
		c.Service.LogLevel = v
	}
	if v := os.Getenv("SERVICE_DEBUG"); v != "" {
		c.Service.Debug = v == "true"
	}
	if v := os.Getenv("SERVICE_TIME_FORMAT"); v != "" {
		c.Service.TimeFormat = v
	}
	if v := os.Getenv("SERVICE_PRETTY_LOG"); v != "" {
		c.Service.PrettyLog = v == "true"
	}

	// Buffer sizes
	if v := os.Getenv("BUFFER_INPUT_CHANNEL"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.BufferSizes.InputChannel = i
		}
	}
	if v := os.Getenv("BUFFER_RESULT_CHANNEL"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.BufferSizes.ResultChannel = i
		}
	}
}

// LoadFromFile loads configuration from a JSON file
func (c *Config) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Pipeline.NumWorkers <= 0 {
		return fmt.Errorf("number of workers must be greater than 0")
	}
	if c.Pipeline.AggregationWindow <= 0 {
		return fmt.Errorf("aggregation window must be greater than 0")
	}
	if c.Pipeline.TasksPerSecond <= 0 {
		return fmt.Errorf("tasks per second must be greater than 0")
	}
	if c.Pipeline.BurstSize < c.Pipeline.TasksPerSecond {
		c.Pipeline.BurstSize = c.Pipeline.TasksPerSecond
	}
	if c.BufferSizes.InputChannel <= 0 {
		return fmt.Errorf("input channel buffer size must be greater than 0")
	}
	if c.BufferSizes.ResultChannel <= 0 {
		return fmt.Errorf("result channel buffer size must be greater than 0")
	}
	return nil
}
