package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PriceClient struct {
	httpClient *http.Client
	cache      map[string]cachedPrice
}

type cachedPrice struct {
	price     float64
	timestamp time.Time
}

func NewPriceClient() *PriceClient {
	return &PriceClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: make(map[string]cachedPrice),
	}
}

// CoinGecko ID mapping for common crypto symbols
var cryptoIDMapping = map[string]string{
	"BTC": "bitcoin",
	"ETH": "ethereum",
	"USDT": "tether",
	"USDC": "usd-coin",
	"BNB": "binancecoin",
	"XRP": "ripple",
	"SOL": "solana",
	"ADA": "cardano",
	"DOGE": "dogecoin",
	"DOT": "polkadot",
	"MATIC": "matic-network",
	"AVAX": "avalanche-2",
}

func (c *PriceClient) GetCryptoPrices(symbols []string) (map[string]float64, error) {
	prices := make(map[string]float64)
	var idsToFetch []string
	var symbolMap = make(map[string]string) // maps coingecko ID to original symbol
	
	// Check cache first
	for _, symbol := range symbols {
		symbol = strings.ToUpper(symbol)
		if cached, ok := c.cache[symbol]; ok {
			if time.Since(cached.timestamp) < 5*time.Minute {
				prices[symbol] = cached.price
				continue
			}
		}
		
		// Get CoinGecko ID
		if id, ok := cryptoIDMapping[symbol]; ok {
			idsToFetch = append(idsToFetch, id)
			symbolMap[id] = symbol
		}
	}
	
	if len(idsToFetch) == 0 {
		return prices, nil
	}
	
	// Batch fetch from CoinGecko
	ids := strings.Join(idsToFetch, ",")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", ids)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return prices, fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return prices, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map results back to symbols and update cache
	for id, priceData := range result {
		if symbol, ok := symbolMap[id]; ok {
			if price, ok := priceData["usd"]; ok {
				prices[symbol] = price
				c.cache[symbol] = cachedPrice{
					price:     price,
					timestamp: time.Now(),
				}
			}
		}
	}

	return prices, nil
}

func (c *PriceClient) GetFiatRates(symbols []string) (map[string]float64, error) {
	rates := make(map[string]float64)
	
	// USD is always 1.0
	for _, symbol := range symbols {
		if strings.ToUpper(symbol) == "USD" {
			rates[symbol] = 1.0
			continue
		}
	}
	
	// Check if we need to fetch any rates
	var needsFetch bool
	for _, symbol := range symbols {
		symbol = strings.ToUpper(symbol)
		if symbol != "USD" {
			if cached, ok := c.cache[symbol]; ok {
				if time.Since(cached.timestamp) < 1*time.Hour {
					rates[symbol] = cached.price
				} else {
					needsFetch = true
				}
			} else {
				needsFetch = true
			}
		}
	}
	
	if !needsFetch {
		return rates, nil
	}
	
	// Fetch from ExchangeRate-API
	url := "https://api.exchangerate-api.com/v4/latest/USD"
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return rates, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return rates, fmt.Errorf("failed to decode response: %w", err)
	}

	// Update rates and cache
	for _, symbol := range symbols {
		symbol = strings.ToUpper(symbol)
		if rate, ok := result.Rates[symbol]; ok {
			rates[symbol] = 1.0 / rate // Convert to USD rate
			c.cache[symbol] = cachedPrice{
				price:     1.0 / rate,
				timestamp: time.Now(),
			}
		}
	}

	return rates, nil
}