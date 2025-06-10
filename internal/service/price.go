package service

import (
	"time"

	"github.com/bioharz/budget/internal/api"
	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/internal/repository"
	"gorm.io/gorm"
)

type PriceService struct {
	client         *api.PriceClient
	assetRepo      *repository.AssetRepository
	priceHistoryRepo *repository.PriceHistoryRepository
}

func NewPriceService() *PriceService {
	return &PriceService{
		client:           api.NewPriceClient(),
		assetRepo:        repository.NewAssetRepository(),
		priceHistoryRepo: repository.NewPriceHistoryRepository(db.DB),
	}
}

func NewPriceServiceWithDB(database *gorm.DB) *PriceService {
	return &PriceService{
		client:           api.NewPriceClient(),
		assetRepo:        repository.NewAssetRepositoryWithDB(database),
		priceHistoryRepo: repository.NewPriceHistoryRepository(database),
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
					// Save price history
					s.UpdatePriceHistory(assetID, price)
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
					// Save price history
					s.UpdatePriceHistory(assetID, rate)
				}
			}
		}
	}
	
	return prices, nil
}

func (s *PriceService) UpdatePriceHistory(assetID uint, price float64) error {
	priceHistory := &models.PriceHistory{
		AssetID:   assetID,
		PriceUSD:  price,
		Timestamp: time.Now(),
	}
	
	return s.priceHistoryRepo.Create(priceHistory)
}

func (s *PriceService) GetPriceHistory(assetID uint, limit int) ([]models.PriceHistory, error) {
	return s.priceHistoryRepo.GetByAssetID(assetID, limit)
}

func (s *PriceService) GetAllPriceHistories(limit int) (map[uint][]models.PriceHistory, error) {
	return s.priceHistoryRepo.GetAllAssetHistories(limit)
}