package repository

import (
	"testing"
	"time"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/test/fixtures"
	"github.com/bioharz/budget/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHoldingRepository_Create(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewHoldingRepositoryWithDB(db)
		
		// Create account and asset first
		account := fixtures.NewAccount().WithName("Test Wallet").Create(t, db)
		asset := fixtures.NewAsset().WithSymbol("BTC").Create(t, db)
		
		holding := &models.Holding{
			AccountID:     account.ID,
			AssetID:       asset.ID,
			Amount:        0.5,
			PurchasePrice: 40000,
			PurchaseDate:  time.Now(),
		}
		
		err := repo.Create(holding)
		require.NoError(t, err)
		assert.NotZero(t, holding.ID)
		
		// Verify it was saved
		var saved models.Holding
		err = db.First(&saved, holding.ID).Error
		require.NoError(t, err)
		assert.Equal(t, account.ID, saved.AccountID)
		assert.Equal(t, asset.ID, saved.AssetID)
		assert.Equal(t, 0.5, saved.Amount)
		assert.Equal(t, float64(40000), saved.PurchasePrice)
	})
}

func TestHoldingRepository_GetAll(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewHoldingRepositoryWithDB(db)
		
		// Create test data
		account1 := fixtures.NewAccount().WithName("hardware wallet").Create(t, db)
		account2 := fixtures.NewAccount().WithName("NeoBank").Create(t, db)
		btc := fixtures.NewAsset().WithSymbol("BTC").Create(t, db)
		eth := fixtures.NewAsset().WithSymbol("ETH").Create(t, db)
		
		// Create holdings
		fixtures.NewHolding().
			WithAccount(account1).
			WithAsset(btc).
			WithAmount(0.5).
			Create(t, db)
			
		fixtures.NewHolding().
			WithAccount(account1).
			WithAsset(eth).
			WithAmount(10).
			Create(t, db)
			
		fixtures.NewHolding().
			WithAccount(account2).
			WithAsset(btc).
			WithAmount(0.1).
			Create(t, db)
		
		// Get all holdings
		holdings, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, holdings, 3)
		
		// Verify preloading worked
		for _, holding := range holdings {
			assert.NotEmpty(t, holding.Account.Name)
			assert.NotEmpty(t, holding.Asset.Symbol)
		}
	})
}

func TestHoldingRepository_GetByID(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewHoldingRepositoryWithDB(db)
		
		// Create test data
		account := fixtures.NewAccount().WithName("Test Account").Create(t, db)
		asset := fixtures.NewAsset().WithSymbol("ETH").Create(t, db)
		created := fixtures.NewHolding().
			WithAccount(account).
			WithAsset(asset).
			WithAmount(5.5).
			WithPurchasePrice(2000).
			Create(t, db)
		
		// Get by ID
		holding, err := repo.GetByID(created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.Amount, holding.Amount)
		assert.Equal(t, account.Name, holding.Account.Name)
		assert.Equal(t, asset.Symbol, holding.Asset.Symbol)
		
		// Test non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestHoldingRepository_GetByAccountID(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewHoldingRepositoryWithDB(db)
		
		// Create test data
		account1 := fixtures.NewAccount().WithName("Account 1").Create(t, db)
		account2 := fixtures.NewAccount().WithName("Account 2").Create(t, db)
		btc := fixtures.NewAsset().WithSymbol("BTC").Create(t, db)
		eth := fixtures.NewAsset().WithSymbol("ETH").Create(t, db)
		
		// Create holdings for account1
		fixtures.NewHolding().
			WithAccount(account1).
			WithAsset(btc).
			WithAmount(1).
			Create(t, db)
			
		fixtures.NewHolding().
			WithAccount(account1).
			WithAsset(eth).
			WithAmount(20).
			Create(t, db)
			
		// Create holding for account2
		fixtures.NewHolding().
			WithAccount(account2).
			WithAsset(btc).
			WithAmount(0.5).
			Create(t, db)
		
		// Get holdings for account1
		holdings, err := repo.GetByAccountID(account1.ID)
		require.NoError(t, err)
		assert.Len(t, holdings, 2)
		
		// Verify all holdings belong to account1
		for _, holding := range holdings {
			assert.Equal(t, account1.ID, holding.AccountID)
			assert.Equal(t, account1.Name, holding.Account.Name)
		}
		
		// Get holdings for account2
		holdings, err = repo.GetByAccountID(account2.ID)
		require.NoError(t, err)
		assert.Len(t, holdings, 1)
		assert.Equal(t, account2.ID, holdings[0].AccountID)
	})
}

func TestHoldingRepository_Update(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewHoldingRepositoryWithDB(db)
		
		// Create test data
		account := fixtures.NewAccount().Create(t, db)
		asset := fixtures.NewAsset().Create(t, db)
		holding := fixtures.NewHolding().
			WithAccount(account).
			WithAsset(asset).
			WithAmount(1).
			WithPurchasePrice(30000).
			Create(t, db)
		
		// Update holding
		holding.Amount = 1.5
		holding.PurchasePrice = 35000
		err := repo.Update(holding)
		require.NoError(t, err)
		
		// Verify update
		var updated models.Holding
		err = db.First(&updated, holding.ID).Error
		require.NoError(t, err)
		assert.Equal(t, 1.5, updated.Amount)
		assert.Equal(t, float64(35000), updated.PurchasePrice)
	})
}

func TestHoldingRepository_Delete(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewHoldingRepositoryWithDB(db)
		
		// Create test data
		account := fixtures.NewAccount().Create(t, db)
		asset := fixtures.NewAsset().Create(t, db)
		holding := fixtures.NewHolding().
			WithAccount(account).
			WithAsset(asset).
			Create(t, db)
		
		// Delete holding
		err := repo.Delete(holding.ID)
		require.NoError(t, err)
		
		// Verify soft delete
		var deleted models.Holding
		err = db.First(&deleted, holding.ID).Error
		assert.Error(t, err)
		
		// Verify with unscoped
		err = db.Unscoped().First(&deleted, holding.ID).Error
		require.NoError(t, err)
		assert.NotNil(t, deleted.DeletedAt)
	})
}

func TestHoldingRepository_MultipleHoldingsSameAsset(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewHoldingRepositoryWithDB(db)
		
		// Create test data
		account1 := fixtures.NewAccount().WithName("Account 1").Create(t, db)
		account2 := fixtures.NewAccount().WithName("Account 2").Create(t, db)
		btc := fixtures.NewAsset().WithSymbol("BTC").Create(t, db)
		
		// Create multiple holdings of same asset in different accounts
		holding1 := &models.Holding{
			AccountID: account1.ID,
			AssetID:   btc.ID,
			Amount:    0.5,
		}
		err := repo.Create(holding1)
		require.NoError(t, err)
		
		holding2 := &models.Holding{
			AccountID: account2.ID,
			AssetID:   btc.ID,
			Amount:    1.0,
		}
		err = repo.Create(holding2)
		require.NoError(t, err)
		
		// Verify both exist
		holdings, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, holdings, 2)
	})
}