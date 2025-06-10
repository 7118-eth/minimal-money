# Budget Tool - Claude Development Guide

## Project Overview
Personal budget tracking tool built with Go, featuring:
- Asset management (BTC, ETH, USD, EUR, etc.)
- Real-time price fetching via APIs
- SQLite database with versioning/history
- Terminal UI using Bubble Tea framework

## Tech Stack
- **Language**: Go 1.24.3
- **Database**: SQLite with GORM v1.30.0
- **Terminal UI**: Bubble Tea v1.3.5 + Bubbles v0.21.0
- **Styling**: Lipgloss v1.1.0
- **Price APIs**: CoinGecko (crypto), ExchangeRate-API (fiat)

## Project Structure
```
budget/
├── cmd/budget/         # Main application entry point
├── internal/
│   ├── models/        # Database models and structs
│   ├── api/           # External API clients
│   ├── ui/            # Bubble Tea UI components
│   └── db/            # Database connection and migrations
├── pkg/utils/         # Shared utilities
└── data/              # SQLite database files
```

## Database Schema
- **assets**: id, symbol, name, type (crypto/fiat/stock)
- **holdings**: id, asset_id, amount, created_at, updated_at
- **price_history**: id, asset_id, price_usd, timestamp
- **portfolio_snapshots**: id, total_value_usd, timestamp, details (JSON)

## Key Commands
```bash
# Run the application
go run cmd/budget/main.go

# Build binary
go build -o budget cmd/budget/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code (install golangci-lint first)
golangci-lint run
```

## Development Guidelines
1. Use GORM for all database operations
2. Implement proper error handling with wrapped errors
3. Use Bubble Tea's Model-Update-View pattern
4. Keep API calls concurrent using goroutines
5. Store sensitive API keys in environment variables

## API Integration
- **CoinGecko**: Free tier allows 10-30 calls/minute
- **ExchangeRate-API**: Free tier allows 1500 requests/month
- Consider caching prices for 1-5 minutes to avoid rate limits

## Testing
- Unit tests for models and business logic
- Integration tests for API clients with mocks
- Manual testing for TUI interactions