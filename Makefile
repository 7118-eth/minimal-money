.PHONY: all build run test test-all test-fast test-coverage test-race test-clean clean fmt lint check install-hooks

# Default target
all: build

# Build the application
build:
	@VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "none"); \
	DATE=$$(date -u '+%Y-%m-%d_%H:%M:%S'); \
	go build -o minimal-money -ldflags="-s -w -X main.version=$$VERSION -X main.commit=$$COMMIT -X main.date=$$DATE" cmd/budget/main.go

# Run the application
run:
	go run cmd/budget/main.go

# Run all tests (including API tests)
test: test-all

# Run all tests including real API calls
test-all:
	go test ./...

# Run fast tests only (skip API calls)
test-fast:
	go test -short ./...

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

# Keep test databases for debugging
test-debug:
	TEST_KEEP_DB=1 go test -v ./...

# Clean test databases
test-clean:
	rm -rf ./test/testdata/*.db
	rm -rf ./test/testdata/*.db-shm
	rm -rf ./test/testdata/*.db-wal

# Clean all artifacts
clean: test-clean
	rm -f minimal-money
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	elif [ -x "$$(go env GOPATH)/bin/golangci-lint" ]; then \
		$$(go env GOPATH)/bin/golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Run formatting and linting
check: fmt lint

# Install git hooks
install-hooks:
	@./scripts/install-hooks.sh

# Install dependencies
deps:
	go mod download
	go mod tidy

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# Test release locally with GoReleaser
release-test:
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean --skip=publish; \
	else \
		echo "goreleaser not found. Install with: go install github.com/goreleaser/goreleaser@latest"; \
		exit 1; \
	fi

# Create a new release tag
release-tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release-tag VERSION=v1.0.0"; \
		exit 1; \
	fi; \
	git tag -a $(VERSION) -m "Release $(VERSION)"; \
	echo "Created tag $(VERSION). Push with: git push origin $(VERSION)"