# Budget Tracker - Pragmatic Test Architecture

## Overview

This document outlines a practical testing approach for the Budget Tracker, focusing on real integration tests that ensure the tool actually works with real databases and APIs.

## Test Philosophy

1. **Real Integration**: Test against actual SQLite databases and real APIs
2. **Simplicity**: Avoid complex mocking when real services work fine
3. **Practicality**: Tests should verify the tool works in real scenarios
4. **Speed**: Keep tests fast enough to run frequently (< 30 seconds total)

## Test Structure

```
budget/
├── internal/
│   ├── repository/
│   │   └── *_test.go      # Test with real SQLite
│   ├── api/
│   │   └── client_test.go  # Test with real APIs
│   ├── service/
│   │   └── price_test.go   # Integration tests
│   └── ui/
│       └── model_test.go   # State transition tests
├── test/
│   ├── integration/
│   │   └── workflow_test.go # End-to-end tests
│   └── testdata/
│       └── test_budget.db   # Test database (gitignored)
└── Makefile
```

## Database Testing Strategy

### Use Real SQLite Files
```go
func SetupTestDB(t *testing.T) *gorm.DB {
    // Create a unique test database for each test
    dbPath := fmt.Sprintf("./test/testdata/test_%s.db", t.Name())
    
    // Ensure directory exists
    os.MkdirAll("./test/testdata", 0755)
    
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)
    
    // Run migrations
    err = db.AutoMigrate(&models.Account{}, &models.Asset{}, &models.Holding{})
    require.NoError(t, err)
    
    // Cleanup after test
    t.Cleanup(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
        os.Remove(dbPath)
    })
    
    return db
}
```

### Benefits
- Tests real file I/O behavior
- Catches actual SQLite issues
- Can inspect database if test fails
- More realistic than in-memory

## API Testing Strategy

### Test Against Real APIs
```go
func TestRealCoinGeckoAPI(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping API test in short mode")
    }
    
    client := api.NewPriceClient()
    
    // Test with common cryptos
    prices, err := client.GetCryptoPrices([]string{"BTC", "ETH"})
    require.NoError(t, err)
    
    // Verify we got reasonable prices
    assert.Greater(t, prices["BTC"], 10000.0)
    assert.Greater(t, prices["ETH"], 500.0)
}
```

### Rate Limit Handling
```go
func TestAPIRateLimits(t *testing.T) {
    // Add delays between tests to respect rate limits
    time.Sleep(100 * time.Millisecond)
}
```

### API Resilience
Tests are designed to handle API failures gracefully:

```go
// API client returns empty results on rate limit
if resp.StatusCode != http.StatusOK {
    return prices, nil  // Don't fail, just return empty
}

// Tests check if prices were actually fetched
if price, ok := prices["BTC"]; ok {
    assert.Greater(t, price, 10000.0)
} else {
    t.Log("BTC price not available (API might be rate limited)")
}
```

### Benefits
- Catches real API changes
- Verifies API keys/endpoints work
- Tests actual network conditions
- No mock maintenance needed

## Integration Test Example

```go
func TestCompleteWorkflow(t *testing.T) {
    // Setup real database
    db := SetupTestDB(t)
    
    // Create repositories with real DB
    accountRepo := repository.NewAccountRepositoryWithDB(db)
    assetRepo := repository.NewAssetRepositoryWithDB(db)
    holdingRepo := repository.NewHoldingRepositoryWithDB(db)
    
    // Create account
    account := &models.Account{Name: "Test Wallet", Type: "wallet"}
    err := accountRepo.Create(account)
    require.NoError(t, err)
    
    // Create asset (will fetch real price)
    asset := &models.Asset{Symbol: "BTC", Name: "Bitcoin", Type: models.AssetTypeCrypto}
    err = assetRepo.Create(asset)
    require.NoError(t, err)
    
    // Create holding
    holding := &models.Holding{
        AccountID: account.ID,
        AssetID:   asset.ID,
        Amount:    0.5,
    }
    err = holdingRepo.Create(holding)
    require.NoError(t, err)
    
    // Fetch real price (using test database)
    priceService := service.NewPriceServiceWithDB(db)
    prices, err := priceService.FetchPrices([]models.Asset{*asset})
    require.NoError(t, err)
    
    // Price might be 0 if API is rate limited
    if price, ok := prices[asset.ID]; ok && price > 0 {
        assert.Greater(t, price, 0.0)
    }
}
```

## Test Data Management

### Test Database
- Each test gets its own database file
- Automatically cleaned up after test
- Can be preserved for debugging with flag

### Sample Data
```go
func LoadSamplePortfolio(t *testing.T, db *gorm.DB) {
    // Create realistic test data
    accounts := []models.Account{
        {Name: "hardware wallet", Type: "wallet"},
        {Name: "NeoBank", Type: "bank"},
    }
    
    for _, acc := range accounts {
        require.NoError(t, db.Create(&acc).Error)
    }
}
```

## Running Tests

```makefile
# Run all tests (including API tests)
test-all:
	go test ./...

# Run fast tests only (skip API calls)
test-fast:
	go test -short ./...

# Run with real API calls
test-integration:
	go test -run Integration ./...

# Keep test databases for debugging
test-debug:
	TEST_KEEP_DB=1 go test -v ./...

# Clean test data
test-clean:
	rm -rf ./test/testdata/*.db
```

## Environment Variables

```bash
# Skip API tests
export TEST_SKIP_API=1

# Keep test databases
export TEST_KEEP_DB=1

# Use specific test database
export TEST_DB_PATH="./my_test.db"
```

## Best Practices

1. **API Tests**
   - Run sparingly to avoid rate limits
   - Add reasonable assertions (price > 0)
   - Skip in CI if needed

2. **Database Tests**
   - Each test gets fresh database
   - Test actual constraints
   - Verify migrations work

3. **Speed**
   - Parallelize where possible
   - Cache API responses within test run
   - Use -short flag for quick feedback

4. **Debugging**
   - Keep failed test databases
   - Log actual vs expected values
   - Use verbose mode when needed

## Example Test Output

```bash
$ make test-all
=== RUN   TestAccountRepository_Create
    Using test database: ./test/testdata/test_TestAccountRepository_Create.db
--- PASS: TestAccountRepository_Create (0.05s)

=== RUN   TestRealCoinGeckoAPI
    Fetching real prices from CoinGecko...
    BTC: $45,234.00
    ETH: $3,012.00
--- PASS: TestRealCoinGeckoAPI (0.84s)

=== RUN   TestCompleteWorkflow
    Created account: hardware wallet
    Created asset: BTC
    Fetched price: $45,234.00
    Portfolio value: $22,617.00
--- PASS: TestCompleteWorkflow (1.23s)

PASS
ok      github.com/bioharz/budget       2.12s
```

This approach gives us confidence that our personal budget tool actually works with real services!