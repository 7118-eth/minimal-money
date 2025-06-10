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

func TestAssetRepository_Create(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewAssetRepositoryWithDB(db)
		
		asset := &models.Asset{
			Symbol: "BTC",
			Name:   "Bitcoin",
			Type:   models.AssetTypeCrypto,
		}
		
		err := repo.Create(asset)
		require.NoError(t, err)
		assert.NotZero(t, asset.ID)
		
		// Verify it was saved
		var saved models.Asset
		err = db.First(&saved, asset.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "BTC", saved.Symbol)
		assert.Equal(t, "Bitcoin", saved.Name)
		assert.Equal(t, models.AssetTypeCrypto, saved.Type)
	})
}

func TestAssetRepository_GetAll(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewAssetRepositoryWithDB(db)
		
		// Create test assets
		btc := fixtures.NewAsset().WithSymbol("BTC").WithName("Bitcoin").Create(t, db)
		eth := fixtures.NewAsset().WithSymbol("ETH").WithName("Ethereum").Create(t, db)
		usd := fixtures.NewAsset().
			WithSymbol("USD").
			WithName("US Dollar").
			WithType(models.AssetTypeFiat).
			Create(t, db)
		
		// Get all assets
		assets, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, assets, 3)
		
		// Create a map for easier verification
		assetMap := make(map[string]models.Asset)
		for _, asset := range assets {
			assetMap[asset.Symbol] = asset
		}
		
		assert.Equal(t, btc.Name, assetMap["BTC"].Name)
		assert.Equal(t, eth.Name, assetMap["ETH"].Name)
		assert.Equal(t, usd.Type, assetMap["USD"].Type)
	})
}

func TestAssetRepository_GetBySymbol(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewAssetRepositoryWithDB(db)
		
		// Create test asset
		fixtures.NewAsset().WithSymbol("ETH").WithName("Ethereum").Create(t, db)
		
		// Get by symbol
		asset, err := repo.GetBySymbol("ETH")
		require.NoError(t, err)
		assert.Equal(t, "ETH", asset.Symbol)
		assert.Equal(t, "Ethereum", asset.Name)
		
		// Test non-existent symbol
		_, err = repo.GetBySymbol("FAKE")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestAssetRepository_Update(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewAssetRepositoryWithDB(db)
		
		// Create test asset
		asset := fixtures.NewAsset().WithSymbol("SOL").WithName("Solana").Create(t, db)
		
		// Update asset
		asset.Name = "Solana Network"
		err := repo.Update(asset)
		require.NoError(t, err)
		
		// Verify update
		var updated models.Asset
		err = db.First(&updated, asset.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Solana Network", updated.Name)
		assert.Equal(t, "SOL", updated.Symbol) // Symbol unchanged
	})
}

func TestAssetRepository_Delete(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewAssetRepositoryWithDB(db)
		
		// Create test asset
		asset := fixtures.NewAsset().Create(t, db)
		
		// Delete asset
		err := repo.Delete(asset.ID)
		require.NoError(t, err)
		
		// Verify soft delete
		var deleted models.Asset
		err = db.First(&deleted, asset.ID).Error
		assert.Error(t, err)
		
		// Verify with unscoped
		err = db.Unscoped().First(&deleted, asset.ID).Error
		require.NoError(t, err)
		assert.NotNil(t, deleted.DeletedAt)
	})
}

func TestAssetRepository_UniqueSymbol(t *testing.T) {
	helpers.WithTestDB(t, func(db *gorm.DB) {
		repo := NewAssetRepositoryWithDB(db)
		
		// Create first asset
		asset1 := &models.Asset{
			Symbol: "BTC",
			Name:   "Bitcoin",
			Type:   models.AssetTypeCrypto,
		}
		err := repo.Create(asset1)
		require.NoError(t, err)
		
		// Try to create duplicate symbol
		asset2 := &models.Asset{
			Symbol: "BTC",
			Name:   "Bitcoin Copy",
			Type:   models.AssetTypeCrypto,
		}
		err = repo.Create(asset2)
		assert.Error(t, err) // Should fail due to unique constraint
	})
}