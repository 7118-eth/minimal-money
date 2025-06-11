package api

import (
	"testing"
	"time"

	"github.com/bioharz/budget/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceClient_GetCryptoPrices_Real(t *testing.T) {
	helpers.SkipIfShort(t)

	client := NewPriceClient()

	t.Run("fetch single crypto price", func(t *testing.T) {
		helpers.RateLimitDelay()

		prices, err := client.GetCryptoPrices([]string{"BTC"})
		require.NoError(t, err)

		assert.Contains(t, prices, "BTC")
		helpers.AssertReasonablePrice(t, "BTC", prices["BTC"])
		t.Logf("BTC price: $%.2f", prices["BTC"])
	})

	t.Run("fetch multiple crypto prices", func(t *testing.T) {
		helpers.RateLimitDelay()

		symbols := []string{"BTC", "ETH", "SOL"}
		prices, err := client.GetCryptoPrices(symbols)
		require.NoError(t, err)

		// We might not get all symbols if they're not in our mapping
		assert.NotEmpty(t, prices)

		for symbol, price := range prices {
			helpers.AssertReasonablePrice(t, symbol, price)
			t.Logf("%s price: $%.2f", symbol, price)
		}
	})

	t.Run("caching works", func(t *testing.T) {
		// First call
		start := time.Now()
		prices1, err := client.GetCryptoPrices([]string{"BTC"})
		require.NoError(t, err)
		firstCallDuration := time.Since(start)

		// Second call should be cached
		start = time.Now()
		prices2, err := client.GetCryptoPrices([]string{"BTC"})
		require.NoError(t, err)
		cachedCallDuration := time.Since(start)

		// Cached call should be faster (at least not making network call)
		assert.Less(t, cachedCallDuration.Milliseconds(), int64(100), "Cached call should be under 100ms")
		assert.Equal(t, prices1["BTC"], prices2["BTC"])
		t.Logf("First call: %v, Cached call: %v", firstCallDuration, cachedCallDuration)
	})

	t.Run("unknown symbol returns empty", func(t *testing.T) {
		helpers.RateLimitDelay()

		prices, err := client.GetCryptoPrices([]string{"FAKECOIN"})
		require.NoError(t, err)
		assert.Empty(t, prices)
	})
}

func TestPriceClient_GetFiatRates_Real(t *testing.T) {
	helpers.SkipIfShort(t)

	client := NewPriceClient()

	t.Run("USD always returns 1", func(t *testing.T) {
		rates, err := client.GetFiatRates([]string{"USD"})
		require.NoError(t, err)

		assert.Equal(t, 1.0, rates["USD"])
	})

	t.Run("fetch common fiat rates", func(t *testing.T) {
		helpers.RateLimitDelay()

		symbols := []string{"EUR", "GBP", "JPY"}
		rates, err := client.GetFiatRates(symbols)
		require.NoError(t, err)

		assert.Len(t, rates, len(symbols))

		for symbol, rate := range rates {
			helpers.AssertReasonablePrice(t, symbol, rate)
			t.Logf("%s/USD rate: %.4f", symbol, rate)
		}

		// EUR and GBP should be worth more than USD
		assert.Greater(t, rates["EUR"], 0.5)
		assert.Greater(t, rates["GBP"], 0.5)

		// JPY should be worth less than USD
		assert.Less(t, rates["JPY"], 0.1)
	})

	t.Run("caching works for fiat", func(t *testing.T) {
		// First call
		start := time.Now()
		rates1, err := client.GetFiatRates([]string{"EUR"})
		require.NoError(t, err)
		firstCallDuration := time.Since(start)

		// Second call should be cached
		start = time.Now()
		rates2, err := client.GetFiatRates([]string{"EUR"})
		require.NoError(t, err)
		cachedCallDuration := time.Since(start)

		// Cached call should be faster (at least not making network call)
		assert.Less(t, cachedCallDuration.Milliseconds(), int64(100), "Cached call should be under 100ms")
		assert.Equal(t, rates1["EUR"], rates2["EUR"])
		t.Logf("First call: %v, Cached call: %v", firstCallDuration, cachedCallDuration)
	})
}

func TestCryptoIDMapping(t *testing.T) {
	// Just verify our mapping has common cryptos
	commonCryptos := []string{"BTC", "ETH", "USDT", "USDC", "BNB", "SOL"}

	for _, symbol := range commonCryptos {
		_, exists := cryptoIDMapping[symbol]
		assert.True(t, exists, "Missing mapping for %s", symbol)
	}
}

func TestPriceClient_Integration(t *testing.T) {
	helpers.SkipIfShort(t)

	client := NewPriceClient()

	// Simulate a real portfolio fetch
	cryptos := []string{"BTC", "ETH"}
	fiats := []string{"USD", "EUR"}

	helpers.RateLimitDelay()
	cryptoPrices, err := client.GetCryptoPrices(cryptos)
	require.NoError(t, err)

	helpers.RateLimitDelay()
	fiatRates, err := client.GetFiatRates(fiats)
	require.NoError(t, err)

	// Calculate a portfolio value
	btcAmount := 0.5
	ethAmount := 10.0
	eurAmount := 1000.0

	portfolioValue := btcAmount*cryptoPrices["BTC"] +
		ethAmount*cryptoPrices["ETH"] +
		eurAmount*fiatRates["EUR"]

	assert.Greater(t, portfolioValue, 10000.0) // Should be worth something!
	t.Logf("Portfolio value: $%.2f", portfolioValue)
	t.Logf("BTC: %.2f @ $%.2f = $%.2f", btcAmount, cryptoPrices["BTC"], btcAmount*cryptoPrices["BTC"])
	t.Logf("ETH: %.2f @ $%.2f = $%.2f", ethAmount, cryptoPrices["ETH"], ethAmount*cryptoPrices["ETH"])
	t.Logf("EUR: %.2f @ $%.2f = $%.2f", eurAmount, fiatRates["EUR"], eurAmount*fiatRates["EUR"])
}
