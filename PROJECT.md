# Budget Tracker - Technical Design Document

## Architecture Overview

### System Design
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Bubble    │────▶│   Business   │────▶│   SQLite    │
│   Tea UI    │     │    Logic     │     │  Database   │
└─────────────┘     └──────────────┘     └─────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │  Price APIs  │
                    │ (CoinGecko,  │
                    │  Exchange)   │
                    └──────────────┘
```

### Tech Stack
- **Language**: Go 1.24.3
- **UI Framework**: Bubble Tea v1.3.5 (TUI)
- **Database**: SQLite via GORM v1.30.0
- **Price APIs**: CoinGecko (crypto), ExchangeRate-API (fiat)

## Data Models

### Account
```go
type Account struct {
    ID        uint      
    Name      string    // "Hardware Wallet", "CityTrust", "CoinBase"
    Type      string    // "bank", "wallet", "exchange"
    Color     string    // For UI display
    CreatedAt time.Time
}
```

### Asset
```go
type Asset struct {
    ID        uint      
    Symbol    string    // "BTC", "ETH", "USD"
    Name      string    // "Bitcoin", "Ethereum", "US Dollar"
    Type      AssetType // "crypto", "fiat", "stock"
    CreatedAt time.Time
}
```

### Holding
```go
type Holding struct {
    ID            uint      
    AccountID     uint      // Link to account
    AssetID       uint      
    Amount        float64   
    PurchasePrice float64   // For P&L calculation
    PurchaseDate  time.Time // For holding period
    CreatedAt     time.Time
}
```

### PriceCache
```go
type PriceCache struct {
    ID         uint      
    AssetID    uint      
    PriceUSD   float64   
    UpdatedAt  time.Time 
}
```

### AuditLog
```go
type AuditLog struct {
    ID         uint      
    Action     string    
    EntityType string    
    EntityID   uint      
    OldValue   string    
    NewValue   string    
    CreatedAt  time.Time 
}
```

### Relationships
```
Account ──1:N──▶ Holding ◀──N:1── Asset
                    │
                    └──────▶ PriceCache
```

## UI/UX Design

### Main Table View
The application displays a tree-structured view of all assets and their holdings across different accounts, with real-time valuations and a total portfolio value at the top.

### Add Asset Modal
```
┌─────────────────────────────────────┐
│       Add New Asset                 │
├─────────────────────────────────────┤
│ Account:  [Hardware Wallet ▼]      │
│ Asset:    [BTC             ▼]      │
│ Amount:   [0.1250______]           │
│ Price:    [$40,600____] (optional) │
│                                     │
│        [Save]  [Cancel]             │
└─────────────────────────────────────┘
```

### Keyboard Navigation
- **Arrow Keys**: Navigate table rows
- **n**: New asset (opens modal)
- **e**: Edit selected row
- **d**: Delete selected row
- **p**: Update prices (manual refresh)
- **h**: View audit trail
- **q**: Quit application
- **Tab**: Navigate modal fields
- **Enter**: Confirm actions
- **Esc**: Cancel/back

## API Integration

### CoinGecko (Crypto Prices)
- **Endpoint**: `https://api.coingecko.com/api/v3/simple/price`
- **Rate Limit**: 10-30 calls/minute (free tier)
- **Caching**: 5 minutes
- **Supported**: BTC, ETH, and 10,000+ cryptocurrencies

### ExchangeRate-API (Fiat Rates)
- **Endpoint**: `https://api.exchangerate-api.com/v4/latest/USD`
- **Rate Limit**: 1,500 requests/month (free tier)
- **Caching**: 1 hour
- **Supported**: 160+ currencies

### Price Update Strategy
1. On startup: Load cached prices from database
2. Manual update: User presses 'p'
3. Smart batching: Single API call for multiple assets
4. Database caching: Store prices with timestamps
5. Show last update time in UI

## Technical Decisions

### Why Table-Based UI?
- **Always visible overview** - No menu navigation needed
- **Familiar interface** - Excel-like experience
- **Efficient workflow** - Quick keyboard shortcuts
- **Real-time updates** - See price changes immediately

### Why Accounts?
- **Real-world mapping** - Matches how users organize assets
- **Risk management** - See exposure per platform
- **Tax reporting** - Track holdings by location
- **Visual grouping** - Easier to scan and understand

### Why GORM?
- **Migrations** - Automatic schema updates
- **Type safety** - Go structs map to tables
- **Relationships** - Easy joins and preloading
- **Soft deletes** - Built-in data recovery

### Error Handling
- **API failures**: Use cached prices, show stale indicator
- **Input validation**: Clear error messages in modal
- **Database errors**: Graceful degradation, user notification
- **Network issues**: Offline mode with cached data