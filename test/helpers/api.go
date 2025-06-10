package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockAPIServer creates a mock HTTP server for testing API calls
type MockAPIServer struct {
	*httptest.Server
	t *testing.T
}

// NewMockAPIServer creates a new mock API server
func NewMockAPIServer(t *testing.T) *MockAPIServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mock := &MockAPIServer{
		Server: server,
		t:      t,
	}

	// Setup default routes
	mock.setupRoutes(mux)

	t.Cleanup(func() {
		server.Close()
	})

	return mock
}

func (m *MockAPIServer) setupRoutes(mux *http.ServeMux) {
	// CoinGecko price endpoint
	mux.HandleFunc("/api/v3/simple/price", func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		response := make(map[string]map[string]float64)

		// Mock prices based on IDs
		priceMap := map[string]float64{
			"bitcoin":   45000,
			"ethereum":  3000,
			"tether":    1,
			"usd-coin":  1,
			"solana":    150,
			"cardano":   0.5,
		}

		for id, price := range priceMap {
			if ids == "" || contains(ids, id) {
				response[id] = map[string]float64{"usd": price}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// ExchangeRate-API endpoint
	mux.HandleFunc("/v4/latest/USD", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Base  string             `json:"base"`
			Date  string             `json:"date"`
			Rates map[string]float64 `json:"rates"`
		}{
			Base: "USD",
			Date: "2025-01-06",
			Rates: map[string]float64{
				"EUR": 0.92,
				"GBP": 0.79,
				"JPY": 148.50,
				"CHF": 0.85,
				"CAD": 1.35,
				"AUD": 1.52,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

// contains checks if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr || len(str) > len(substr) && contains(str[1:], substr)
}

// MockHTTPClient creates a mock HTTP client with predefined responses
type MockHTTPClient struct {
	Responses map[string]MockResponse
	t         *testing.T
}

type MockResponse struct {
	StatusCode int
	Body       interface{}
	Error      error
}

func NewMockHTTPClient(t *testing.T) *MockHTTPClient {
	return &MockHTTPClient{
		Responses: make(map[string]MockResponse),
		t:         t,
	}
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	if resp, ok := c.Responses[url]; ok {
		if resp.Error != nil {
			return nil, resp.Error
		}
		// Create response with body
		// Implementation depends on needs
	}
	c.t.Fatalf("Unexpected request to: %s", url)
	return nil, nil
}