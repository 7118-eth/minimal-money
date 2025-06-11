# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version info
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

RUN CGO_ENABLED=1 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
    -o minimal-money \
    cmd/budget/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Create non-root user
RUN addgroup -g 1000 minimal && \
    adduser -D -u 1000 -G minimal -h /home/minimal -s /bin/sh minimal

# Copy binary from builder
COPY --from=builder /app/minimal-money /usr/local/bin/minimal-money

# Set ownership
RUN chown minimal:minimal /usr/local/bin/minimal-money

# Switch to non-root user
USER minimal
WORKDIR /home/minimal

# Create data directory
RUN mkdir -p /home/minimal/data

# Volume for persistent data
VOLUME ["/home/minimal/data"]

# Set environment variable for data location
ENV MINIMAL_MONEY_DB_PATH=/home/minimal/data/budget.db

# Default command
ENTRYPOINT ["minimal-money"]

# Labels
LABEL org.opencontainers.image.title="Minimal Money"
LABEL org.opencontainers.image.description="Beautiful terminal-based portfolio tracker"
LABEL org.opencontainers.image.url="https://github.com/7118-eth/minimal-money"
LABEL org.opencontainers.image.source="https://github.com/7118-eth/minimal-money"
LABEL org.opencontainers.image.licenses="MIT"