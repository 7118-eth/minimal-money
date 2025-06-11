package helpers

import (
	"testing"
	"time"
)

// RateLimitDelay adds a delay between API tests to respect rate limits
func RateLimitDelay() {
	time.Sleep(100 * time.Millisecond)
}

// SkipIfShort skips API tests when running in short mode
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}
}

// AssertReasonablePrice checks if a price is within reasonable bounds
func AssertReasonablePrice(t *testing.T, symbol string, price float64) {
	t.Helper()

	// Define reasonable price ranges for common assets
	minPrices := map[string]float64{
		"BTC":  10000.0, // Bitcoin should be > $10k
		"ETH":  500.0,   // Ethereum should be > $500
		"USD":  0.9,     // USD rate should be near 1
		"EUR":  0.5,     // EUR rate should be > 0.5
		"USDT": 0.9,     // Stablecoins near $1
		"USDC": 0.9,
	}

	maxPrices := map[string]float64{
		"BTC":  1000000.0, // Bitcoin < $1M (for now!)
		"ETH":  100000.0,  // Ethereum < $100k
		"USD":  1.1,       // USD rate near 1
		"EUR":  2.0,       // EUR rate < 2
		"USDT": 1.1,       // Stablecoins near $1
		"USDC": 1.1,
	}

	if min, ok := minPrices[symbol]; ok {
		if price < min {
			t.Errorf("%s price %.2f is below reasonable minimum %.2f", symbol, price, min)
		}
	}

	if max, ok := maxPrices[symbol]; ok {
		if price > max {
			t.Errorf("%s price %.2f is above reasonable maximum %.2f", symbol, price, max)
		}
	}

	// All prices should be positive
	if price <= 0 {
		t.Errorf("%s price %.2f should be positive", symbol, price)
	}
}
