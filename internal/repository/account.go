package repository

import (
	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
)

type AccountRepository struct{}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (r *AccountRepository) Create(account *models.Account) error {
	return db.DB.Create(account).Error
}

func (r *AccountRepository) GetAll() ([]models.Account, error) {
	var accounts []models.Account
	err := db.DB.Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) GetByID(id uint) (models.Account, error) {
	var account models.Account
	err := db.DB.First(&account, id).Error
	return account, err
}

func (r *AccountRepository) GetByName(name string) (models.Account, error) {
	var account models.Account
	err := db.DB.Where("name = ?", name).First(&account).Error
	return account, err
}

func (r *AccountRepository) Update(account *models.Account) error {
	return db.DB.Save(account).Error
}

func (r *AccountRepository) Delete(id uint) error {
	return db.DB.Delete(&models.Account{}, id).Error
}