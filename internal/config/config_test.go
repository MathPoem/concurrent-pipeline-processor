package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Test pipeline defaults
	if cfg.Pipeline.NumWorkers != 10 {
		t.Errorf("Expected NumWorkers=10, got %d", cfg.Pipeline.NumWorkers)
	}
	if cfg.Pipeline.AggregationWindow != 50 {
		t.Errorf("Expected AggregationWindow=50, got %d", cfg.Pipeline.AggregationWindow)
	}
	if cfg.Pipeline.TasksPerSecond != 100 {
		t.Errorf("Expected TasksPerSecond=100, got %d", cfg.Pipeline.TasksPerSecond)
	}
	if cfg.Pipeline.BurstSize != 200 {
		t.Errorf("Expected BurstSize=200, got %d", cfg.Pipeline.BurstSize)
	}

	// Test service defaults
	if cfg.Service.LogLevel != "info" {
		t.Errorf("Expected LogLevel=info, got %s", cfg.Service.LogLevel)
	}
	if cfg.Service.Debug {
		t.Error("Expected Debug=false")
	}
	if cfg.Service.TimeFormat != "2006-01-02T15:04:05.000Z07:00" {
		t.Errorf("Expected TimeFormat=2006-01-02T15:04:05.000Z07:00, got %s", cfg.Service.TimeFormat)
	}
	if !cfg.Service.PrettyLog {
		t.Error("Expected PrettyLog=true")
	}

	// Test buffer defaults
	if cfg.BufferSizes.InputChannel != 1000 {
		t.Errorf("Expected InputChannel=1000, got %d", cfg.BufferSizes.InputChannel)
	}
	if cfg.BufferSizes.ResultChannel != 1000 {
		t.Errorf("Expected ResultChannel=1000, got %d", cfg.BufferSizes.ResultChannel)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"PIPELINE_NUM_WORKERS":        "20",
		"PIPELINE_AGGREGATION_WINDOW": "100",
		"PIPELINE_TASKS_PER_SECOND":   "200",
		"PIPELINE_BURST_SIZE":         "400",
		"SERVICE_LOG_LEVEL":           "debug",
		"SERVICE_DEBUG":               "true",
		"SERVICE_TIME_FORMAT":         "2006-01-02",
		"SERVICE_PRETTY_LOG":          "false",
		"BUFFER_INPUT_CHANNEL":        "2000",
		"BUFFER_RESULT_CHANNEL":       "2000",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	cfg := DefaultConfig()
	cfg.LoadFromEnv()

	// Test pipeline values
	if cfg.Pipeline.NumWorkers != 20 {
		t.Errorf("Expected NumWorkers=20, got %d", cfg.Pipeline.NumWorkers)
	}
	if cfg.Pipeline.AggregationWindow != 100 {
		t.Errorf("Expected AggregationWindow=100, got %d", cfg.Pipeline.AggregationWindow)
	}
	if cfg.Pipeline.TasksPerSecond != 200 {
		t.Errorf("Expected TasksPerSecond=200, got %d", cfg.Pipeline.TasksPerSecond)
	}
	if cfg.Pipeline.BurstSize != 400 {
		t.Errorf("Expected BurstSize=400, got %d", cfg.Pipeline.BurstSize)
	}

	// Test service values
	if cfg.Service.LogLevel != "debug" {
		t.Errorf("Expected LogLevel=debug, got %s", cfg.Service.LogLevel)
	}
	if !cfg.Service.Debug {
		t.Error("Expected Debug=true")
	}
	if cfg.Service.TimeFormat != "2006-01-02" {
		t.Errorf("Expected TimeFormat=2006-01-02, got %s", cfg.Service.TimeFormat)
	}
	if cfg.Service.PrettyLog {
		t.Error("Expected PrettyLog=false")
	}

	// Test buffer values
	if cfg.BufferSizes.InputChannel != 2000 {
		t.Errorf("Expected InputChannel=2000, got %d", cfg.BufferSizes.InputChannel)
	}
	if cfg.BufferSizes.ResultChannel != 2000 {
		t.Errorf("Expected ResultChannel=2000, got %d", cfg.BufferSizes.ResultChannel)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid num workers",
			cfg: &Config{
				Pipeline: struct {
					NumWorkers        int "json:\"num_workers\""
					AggregationWindow int "json:\"aggregation_window\""
					TasksPerSecond    int "json:\"tasks_per_second\""
					BurstSize         int "json:\"burst_size\""
				}{
					NumWorkers:        0,
					AggregationWindow: 50,
					TasksPerSecond:    100,
					BurstSize:         200,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid aggregation window",
			cfg: &Config{
				Pipeline: struct {
					NumWorkers        int "json:\"num_workers\""
					AggregationWindow int "json:\"aggregation_window\""
					TasksPerSecond    int "json:\"tasks_per_second\""
					BurstSize         int "json:\"burst_size\""
				}{
					NumWorkers:        10,
					AggregationWindow: 0,
					TasksPerSecond:    100,
					BurstSize:         200,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tasks per second",
			cfg: &Config{
				Pipeline: struct {
					NumWorkers        int "json:\"num_workers\""
					AggregationWindow int "json:\"aggregation_window\""
					TasksPerSecond    int "json:\"tasks_per_second\""
					BurstSize         int "json:\"burst_size\""
				}{
					NumWorkers:        10,
					AggregationWindow: 50,
					TasksPerSecond:    0,
					BurstSize:         200,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary config file
	content := []byte(`{
		"pipeline": {
			"num_workers": 20,
			"aggregation_window": 100,
			"tasks_per_second": 200,
			"burst_size": 400
		},
		"service": {
			"log_level": "debug",
			"debug": true,
			"time_format": "2006-01-02",
			"pretty_log": false
		},
		"buffer_sizes": {
			"input_channel": 2000,
			"result_channel": 2000
		}
	}`)

	tmpfile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading from file
	cfg := DefaultConfig()
	err = cfg.LoadFromFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	// Verify loaded values
	if cfg.Pipeline.NumWorkers != 20 {
		t.Errorf("Expected NumWorkers=20, got %d", cfg.Pipeline.NumWorkers)
	}
	if cfg.Service.LogLevel != "debug" {
		t.Errorf("Expected LogLevel=debug, got %s", cfg.Service.LogLevel)
	}
	if cfg.BufferSizes.InputChannel != 2000 {
		t.Errorf("Expected InputChannel=2000, got %d", cfg.BufferSizes.InputChannel)
	}

	// Test loading from non-existent file
	err = cfg.LoadFromFile("non-existent.json")
	if err == nil {
		t.Error("Expected error when loading from non-existent file")
	}
}
