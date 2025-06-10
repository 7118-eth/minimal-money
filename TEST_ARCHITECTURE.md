# Budget Tracker - Test Architecture

## Overview

This document outlines the comprehensive testing strategy for the Budget Tracker application, covering unit tests, integration tests, and test infrastructure.

## Test Philosophy

1. **Fast Feedback**: Unit tests should run in milliseconds
2. **Isolation**: Tests should not depend on external services or state
3. **Clarity**: Test names should describe what and why
4. **Maintainability**: DRY principles apply to tests too
5. **Coverage**: Aim for 80%+ coverage on business logic

## Test Structure

```
budget/
├── internal/
│   ├── models/
│   │   └── asset_test.go
│   ├── repository/
│   │   ├── account_test.go
│   │   ├── asset_test.go
│   │   └── holding_test.go
│   ├── service/
│   │   └── price_test.go
│   ├── api/
│   │   └── client_test.go
│   └── ui/
│       ├── model_test.go
│       ├── table_test.go
│       └── modal_test.go
├── test/
│   ├── integration/
│   │   ├── workflow_test.go
│   │   └── api_test.go
│   ├── fixtures/
│   │   └── test_data.go
│   └── helpers/
│       ├── db.go
│       ├── api.go
│       └── ui.go
└── Makefile
```

## Testing Layers

### 1. Unit Tests

#### Repository Layer
```go
// Test with in-memory SQLite
func TestAccountRepository_Create(t *testing.T) {
    db := test.SetupTestDB(t)
    repo := repository.NewAccountRepository(db)
    
    account := &models.Account{
        Name: "Test Account",
        Type: "wallet",
    }
    
    err := repo.Create(account)
    assert.NoError(t, err)
    assert.NotZero(t, account.ID)
}
```

#### Service Layer
```go
// Test with mocked dependencies
func TestPriceService_FetchPrices(t *testing.T) {
    mockClient := mocks.NewMockPriceClient(t)
    service := service.NewPriceService(mockClient)
    
    mockClient.On("GetCryptoPrices", []string{"BTC"}).
        Return(map[string]float64{"BTC": 45000}, nil)
    
    prices, err := service.FetchPrices([]models.Asset{
        {ID: 1, Symbol: "BTC", Type: models.AssetTypeCrypto},
    })
    
    assert.NoError(t, err)
    assert.Equal(t, 45000.0, prices[1])
}
```

#### API Layer
```go
// Test with httptest
func TestPriceClient_GetCryptoPrices(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]map[string]float64{
            "bitcoin": {"usd": 45000},
        })
    }))
    defer server.Close()
    
    client := api.NewPriceClientWithURL(server.URL)
    prices, err := client.GetCryptoPrices([]string{"BTC"})
    
    assert.NoError(t, err)
    assert.Equal(t, 45000.0, prices["BTC"])
}
```

#### UI Layer
```go
// Test state transitions
func TestModel_AddAsset(t *testing.T) {
    model := ui.InitialModel()
    model = model.HandleKey("n") // Open add asset modal
    
    assert.Equal(t, ui.ViewAddAsset, model.View)
    assert.True(t, model.InputMode)
}
```

### 2. Integration Tests

```go
func TestAddAssetWorkflow(t *testing.T) {
    // Setup
    db := test.SetupTestDB(t)
    mockAPI := test.SetupMockAPI(t)
    defer mockAPI.Close()
    
    app := NewTestApp(db, mockAPI.URL)
    
    // Execute workflow
    app.PressKey("n")
    app.TypeText("hardware wallet")
    app.PressTab()
    app.TypeText("BTC")
    app.PressTab()
    app.TypeText("0.5")
    app.PressEnter()
    
    // Verify
    holdings := app.GetHoldings()
    assert.Len(t, holdings, 1)
    assert.Equal(t, "BTC", holdings[0].Asset.Symbol)
    assert.Equal(t, 0.5, holdings[0].Amount)
}
```

### 3. Table-Driven Tests

```go
func TestAssetType_Guess(t *testing.T) {
    tests := []struct {
        name     string
        symbol   string
        expected models.AssetType
    }{
        {"Bitcoin", "BTC", models.AssetTypeCrypto},
        {"US Dollar", "USD", models.AssetTypeFiat},
        {"Euro", "EUR", models.AssetTypeFiat},
        {"Unknown", "XYZ", models.AssetTypeCrypto},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := guessAssetType(tt.symbol)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Test Helpers

### Database Helper
```go
package test

func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)
    
    err = db.AutoMigrate(
        &models.Account{},
        &models.Asset{},
        &models.Holding{},
    )
    require.NoError(t, err)
    
    t.Cleanup(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    })
    
    return db
}
```

### API Mock Helper
```go
func SetupMockAPI(t *testing.T) *httptest.Server {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/api/v3/simple/price", func(w http.ResponseWriter, r *http.Request) {
        response := map[string]map[string]float64{
            "bitcoin": {"usd": 45000},
            "ethereum": {"usd": 3000},
        }
        json.NewEncoder(w).Encode(response)
    })
    
    return httptest.NewServer(mux)
}
```

### Test Fixtures
```go
func CreateTestAccount(t *testing.T, db *gorm.DB, name string) *models.Account {
    account := &models.Account{
        Name: name,
        Type: "wallet",
    }
    require.NoError(t, db.Create(account).Error)
    return account
}

func CreateTestHolding(t *testing.T, db *gorm.DB, opts ...HoldingOption) *models.Holding {
    // Builder pattern for test data
}
```

## Mocking Strategy

### Interface-Based Mocking
```go
type PriceClient interface {
    GetCryptoPrices(symbols []string) (map[string]float64, error)
    GetFiatRates(symbols []string) (map[string]float64, error)
}

//go:generate mockery --name=PriceClient --output=mocks
```

### Time Mocking
```go
type Clock interface {
    Now() time.Time
}

type MockClock struct {
    CurrentTime time.Time
}

func (m *MockClock) Now() time.Time {
    return m.CurrentTime
}
```

## Test Commands

```makefile
# Makefile
.PHONY: test test-unit test-integration test-coverage test-race

test: test-unit test-integration

test-unit:
	go test -short ./...

test-integration:
	go test -run Integration ./test/integration

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-race:
	go test -race ./...

test-bench:
	go test -bench=. -benchmem ./...
```

## Continuous Integration

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run tests
        run: make test-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
```

## Testing Best Practices

1. **Test Naming**: `Test<Type>_<Method>_<Scenario>`
2. **Arrange-Act-Assert**: Clear test structure
3. **One Assertion Per Test**: When possible
4. **Parallel Tests**: Use `t.Parallel()` for independent tests
5. **Cleanup**: Always cleanup resources with `t.Cleanup()`
6. **Error Messages**: Provide context in assertions
7. **Test Data**: Use builders and fixtures for complex data

## Performance Testing

```go
func BenchmarkPriceService_FetchPrices(b *testing.B) {
    service := setupTestService(b)
    assets := generateTestAssets(100)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.FetchPrices(assets)
    }
}
```

## UI Testing Strategy

Since Bubble Tea is interactive, we'll use:

1. **State Testing**: Test model state transitions
2. **View Snapshots**: Golden file testing for rendered views
3. **Command Testing**: Verify commands are triggered correctly
4. **Message Handling**: Test update logic for all message types

```go
func TestTableView_Render(t *testing.T) {
    model := setupTestModel(t)
    view := model.View()
    
    golden.Assert(t, view, "table_view.golden")
}
```