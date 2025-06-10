package repository

import (
	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
	"gorm.io/gorm"
)

type AssetRepository struct{
	db *gorm.DB
}

func NewAssetRepository() *AssetRepository {
	return &AssetRepository{db: db.DB}
}

func NewAssetRepositoryWithDB(database *gorm.DB) *AssetRepository {
	return &AssetRepository{db: database}
}

func (r *AssetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

func (r *AssetRepository) GetAll() ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.Find(&assets).Error
	return assets, err
}

func (r *AssetRepository) GetByID(id uint) (models.Asset, error) {
	var asset models.Asset
	err := r.db.First(&asset, id).Error
	return asset, err
}

func (r *AssetRepository) GetBySymbol(symbol string) (models.Asset, error) {
	var asset models.Asset
	err := r.db.Where("symbol = ?", symbol).First(&asset).Error
	return asset, err
}

func (r *AssetRepository) Update(asset *models.Asset) error {
	return r.db.Save(asset).Error
}

func (r *AssetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Asset{}, id).Error
}