package repository

import (
	"time"

	"github.com/bioharz/budget/internal/models"
	"gorm.io/gorm"
)

type PriceHistoryRepository struct {
	db *gorm.DB
}

func NewPriceHistoryRepository(db *gorm.DB) *PriceHistoryRepository {
	return &PriceHistoryRepository{db: db}
}

func (r *PriceHistoryRepository) Create(history *models.PriceHistory) error {
	return r.db.Create(history).Error
}

func (r *PriceHistoryRepository) GetByAssetID(assetID uint, limit int) ([]models.PriceHistory, error) {
	var histories []models.PriceHistory
	err := r.db.Where("asset_id = ?", assetID).
		Order("timestamp desc").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *PriceHistoryRepository) GetByAssetIDAndTimeRange(assetID uint, start, end time.Time) ([]models.PriceHistory, error) {
	var histories []models.PriceHistory
	err := r.db.Where("asset_id = ? AND timestamp BETWEEN ? AND ?", assetID, start, end).
		Order("timestamp asc").
		Find(&histories).Error
	return histories, err
}

func (r *PriceHistoryRepository) GetLatestByAssetID(assetID uint) (*models.PriceHistory, error) {
	var history models.PriceHistory
	err := r.db.Where("asset_id = ?", assetID).
		Order("timestamp desc").
		First(&history).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &history, err
}

func (r *PriceHistoryRepository) GetAllAssetHistories(limit int) (map[uint][]models.PriceHistory, error) {
	// Get distinct asset IDs
	var assetIDs []uint
	err := r.db.Model(&models.PriceHistory{}).
		Distinct("asset_id").
		Pluck("asset_id", &assetIDs).Error
	if err != nil {
		return nil, err
	}
	
	result := make(map[uint][]models.PriceHistory)
	
	// Get history for each asset
	for _, assetID := range assetIDs {
		assetHistories, err := r.GetByAssetID(assetID, limit)
		if err != nil {
			return nil, err
		}
		result[assetID] = assetHistories
	}
	
	return result, nil
}