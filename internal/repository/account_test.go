package repository

import (
	"testing"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/test/fixtures"
	"github.com/bioharz/budget/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAccountRepository_Create(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewAccountRepositoryWithDB(db)

	account := &models.Account{
		Name:  "Test Wallet",
		Type:  "wallet",
		Color: "#FF5733",
	}

	err := repo.Create(account)
	require.NoError(t, err)
	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	// Verify it was saved
	var saved models.Account
	err = db.First(&saved, account.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Test Wallet", saved.Name)
	assert.Equal(t, "wallet", saved.Type)
	assert.Equal(t, "#FF5733", saved.Color)
}

func TestAccountRepository_GetAll(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewAccountRepositoryWithDB(db)

	// Create test accounts
	account1 := fixtures.NewAccount().WithName("hardware wallet").Create(t, db)
	account2 := fixtures.NewAccount().WithName("NeoBank").WithType("bank").Create(t, db)

	// Get all accounts
	accounts, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, accounts, 2)

	// Verify order and content
	assert.Equal(t, account1.Name, accounts[0].Name)
	assert.Equal(t, account2.Name, accounts[1].Name)
}

func TestAccountRepository_GetByID(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewAccountRepositoryWithDB(db)

	// Create test account
	created := fixtures.NewAccount().WithName("Test Account").Create(t, db)

	// Get by ID
	account, err := repo.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.Name, account.Name)
	assert.Equal(t, created.Type, account.Type)

	// Test non-existent ID
	_, err = repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestAccountRepository_GetByName(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewAccountRepositoryWithDB(db)

	// Create test account
	fixtures.NewAccount().WithName("Unique Name").Create(t, db)

	// Get by name
	account, err := repo.GetByName("Unique Name")
	require.NoError(t, err)
	assert.Equal(t, "Unique Name", account.Name)

	// Test non-existent name
	_, err = repo.GetByName("Non Existent")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestAccountRepository_Update(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewAccountRepositoryWithDB(db)

	// Create test account
	account := fixtures.NewAccount().Create(t, db)
	originalName := account.Name

	// Update account
	account.Name = "Updated Name"
	account.Color = "#00FF00"
	err := repo.Update(account)
	require.NoError(t, err)

	// Verify update
	var updated models.Account
	err = db.First(&updated, account.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "#00FF00", updated.Color)
	assert.NotEqual(t, originalName, updated.Name)
}

func TestAccountRepository_Delete(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewAccountRepositoryWithDB(db)

	// Create test account
	account := fixtures.NewAccount().Create(t, db)

	// Delete account
	err := repo.Delete(account.ID)
	require.NoError(t, err)

	// Verify soft delete
	var deleted models.Account
	err = db.First(&deleted, account.ID).Error
	assert.Error(t, err) // Should not find with normal query

	// Verify with unscoped
	err = db.Unscoped().First(&deleted, account.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, deleted.DeletedAt)
}
