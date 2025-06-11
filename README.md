# 💰 Minimal Money

<p align="center">
  <strong>A beautiful terminal-based portfolio tracker that respects your time and privacy</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat" alt="License">
  <img src="https://img.shields.io/badge/Platform-macOS%20|%20Linux-blue?style=flat" alt="Platform">
  <img src="https://github.com/7118-eth/minimal-money/workflows/CI/badge.svg" alt="CI Status">
  <img src="https://github.com/7118-eth/minimal-money/workflows/golangci-lint/badge.svg" alt="Lint Status">
</p>

## ✨ Features

### 🌳 **Intuitive Tree View**
- Asset-first organization with accounts as branches
- Visual hierarchy inspired by `htop`
- Full terminal width utilization

### 💸 **Real-Time Pricing**
- Live crypto prices via CoinGecko
- Fiat exchange rates via ExchangeRate-API
- Smart caching to minimize API calls
- Manual refresh with `p` key

### 📊 **Multi-Account Support**
- Organize holdings by exchange, wallet, or bank
- Track assets across multiple platforms
- See total value per asset across all accounts

### 🔍 **Complete Audit Trail**
- Track every portfolio change
- Know exactly when and what was added/edited/deleted
- Essential for tax reporting

### ⚡ **Lightning Fast**
- SQLite for instant data access
- Keyboard-driven interface
- No unnecessary features or bloat

## 📸 Screenshots

```
💰 Minimal Money                               Total: $51,284.43
                                    Last Update: 2025-01-11 14:22:18

Asset/Account                    Amount                Value
BTC                              0.7250                $29,450.00
  ├─ Hardware Wallet             0.4500                $18,270.00
  ├─ CoinBase                    0.1800                $7,308.00
  └─ Gemini                      0.0950                $3,872.00
ETH                              4.2000                $10,080.00
  ├─ Hardware Wallet             2.8000                $6,720.00
  └─ Binance                     1.4000                $3,360.00
USD                              8,750.00              $8,750.00
  ├─ CityTrust                   5,200.00              $5,200.00
  ├─ FirstBank                   2,100.00              $2,100.00
  └─ GlobalBank                  1,450.00              $1,450.00
GBP                              2,100.00              $2,604.00
  ├─ MonzoBank                   1,400.00              $1,736.00
  └─ BarclaysBank                700.00                $868.00

[n]ew  [e]dit  [d]elete  [p]rice update  [h]istory  [q]uit
```

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/bioharz/minimal-money.git
cd minimal-money

# Build the application
make build

# Run it!
./minimal-money
```

## ⌨️ Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `n` | Add new holding |
| `e` | Edit selected |
| `d` | Delete selected |
| `p` | Update prices |
| `h` | View audit history |
| `q` | Quit |
| `↑↓` | Navigate |
| `Tab` | Next field in forms |
| `Esc` | Cancel/Go back |

## 🛠 Development

### Setup

```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install git hooks (optional but recommended)
make install-hooks
```

### Commands

```bash
# Run directly
make run

# Run tests
make test          # All tests (including API calls)
make test-fast     # Skip API tests
make test-coverage # With coverage report

# Clean test databases
make test-clean

# Format code
make fmt

# Run linter
make lint

# Run both formatting and linting
make check
```

### Git Hooks

This project includes pre-commit hooks that automatically:
- Check code formatting with `gofmt`
- Run linting with `golangci-lint`
- Execute fast tests

To skip hooks temporarily: `git commit --no-verify`

### Code Quality Tools

#### golangci-lint

We use [golangci-lint](https://golangci-lint.run/) for comprehensive code analysis:

```bash
# Install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all linters
golangci-lint run

# Run with specific linter
golangci-lint run --enable-only gofmt

# Auto-fix issues (where possible)
golangci-lint run --fix
```

Configuration: `.golangci.yml` (simplified for CI compatibility)

## 🏗 Architecture

Built with modern Go practices:
- **Bubble Tea** - Delightful TUI framework
- **GORM** - Type-safe database operations
- **SQLite** - Zero-config persistence
- **Repository Pattern** - Clean data access
- **Service Layer** - Business logic separation

See [TEST_ARCHITECTURE.md](TEST_ARCHITECTURE.md) for our pragmatic testing approach.

## 🔄 CI/CD

This project uses GitHub Actions for continuous integration and deployment:

### Workflows

- **CI** - Runs on every push and PR
  - Linting and formatting checks
  - Tests on multiple Go versions (1.21, 1.22, 1.23)
  - Cross-platform testing (Linux, macOS)
  - Security scanning with govulncheck
  - Code coverage reporting
  - Binary builds for all platforms

- **Release** - Automated releases on version tags
  - Builds binaries for Linux and macOS (amd64, arm64)
  - Creates GitHub releases with checksums
  - Automated changelog in release notes

- **Code Quality** - golangci-lint for comprehensive checks

### Status Badges

![CI](https://github.com/7118-eth/minimal-money/workflows/CI/badge.svg)
![golangci-lint](https://github.com/7118-eth/minimal-money/workflows/golangci-lint/badge.svg)

### Testing Workflows Locally

You can test GitHub Actions workflows locally before pushing:

#### Using act

[act](https://github.com/nektos/act) runs your workflows locally in Docker containers:

```bash
# Install act
brew install act              # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash  # Linux

# List available workflows
act -l

# Run the default push event
act

# Run a specific workflow
act -W .github/workflows/ci.yml

# Run with specific event
act pull_request

# Run a specific job
act -j test-fast

# See what actions would run (dry run)
act -n

# For Apple Silicon (M1/M2/M3) Macs
act -P macos-latest=-self-hosted --container-architecture linux/amd64
```

**Note for Apple Silicon users**: The default Docker images may not work correctly on ARM64 architecture. Use the command above to specify x86_64 emulation.

#### Using GitHub CLI

```bash
# Manually trigger a workflow run
gh workflow run ci.yml

# View workflow runs
gh run list

# Watch a workflow run
gh run watch

# View workflow run details
gh run view <run-id>

# Download workflow artifacts
gh run download <run-id>
```

#### Docker Alternative

If you prefer Docker directly:

```bash
# Run act in Docker
docker run --rm -it \
  -v "$PWD:/workspace" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /workspace \
  nektos/act
```

## 🤝 Contributing

Contributions are welcome! This project values:
- Clean, readable code
- Pragmatic solutions
- Real-world testing
- User experience

All PRs must pass CI checks before merging.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

<p align="center">
  Made with ❤️ by developers who track their portfolios daily
</p>