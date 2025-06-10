package repository

import (
	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
)

type AssetRepository struct{}

func NewAssetRepository() *AssetRepository {
	return &AssetRepository{}
}

func (r *AssetRepository) Create(asset *models.Asset) error {
	return db.DB.Create(asset).Error
}

func (r *AssetRepository) GetAll() ([]models.Asset, error) {
	var assets []models.Asset
	err := db.DB.Find(&assets).Error
	return assets, err
}

func (r *AssetRepository) GetByID(id uint) (models.Asset, error) {
	var asset models.Asset
	err := db.DB.First(&asset, id).Error
	return asset, err
}

func (r *AssetRepository) GetBySymbol(symbol string) (models.Asset, error) {
	var asset models.Asset
	err := db.DB.Where("symbol = ?", symbol).First(&asset).Error
	return asset, err
}

func (r *AssetRepository) Update(asset *models.Asset) error {
	return db.DB.Save(asset).Error
}

func (r *AssetRepository) Delete(id uint) error {
	return db.DB.Delete(&models.Asset{}, id).Error
}