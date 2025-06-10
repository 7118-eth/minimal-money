package repository

import (
	"testing"
	"time"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceCacheRepository_Upsert(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceCacheRepositoryWithDB(db)
	assetRepo := NewAssetRepositoryWithDB(db)

	// Create test asset
	asset := models.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
		Type:   models.AssetTypeCrypto,
	}
	err := assetRepo.Create(&asset)
	require.NoError(t, err)

	// Test insert
	err = repo.Upsert(asset.ID, 50000.0)
	require.NoError(t, err)

	// Verify insert
	cache, err := repo.GetByAssetID(asset.ID)
	require.NoError(t, err)
	assert.Equal(t, 50000.0, cache.PriceUSD)
	assert.Equal(t, asset.ID, cache.AssetID)
	assert.WithinDuration(t, time.Now(), cache.UpdatedAt, 2*time.Second)

	// Test update
	time.Sleep(100 * time.Millisecond) // Ensure different timestamp
	err = repo.Upsert(asset.ID, 51000.0)
	require.NoError(t, err)

	// Verify update
	cache, err = repo.GetByAssetID(asset.ID)
	require.NoError(t, err)
	assert.Equal(t, 51000.0, cache.PriceUSD)
	assert.Equal(t, asset.ID, cache.AssetID)
}

func TestPriceCacheRepository_UpsertBatch(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceCacheRepositoryWithDB(db)
	assetRepo := NewAssetRepositoryWithDB(db)

	// Create test assets
	btc := models.Asset{Symbol: "BTC", Name: "Bitcoin", Type: models.AssetTypeCrypto}
	eth := models.Asset{Symbol: "ETH", Name: "Ethereum", Type: models.AssetTypeCrypto}
	usd := models.Asset{Symbol: "USD", Name: "US Dollar", Type: models.AssetTypeFiat}

	require.NoError(t, assetRepo.Create(&btc))
	require.NoError(t, assetRepo.Create(&eth))
	require.NoError(t, assetRepo.Create(&usd))

	// Test batch upsert
	prices := map[uint]float64{
		btc.ID: 50000.0,
		eth.ID: 3000.0,
		usd.ID: 1.0,
	}

	err := repo.UpsertBatch(prices)
	require.NoError(t, err)

	// Verify all prices
	caches, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, caches, 3)

	priceMap, err := repo.GetPricesMap()
	require.NoError(t, err)
	assert.Equal(t, 50000.0, priceMap[btc.ID])
	assert.Equal(t, 3000.0, priceMap[eth.ID])
	assert.Equal(t, 1.0, priceMap[usd.ID])
}

func TestPriceCacheRepository_GetLastUpdateTime(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceCacheRepositoryWithDB(db)
	assetRepo := NewAssetRepositoryWithDB(db)

	// Test with no data
	lastUpdate, err := repo.GetLastUpdateTime()
	require.NoError(t, err)
	assert.Nil(t, lastUpdate)

	// Create test assets
	btc := models.Asset{Symbol: "BTC", Name: "Bitcoin", Type: models.AssetTypeCrypto}
	eth := models.Asset{Symbol: "ETH", Name: "Ethereum", Type: models.AssetTypeCrypto}
	require.NoError(t, assetRepo.Create(&btc))
	require.NoError(t, assetRepo.Create(&eth))

	// Add prices with different times
	err = repo.Upsert(btc.ID, 50000.0)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	
	err = repo.Upsert(eth.ID, 3000.0)
	require.NoError(t, err)

	// Get last update time
	lastUpdate, err = repo.GetLastUpdateTime()
	require.NoError(t, err)
	require.NotNil(t, lastUpdate)
	
	// Should be close to now (ETH was updated last)
	assert.WithinDuration(t, time.Now(), *lastUpdate, 2*time.Second)
}

func TestPriceCacheRepository_GetPricesMap_Empty(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceCacheRepositoryWithDB(db)

	priceMap, err := repo.GetPricesMap()
	require.NoError(t, err)
	assert.NotNil(t, priceMap)
	assert.Len(t, priceMap, 0)
}

func TestPriceCacheRepository_GetByAssetID_NotFound(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceCacheRepositoryWithDB(db)

	cache, err := repo.GetByAssetID(999)
	assert.Error(t, err)
	assert.Nil(t, cache)
}