package repository

import (
	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
	"gorm.io/gorm"
)

type AccountRepository struct{
	db *gorm.DB
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{db: db.DB}
}

func NewAccountRepositoryWithDB(database *gorm.DB) *AccountRepository {
	return &AccountRepository{db: database}
}

func (r *AccountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *AccountRepository) GetAll() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) GetByID(id uint) (models.Account, error) {
	var account models.Account
	err := r.db.First(&account, id).Error
	return account, err
}

func (r *AccountRepository) GetByName(name string) (models.Account, error) {
	var account models.Account
	err := r.db.Where("name = ?", name).First(&account).Error
	return account, err
}

func (r *AccountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *AccountRepository) Delete(id uint) error {
	return r.db.Delete(&models.Account{}, id).Error
}