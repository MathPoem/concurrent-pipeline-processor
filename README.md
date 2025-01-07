# Concurrent Pipeline Processor

A high-performance, concurrent task processing pipeline implemented in Go. The service processes mathematical operations in a pipeline pattern with fan-out/fan-in concurrency, rate limiting, and window-based result aggregation.

## Features

- Concurrent task processing with configurable worker pool
- Fan-out/fan-in concurrency pattern
- Rate limiting with burst support
- Window-based result aggregation
- Graceful shutdown handling
- Structured logging with multiple output formats
- Configurable via environment variables and JSON files
- Comprehensive error handling and validation
- Full test coverage
- Static code analysis with multiple linters
- Docker containerization with multi-stage builds
- Docker Compose support for easy deployment
- Git-based version control with comprehensive ignore patterns

## Architecture

The service implements a pipeline with four main stages:

1. **Input Stage**: Receives tasks and applies rate limiting
2. **Validator Stage**: Validates tasks and their operations
3. **Processor Stage**: Processes tasks concurrently with multiple workers
4. **Aggregator Stage**: Aggregates results in configurable windows

### Task Processing

Tasks consist of a base value and a series of mathematical operations:
- Addition
- Subtraction
- Multiplication
- Division (with zero division protection)

## Configuration

Configuration can be provided through environment variables. Default values are set in the Dockerfile and can be overridden through docker-compose.yml or environment variables.

### Environment Variables

```bash
# Pipeline Configuration
PIPELINE_NUM_WORKERS=10
PIPELINE_AGGREGATION_WINDOW=50
PIPELINE_TASKS_PER_SECOND=100
PIPELINE_BURST_SIZE=200

# Service Configuration
SERVICE_LOG_LEVEL=info
SERVICE_DEBUG=false
SERVICE_TIME_FORMAT="2006-01-02T15:04:05.000Z07:00"
SERVICE_PRETTY_LOG=true

# Buffer Configuration
BUFFER_INPUT_CHANNEL=1000
BUFFER_RESULT_CHANNEL=1000
```

### Configuration File (config.json)

```json
{
    "pipeline": {
        "num_workers": 10,
        "aggregation_window": 50,
        "tasks_per_second": 100,
        "burst_size": 200
    },
    "service": {
        "log_level": "info",
        "debug": false,
        "time_format": "2006-01-02T15:04:05.000Z07:00",
        "pretty_log": true
    },
    "buffer_sizes": {
        "input_channel": 1000,
        "result_channel": 1000
    }
}
```

## Usage

### Building and Running Locally

```bash
# Copy environment file
cp .env.example .env

# Build the service
go build ./cmd/main.go

# Run with environment variables
./main
```

### Running with Docker

```bash
# Build the Docker image
docker build -t pipeline-processor .

# Run the container with environment variables
docker run -p 8080:8080 \
  -e PIPELINE_NUM_WORKERS=10 \
  -e PIPELINE_AGGREGATION_WINDOW=50 \
  -e PIPELINE_TASKS_PER_SECOND=100 \
  pipeline-processor
```

### Running with Docker Compose

```bash
# Start the service (uses environment variables from docker-compose.yml)
docker compose up -d

# Start with custom environment file
docker compose --env-file .env.custom up -d

# View logs
docker compose logs -f

# Stop the service
docker compose down
```

### Example Task

```go
task := models.Task{
    Value: 10,
    Operations: []models.Operation{
        {Operator: models.OperatorPlus, Value: 5},    // 15
        {Operator: models.OperatorMultiply, Value: 2}, // 30
        {Operator: models.OperatorMinus, Value: 5},    // 25
    },
}
```

## Error Handling

The service handles various error conditions:
- Division by zero
- Invalid operators
- Rate limit exceeded
- Buffer full conditions
- Invalid configuration
- Graceful shutdown

## Logging

Structured logging is implemented using zerolog with support for:
- Multiple log levels (debug, info, warn, error, fatal)
- JSON and pretty-printed formats
- Custom time formats
- Contextual fields
- Error tracking

### Testing

The service includes comprehensive tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests verbosely
go test ./... -v
```

## Performance Considerations

- Uses buffered channels for improved throughput
- Configurable worker pool size
- Rate limiting to prevent overload
- Window-based aggregation for efficient processing
- Optimized buffer sizes for different stages
