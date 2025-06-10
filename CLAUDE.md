# Minimal Money - Quick Reference

## Commands
```bash
# Run
go run cmd/budget/main.go

# Build
go build -o minimal-money cmd/budget/main.go

# Test (all with real APIs)
make test

# Test (fast, skip APIs)
make test-fast

# Test with coverage
make test-coverage

# Format
go fmt ./...

# Clean test databases
make test-clean
```

## Project Structure
```
minimal-money/
├── cmd/budget/         # Main entry point
├── internal/           # Core logic
│   ├── api/           # Price API clients
│   ├── db/            # Database connection
│   ├── models/        # Data models
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   └── ui/            # Terminal UI
├── test/              # Test infrastructure
│   ├── fixtures/      # Test data builders
│   ├── helpers/       # Test utilities
│   └── integration/   # E2E tests
├── data/              # SQLite database
├── Makefile           # Build commands
├── PROJECT.md         # Technical design
├── PROGRESS.md        # Current state
├── README.md          # User documentation
└── TEST_ARCHITECTURE.md # Testing strategy
```

## Key Files
- Technical design: `PROJECT.md`
- Current progress: `PROGRESS.md`
- Test strategy: `TEST_ARCHITECTURE.md`
- Main app: `cmd/budget/main.go`

## Keyboard Controls
- `n` - Add new asset
- `e` - Edit selected (not connected)
- `d` - Delete selected (not connected)
- `p` - Update prices (manual refresh)
- `h` - View history
- `q` - Quit
- `ESC` - Go back
- `Tab` - Navigate modal fields
- Arrow keys - Navigate table

## Current Features
- ✅ Multi-account portfolio tracking
- ✅ Asset-first tree view with account grouping
- ✅ Database-cached prices with timestamps
- ✅ Manual price updates (press 'p')
- ✅ P&L calculation
- ✅ SQLite persistence
- ✅ Full terminal width table UI
- ✅ Last price update timestamp display
- ✅ Audit trail for portfolio changes
- ✅ Input validation
- ✅ Comprehensive tests

## API Integrations
- CoinGecko: Crypto prices (cached in DB)
- ExchangeRate-API: Fiat rates (cached in DB)
- Manual updates only (press 'p')