package service

import (
	"testing"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceService_FetchPrices(t *testing.T) {
	helpers.SkipIfShort(t)

	// Set up test database
	testDB := helpers.SetupTestDB(t)

	service := NewPriceServiceWithDB(testDB)

	t.Run("fetch mixed asset types", func(t *testing.T) {
		helpers.RateLimitDelay()

		assets := []models.Asset{
			{ID: 1, Symbol: "BTC", Type: models.AssetTypeCrypto},
			{ID: 2, Symbol: "ETH", Type: models.AssetTypeCrypto},
			{ID: 3, Symbol: "USD", Type: models.AssetTypeFiat},
			{ID: 4, Symbol: "EUR", Type: models.AssetTypeFiat},
			{ID: 5, Symbol: "AAPL", Type: models.AssetTypeStock}, // Should return 0
		}

		prices, err := service.FetchPrices(assets)
		require.NoError(t, err)

		// Check crypto prices
		assert.Contains(t, prices, uint(1))
		assert.Contains(t, prices, uint(2))
		helpers.AssertReasonablePrice(t, "BTC", prices[1])
		helpers.AssertReasonablePrice(t, "ETH", prices[2])

		// Check fiat prices
		assert.Equal(t, 1.0, prices[3])   // USD should be 1
		assert.Greater(t, prices[4], 0.5) // EUR should be > 0.5

		// Stock should return 0 (not implemented)
		assert.Equal(t, 0.0, prices[5])

		t.Logf("Prices fetched: BTC=$%.2f, ETH=$%.2f, USD=$%.2f, EUR=$%.2f",
			prices[1], prices[2], prices[3], prices[4])
	})

	t.Run("handle empty assets", func(t *testing.T) {
		prices, err := service.FetchPrices([]models.Asset{})
		require.NoError(t, err)
		assert.Empty(t, prices)
	})

	t.Run("handle unknown crypto symbols", func(t *testing.T) {
		helpers.RateLimitDelay()

		assets := []models.Asset{
			{ID: 1, Symbol: "FAKECOIN", Type: models.AssetTypeCrypto},
		}

		prices, err := service.FetchPrices(assets)
		require.NoError(t, err)

		// Unknown crypto should not be in results
		assert.NotContains(t, prices, uint(1))
	})

	t.Run("continue on API errors", func(t *testing.T) {
		// This tests that partial failures don't break everything
		assets := []models.Asset{
			{ID: 1, Symbol: "BTC", Type: models.AssetTypeCrypto},
			{ID: 2, Symbol: "INVALID", Type: models.AssetTypeCrypto},
			{ID: 3, Symbol: "USD", Type: models.AssetTypeFiat},
		}

		prices, err := service.FetchPrices(assets)
		require.NoError(t, err) // Should not error

		// Should have at least USD
		assert.Contains(t, prices, uint(3))
		assert.Equal(t, 1.0, prices[3])
	})
}

func TestPriceService_RealPortfolio(t *testing.T) {
	helpers.SkipIfShort(t)

	// Set up test database
	testDB := helpers.SetupTestDB(t)

	service := NewPriceServiceWithDB(testDB)

	// Simulate a realistic portfolio
	portfolio := []models.Asset{
		{ID: 1, Symbol: "BTC", Name: "Bitcoin", Type: models.AssetTypeCrypto},
		{ID: 2, Symbol: "ETH", Name: "Ethereum", Type: models.AssetTypeCrypto},
		{ID: 3, Symbol: "SOL", Name: "Solana", Type: models.AssetTypeCrypto},
		{ID: 4, Symbol: "USDT", Name: "Tether", Type: models.AssetTypeCrypto},
		{ID: 5, Symbol: "USD", Name: "US Dollar", Type: models.AssetTypeFiat},
		{ID: 6, Symbol: "EUR", Name: "Euro", Type: models.AssetTypeFiat},
		{ID: 7, Symbol: "GBP", Name: "British Pound", Type: models.AssetTypeFiat},
	}

	helpers.RateLimitDelay()
	prices, err := service.FetchPrices(portfolio)
	require.NoError(t, err)

	// Verify we got most prices
	assert.GreaterOrEqual(t, len(prices), 5, "Should fetch most prices")

	// Log portfolio values
	t.Log("Portfolio Prices:")
	for _, asset := range portfolio {
		if price, ok := prices[asset.ID]; ok {
			t.Logf("  %s (%s): $%.2f", asset.Name, asset.Symbol, price)
		} else {
			t.Logf("  %s (%s): No price available", asset.Name, asset.Symbol)
		}
	}

	// Calculate example portfolio value
	holdings := map[uint]float64{
		1: 0.5,    // 0.5 BTC
		2: 10.0,   // 10 ETH
		3: 100.0,  // 100 SOL
		4: 1000.0, // 1000 USDT
		5: 5000.0, // 5000 USD
		6: 2000.0, // 2000 EUR
	}

	totalValue := 0.0
	for assetID, amount := range holdings {
		if price, ok := prices[assetID]; ok {
			value := amount * price
			totalValue += value

			// Find asset name
			assetName := ""
			for _, a := range portfolio {
				if a.ID == assetID {
					assetName = a.Symbol
					break
				}
			}
			t.Logf("  %s: %.2f × $%.2f = $%.2f", assetName, amount, price, value)
		}
	}

	t.Logf("Total Portfolio Value: $%.2f", totalValue)
	assert.Greater(t, totalValue, 50000.0, "Portfolio should be worth > $50k")
}
