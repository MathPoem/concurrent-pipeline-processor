package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "default config",
			cfg: Config{
				Level:      "info",
				Debug:      false,
				TimeFormat: time.RFC3339,
				Pretty:     false,
			},
		},
		{
			name: "debug config",
			cfg: Config{
				Level:      "debug",
				Debug:      true,
				TimeFormat: time.RFC3339,
				Pretty:     false,
			},
		},
		{
			name: "custom time format",
			cfg: Config{
				Level:      "info",
				Debug:      false,
				TimeFormat: "2006-01-02",
				Pretty:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Initialize(tt.cfg)
			logger := GetLogger()

			if tt.cfg.Debug && logger.GetLevel() != zerolog.DebugLevel {
				t.Error("Expected debug level when debug is true")
			}

			if !tt.cfg.Debug && logger.GetLevel() != getLogLevel(tt.cfg.Level) {
				t.Error("Expected configured level when debug is false")
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	log = zerolog.New(&buf)

	tests := []struct {
		name     string
		logFunc  func() *zerolog.Event
		level    string
		message  string
		wantText string
	}{
		{
			name:     "debug level",
			logFunc:  Debug,
			level:    "debug",
			message:  "debug message",
			wantText: "debug message",
		},
		{
			name:     "info level",
			logFunc:  Info,
			level:    "info",
			message:  "info message",
			wantText: "info message",
		},
		{
			name:     "warn level",
			logFunc:  Warn,
			level:    "warn",
			message:  "warn message",
			wantText: "warn message",
		},
		{
			name:     "error level",
			logFunc:  Error,
			level:    "error",
			message:  "error message",
			wantText: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc().Msg(tt.message)

			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("Failed to parse log entry: %v", err)
			}

			if level, ok := logEntry["level"].(string); !ok || level != tt.level {
				t.Errorf("Expected level %s, got %s", tt.level, level)
			}

			if msg, ok := logEntry["message"].(string); !ok || msg != tt.wantText {
				t.Errorf("Expected message %s, got %s", tt.wantText, msg)
			}
		})
	}
}

func TestWithField(t *testing.T) {
	var buf bytes.Buffer
	log = zerolog.New(&buf)

	log.Info().Str("key", "value").Msg("test message")

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if value, ok := logEntry["key"].(string); !ok || value != "value" {
		t.Errorf("Expected field value 'value', got %v", value)
	}
}

func TestWithError(t *testing.T) {
	var buf bytes.Buffer
	log = zerolog.New(&buf)

	testErr := &testError{"test error"}
	WithError(testErr).Msg("error occurred")

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if errMsg, ok := logEntry["error"].(string); !ok || !strings.Contains(errMsg, "test error") {
		t.Errorf("Expected error message containing 'test error', got %v", errMsg)
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"fatal", zerolog.FatalLevel},
		{"invalid", zerolog.InfoLevel},
		{"", zerolog.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			if got := getLogLevel(tt.level); got != tt.expected {
				t.Errorf("getLogLevel(%q) = %v, want %v", tt.level, got, tt.expected)
			}
		})
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
