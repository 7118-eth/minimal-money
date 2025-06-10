package repository

import (
	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
)

type HoldingRepository struct{}

func NewHoldingRepository() *HoldingRepository {
	return &HoldingRepository{}
}

func (r *HoldingRepository) Create(holding *models.Holding) error {
	return db.DB.Create(holding).Error
}

func (r *HoldingRepository) GetAll() ([]models.Holding, error) {
	var holdings []models.Holding
	err := db.DB.Preload("Account").Preload("Asset").Find(&holdings).Error
	return holdings, err
}

func (r *HoldingRepository) GetByID(id uint) (models.Holding, error) {
	var holding models.Holding
	err := db.DB.Preload("Account").Preload("Asset").First(&holding, id).Error
	return holding, err
}

func (r *HoldingRepository) GetByAccountID(accountID uint) ([]models.Holding, error) {
	var holdings []models.Holding
	err := db.DB.Preload("Account").Preload("Asset").Where("account_id = ?", accountID).Find(&holdings).Error
	return holdings, err
}

func (r *HoldingRepository) Update(holding *models.Holding) error {
	return db.DB.Save(holding).Error
}

func (r *HoldingRepository) Delete(id uint) error {
	return db.DB.Delete(&models.Holding{}, id).Error
}