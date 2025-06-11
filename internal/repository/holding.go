package repository

import (
	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
	"gorm.io/gorm"
)

type HoldingRepository struct {
	db *gorm.DB
}

func NewHoldingRepository() *HoldingRepository {
	return &HoldingRepository{db: db.DB}
}

func NewHoldingRepositoryWithDB(database *gorm.DB) *HoldingRepository {
	return &HoldingRepository{db: database}
}

func (r *HoldingRepository) Create(holding *models.Holding) error {
	return r.db.Create(holding).Error
}

func (r *HoldingRepository) GetAll() ([]models.Holding, error) {
	var holdings []models.Holding
	err := r.db.Preload("Account").Preload("Asset").Find(&holdings).Error
	return holdings, err
}

func (r *HoldingRepository) GetByID(id uint) (models.Holding, error) {
	var holding models.Holding
	err := r.db.Preload("Account").Preload("Asset").First(&holding, id).Error
	return holding, err
}

func (r *HoldingRepository) GetByAccountID(accountID uint) ([]models.Holding, error) {
	var holdings []models.Holding
	err := r.db.Preload("Account").Preload("Asset").Where("account_id = ?", accountID).Find(&holdings).Error
	return holdings, err
}

func (r *HoldingRepository) Update(holding *models.Holding) error {
	return r.db.Save(holding).Error
}

func (r *HoldingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Holding{}, id).Error
}
