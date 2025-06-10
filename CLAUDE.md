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
- `e` - Edit selected holding
- `d` - Delete selected holding
- `p` - Update prices (manual refresh)
- `h` - View audit trail
- `q` - Quit
- `ESC` - Go back / Cancel
- `Tab` - Navigate modal fields
- Arrow keys - Navigate table
- `Enter` - Confirm in modals

## Current Features
- ✅ Multi-account portfolio tracking
- ✅ Asset-first tree view with htop-style visualization
- ✅ Tree structure across all columns (Asset/Account, Amount, Value)
- ✅ Database-cached prices with timestamps
- ✅ Manual price updates (press 'p')
- ✅ Sorting by highest value first
- ✅ SQLite persistence with GORM
- ✅ Full terminal width responsive UI
- ✅ Last price update timestamp display
- ✅ Audit trail for all portfolio changes
- ✅ Complete CRUD operations (Create, Read, Update, Delete)
- ✅ Input validation and error handling
- ✅ Comprehensive test suite with real APIs

## API Integrations
- CoinGecko: Crypto prices (cached in DB)
- ExchangeRate-API: Fiat rates (cached in DB)
- Manual updates only (press 'p')