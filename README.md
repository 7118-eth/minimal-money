# 💰 Minimal Money

<p align="center">
  <strong>A beautiful terminal-based portfolio tracker that respects your time and privacy</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat" alt="License">
  <img src="https://img.shields.io/badge/Platform-macOS%20|%20Linux-blue?style=flat" alt="Platform">
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
```

## 🏗 Architecture

Built with modern Go practices:
- **Bubble Tea** - Delightful TUI framework
- **GORM** - Type-safe database operations
- **SQLite** - Zero-config persistence
- **Repository Pattern** - Clean data access
- **Service Layer** - Business logic separation

See [TEST_ARCHITECTURE.md](TEST_ARCHITECTURE.md) for our pragmatic testing approach.

## 🤝 Contributing

Contributions are welcome! This project values:
- Clean, readable code
- Pragmatic solutions
- Real-world testing
- User experience

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

<p align="center">
  Made with ❤️ by developers who track their portfolios daily
</p>