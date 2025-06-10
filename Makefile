.PHONY: all build run test test-unit test-integration test-coverage test-race clean

# Default target
all: build

# Build the application
build:
	go build -o budget cmd/budget/main.go

# Run the application
run:
	go run cmd/budget/main.go

# Run all tests
test: test-unit test-integration

# Run unit tests only (short mode skips integration tests)
test-unit:
	go test -short ./...

# Run integration tests only
test-integration:
	go test -run Integration ./test/integration

# Run tests with coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detector
test-race:
	go test -race ./...

# Run benchmarks
test-bench:
	go test -bench=. -benchmem ./...

# Clean build artifacts
clean:
	rm -f budget
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy