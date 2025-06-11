package service

import (
	"time"

	"github.com/bioharz/budget/internal/api"
	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/internal/repository"
	"gorm.io/gorm"
)

type PriceService struct {
	client    *api.PriceClient
	assetRepo *repository.AssetRepository
	cacheRepo *repository.PriceCacheRepository
}

func NewPriceService() *PriceService {
	return &PriceService{
		client:    api.NewPriceClient(),
		assetRepo: repository.NewAssetRepository(),
		cacheRepo: repository.NewPriceCacheRepository(),
	}
}

func NewPriceServiceWithDB(database *gorm.DB) *PriceService {
	return &PriceService{
		client:    api.NewPriceClient(),
		assetRepo: repository.NewAssetRepositoryWithDB(database),
		cacheRepo: repository.NewPriceCacheRepositoryWithDB(database),
	}
}

func (s *PriceService) FetchPrices(assets []models.Asset) (map[uint]float64, error) {
	prices := make(map[uint]float64)

	// Separate crypto and fiat assets
	var cryptoSymbols []string
	var fiatSymbols []string
	cryptoAssetMap := make(map[string]uint)
	fiatAssetMap := make(map[string]uint)

	for _, asset := range assets {
		switch asset.Type {
		case models.AssetTypeCrypto:
			cryptoSymbols = append(cryptoSymbols, asset.Symbol)
			cryptoAssetMap[asset.Symbol] = asset.ID
		case models.AssetTypeFiat:
			fiatSymbols = append(fiatSymbols, asset.Symbol)
			fiatAssetMap[asset.Symbol] = asset.ID
		default:
			// For stocks and others, default to 0 for now
			prices[asset.ID] = 0
		}
	}

	// Fetch crypto prices
	if len(cryptoSymbols) > 0 {
		cryptoPrices, err := s.client.GetCryptoPrices(cryptoSymbols)
		if err != nil {
			// Continue with partial results
		} else {
			for symbol, price := range cryptoPrices {
				if assetID, ok := cryptoAssetMap[symbol]; ok {
					prices[assetID] = price
				}
			}
		}
	}

	// Fetch fiat rates
	if len(fiatSymbols) > 0 {
		fiatRates, err := s.client.GetFiatRates(fiatSymbols)
		if err != nil {
			// Continue with partial results
		} else {
			for symbol, rate := range fiatRates {
				if assetID, ok := fiatAssetMap[symbol]; ok {
					prices[assetID] = rate
				}
			}
		}
	}

	// Save prices to cache
	if s.cacheRepo != nil {
		if err := s.cacheRepo.UpsertBatch(prices); err != nil {
			// Log error but don't fail the operation
			_ = err
			// Log error but don't fail the operation
		}
	}

	return prices, nil
}

// GetCachedPrices returns prices from the cache
func (s *PriceService) GetCachedPrices() (map[uint]float64, error) {
	if s.cacheRepo == nil {
		return make(map[uint]float64), nil
	}
	return s.cacheRepo.GetPricesMap()
}

// GetLastUpdateTime returns when prices were last updated
func (s *PriceService) GetLastUpdateTime() (*time.Time, error) {
	if s.cacheRepo == nil {
		return nil, nil
	}
	return s.cacheRepo.GetLastUpdateTime()
}
