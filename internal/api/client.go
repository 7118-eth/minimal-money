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

func (c *PriceClient) GetCryptoPrice(symbol string) (float64, error) {
	symbol = strings.ToLower(symbol)
	
	if cached, ok := c.cache[symbol]; ok {
		if time.Since(cached.timestamp) < 5*time.Minute {
			return cached.price, nil
		}
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", symbol)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch price: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	price, ok := result[symbol]["usd"]
	if !ok {
		return 0, fmt.Errorf("price not found for %s", symbol)
	}

	c.cache[symbol] = cachedPrice{
		price:     price,
		timestamp: time.Now(),
	}

	return price, nil
}

func (c *PriceClient) GetFiatRate(from, to string) (float64, error) {
	return 1.0, nil
}