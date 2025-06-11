package repository

import (
	"time"

	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PriceCacheRepository struct {
	db *gorm.DB
}

func NewPriceCacheRepository() *PriceCacheRepository {
	return &PriceCacheRepository{
		db: db.DB,
	}
}

func NewPriceCacheRepositoryWithDB(database *gorm.DB) *PriceCacheRepository {
	return &PriceCacheRepository{
		db: database,
	}
}

// GetByAssetID retrieves cached price for a specific asset
func (r *PriceCacheRepository) GetByAssetID(assetID uint) (*models.PriceCache, error) {
	var cache models.PriceCache
	err := r.db.Where("asset_id = ?", assetID).First(&cache).Error
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// GetAll retrieves all cached prices
func (r *PriceCacheRepository) GetAll() ([]models.PriceCache, error) {
	var caches []models.PriceCache
	err := r.db.Preload("Asset").Find(&caches).Error
	return caches, err
}

// GetPricesMap returns a map of asset_id -> price
func (r *PriceCacheRepository) GetPricesMap() (map[uint]float64, error) {
	if r.db == nil {
		return make(map[uint]float64), nil
	}
	var caches []models.PriceCache
	err := r.db.Find(&caches).Error
	if err != nil {
		return nil, err
	}

	priceMap := make(map[uint]float64)
	for _, cache := range caches {
		priceMap[cache.AssetID] = cache.PriceUSD
	}
	return priceMap, nil
}

// Upsert creates or updates a price cache entry
func (r *PriceCacheRepository) Upsert(assetID uint, priceUSD float64) error {
	cache := models.PriceCache{
		AssetID:   assetID,
		PriceUSD:  priceUSD,
		UpdatedAt: time.Now(),
	}

	// Use GORM's Upsert functionality
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "asset_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"price_usd", "updated_at"}),
	}).Create(&cache).Error
}

// UpsertBatch creates or updates multiple price cache entries
func (r *PriceCacheRepository) UpsertBatch(prices map[uint]float64) error {
	if len(prices) == 0 {
		return nil
	}

	now := time.Now()
	caches := make([]models.PriceCache, 0, len(prices))
	for assetID, price := range prices {
		caches = append(caches, models.PriceCache{
			AssetID:   assetID,
			PriceUSD:  price,
			UpdatedAt: now,
		})
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "asset_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"price_usd", "updated_at"}),
	}).CreateInBatches(&caches, 100).Error
}

// GetLastUpdateTime returns the most recent update time
func (r *PriceCacheRepository) GetLastUpdateTime() (*time.Time, error) {
	if r.db == nil {
		return nil, nil
	}
	var cache models.PriceCache
	err := r.db.Order("updated_at desc").First(&cache).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cache.UpdatedAt, nil
}
