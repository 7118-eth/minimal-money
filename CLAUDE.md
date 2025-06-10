# Budget Tracker - Quick Reference

## Commands
```bash
# Run
go run cmd/budget/main.go

# Build
go build -o budget cmd/budget/main.go

# Test
go test ./...

# Format
go fmt ./...
```

## Project Structure
```
budget/
├── cmd/budget/      # Main entry
├── internal/        # Core logic
├── PROJECT.md       # Technical design
└── PROGRESS.md      # Current state
```

## Key Files
- Technical design: `PROJECT.md`
- Current progress: `PROGRESS.md`
- Main app: `cmd/budget/main.go`