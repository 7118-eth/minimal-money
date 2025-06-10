package repository

import (
	"testing"
	"time"

	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/test/fixtures"
	"github.com/bioharz/budget/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceHistoryRepository_Create(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceHistoryRepository(db)
	assetRepo := NewAssetRepositoryWithDB(db)
	
	// Create an asset first
	asset := fixtures.NewAsset().WithSymbol("BTC").Build()
	err := assetRepo.Create(&asset)
	require.NoError(t, err)
	
	// Create price history
	history := &models.PriceHistory{
		AssetID:   asset.ID,
		PriceUSD:  50000.00,
		Timestamp: time.Now(),
	}
	
	err = repo.Create(history)
	require.NoError(t, err)
	assert.NotZero(t, history.ID)
	
	// Verify it was saved
	var saved models.PriceHistory
	err = db.First(&saved, history.ID).Error
	require.NoError(t, err)
	assert.Equal(t, asset.ID, saved.AssetID)
	assert.Equal(t, 50000.00, saved.PriceUSD)
	assert.WithinDuration(t, history.Timestamp, saved.Timestamp, time.Second)
}

func TestPriceHistoryRepository_GetByAssetID(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceHistoryRepository(db)
	assetRepo := NewAssetRepositoryWithDB(db)
	
	// Create assets
	btc := fixtures.NewAsset().WithSymbol("BTC").Build()
	eth := fixtures.NewAsset().WithSymbol("ETH").Build()
	err := assetRepo.Create(&btc)
	require.NoError(t, err)
	err = assetRepo.Create(&eth)
	require.NoError(t, err)
	
	// Create price histories
	now := time.Now()
	histories := []models.PriceHistory{
		{AssetID: btc.ID, PriceUSD: 50000.00, Timestamp: now.Add(-2 * time.Hour)},
		{AssetID: btc.ID, PriceUSD: 51000.00, Timestamp: now.Add(-1 * time.Hour)},
		{AssetID: btc.ID, PriceUSD: 52000.00, Timestamp: now},
		{AssetID: eth.ID, PriceUSD: 3000.00, Timestamp: now},
	}
	
	for i := range histories {
		err = repo.Create(&histories[i])
		require.NoError(t, err)
	}
	
	// Get BTC history with limit
	btcHistories, err := repo.GetByAssetID(btc.ID, 2)
	require.NoError(t, err)
	assert.Len(t, btcHistories, 2)
	
	// Should be ordered by timestamp desc
	assert.Equal(t, 52000.00, btcHistories[0].PriceUSD)
	assert.Equal(t, 51000.00, btcHistories[1].PriceUSD)
}

func TestPriceHistoryRepository_GetLatestByAssetID(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceHistoryRepository(db)
	assetRepo := NewAssetRepositoryWithDB(db)
	
	// Create asset
	asset := fixtures.NewAsset().WithSymbol("BTC").Build()
	err := assetRepo.Create(&asset)
	require.NoError(t, err)
	
	// No history yet
	latest, err := repo.GetLatestByAssetID(asset.ID)
	require.NoError(t, err)
	assert.Nil(t, latest)
	
	// Create price histories
	now := time.Now()
	histories := []models.PriceHistory{
		{AssetID: asset.ID, PriceUSD: 50000.00, Timestamp: now.Add(-2 * time.Hour)},
		{AssetID: asset.ID, PriceUSD: 51000.00, Timestamp: now.Add(-1 * time.Hour)},
		{AssetID: asset.ID, PriceUSD: 52000.00, Timestamp: now},
	}
	
	for i := range histories {
		err = repo.Create(&histories[i])
		require.NoError(t, err)
	}
	
	// Get latest
	latest, err = repo.GetLatestByAssetID(asset.ID)
	require.NoError(t, err)
	require.NotNil(t, latest)
	assert.Equal(t, 52000.00, latest.PriceUSD)
}

func TestPriceHistoryRepository_GetAllAssetHistories(t *testing.T) {
	db := helpers.SetupTestDB(t)
	repo := NewPriceHistoryRepository(db)
	assetRepo := NewAssetRepositoryWithDB(db)
	
	// Create assets
	btc := fixtures.NewAsset().WithSymbol("BTC").Build()
	eth := fixtures.NewAsset().WithSymbol("ETH").Build()
	err := assetRepo.Create(&btc)
	require.NoError(t, err)
	err = assetRepo.Create(&eth)
	require.NoError(t, err)
	
	// Create price histories
	now := time.Now()
	histories := []models.PriceHistory{
		{AssetID: btc.ID, PriceUSD: 50000.00, Timestamp: now.Add(-2 * time.Hour)},
		{AssetID: btc.ID, PriceUSD: 51000.00, Timestamp: now.Add(-1 * time.Hour)},
		{AssetID: btc.ID, PriceUSD: 52000.00, Timestamp: now},
		{AssetID: eth.ID, PriceUSD: 3000.00, Timestamp: now.Add(-1 * time.Hour)},
		{AssetID: eth.ID, PriceUSD: 3100.00, Timestamp: now},
	}
	
	for i := range histories {
		err = repo.Create(&histories[i])
		require.NoError(t, err)
	}
	
	// Get all histories with limit
	allHistories, err := repo.GetAllAssetHistories(2)
	require.NoError(t, err)
	assert.Len(t, allHistories, 2) // 2 assets
	
	// Check BTC histories
	assert.Len(t, allHistories[btc.ID], 2)
	assert.Equal(t, 52000.00, allHistories[btc.ID][0].PriceUSD)
	assert.Equal(t, 51000.00, allHistories[btc.ID][1].PriceUSD)
	
	// Check ETH histories
	assert.Len(t, allHistories[eth.ID], 2)
	assert.Equal(t, 3100.00, allHistories[eth.ID][0].PriceUSD)
	assert.Equal(t, 3000.00, allHistories[eth.ID][1].PriceUSD)
}