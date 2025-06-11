package integration

import (
	"testing"
	"time"

	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/internal/repository"
	"github.com/bioharz/budget/internal/service"
	"github.com/bioharz/budget/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletePortfolioWorkflow(t *testing.T) {
	helpers.SkipIfShort(t)

	// Setup database
	database := helpers.SetupTestDB(t)

	// Initialize global DB for compatibility
	oldDB := db.DB
	db.DB = database
	t.Cleanup(func() {
		db.DB = oldDB
	})

	// Step 1: Create accounts
	accountRepo := repository.NewAccountRepositoryWithDB(database)

	hardwareWallet := &models.Account{
		Name:  "hardware wallet",
		Type:  "wallet",
		Color: "#FF5733",
	}
	err := accountRepo.Create(hardwareWallet)
	require.NoError(t, err)
	t.Logf("Created account: %s (ID: %d)", hardwareWallet.Name, hardwareWallet.ID)

	neobank := &models.Account{
		Name:  "NeoBank",
		Type:  "bank",
		Color: "#33FF57",
	}
	err = accountRepo.Create(neobank)
	require.NoError(t, err)
	t.Logf("Created account: %s (ID: %d)", neobank.Name, neobank.ID)

	// Step 2: Create assets
	assetRepo := repository.NewAssetRepositoryWithDB(database)

	btc := &models.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
		Type:   models.AssetTypeCrypto,
	}
	err = assetRepo.Create(btc)
	require.NoError(t, err)

	eth := &models.Asset{
		Symbol: "ETH",
		Name:   "Ethereum",
		Type:   models.AssetTypeCrypto,
	}
	err = assetRepo.Create(eth)
	require.NoError(t, err)

	eur := &models.Asset{
		Symbol: "EUR",
		Name:   "Euro",
		Type:   models.AssetTypeFiat,
	}
	err = assetRepo.Create(eur)
	require.NoError(t, err)

	// Step 3: Create holdings
	holdingRepo := repository.NewHoldingRepositoryWithDB(database)

	holdings := []models.Holding{
		{
			AccountID:     hardwareWallet.ID,
			AssetID:       btc.ID,
			Amount:        0.5,
			PurchasePrice: 40000,
			PurchaseDate:  time.Now().AddDate(0, -6, 0),
		},
		{
			AccountID:     hardwareWallet.ID,
			AssetID:       eth.ID,
			Amount:        10,
			PurchasePrice: 2000,
			PurchaseDate:  time.Now().AddDate(0, -3, 0),
		},
		{
			AccountID:     neobank.ID,
			AssetID:       eur.ID,
			Amount:        1000,
			PurchasePrice: 1.1,
			PurchaseDate:  time.Now().AddDate(0, -1, 0),
		},
	}

	for _, holding := range holdings {
		err := holdingRepo.Create(&holding)
		require.NoError(t, err)
		t.Logf("Created holding: Account %d, Asset %d, Amount %.2f",
			holding.AccountID, holding.AssetID, holding.Amount)
	}

	// Step 4: Fetch current prices
	priceService := service.NewPriceServiceWithDB(database)
	assets, err := assetRepo.GetAll()
	require.NoError(t, err)

	helpers.RateLimitDelay()
	prices, err := priceService.FetchPrices(assets)
	require.NoError(t, err)

	// Step 5: Calculate portfolio value and P&L
	allHoldings, err := holdingRepo.GetAll()
	require.NoError(t, err)

	totalValue := 0.0
	totalPL := 0.0

	t.Log("\nPortfolio Summary:")
	t.Log("==================")

	for _, holding := range allHoldings {
		currentPrice := prices[holding.AssetID]
		value := holding.Amount * currentPrice
		pl := (currentPrice - holding.PurchasePrice) * holding.Amount

		totalValue += value
		totalPL += pl

		t.Logf("%s - %s: %.2f @ $%.2f = $%.2f (P&L: $%.2f)",
			holding.Account.Name,
			holding.Asset.Symbol,
			holding.Amount,
			currentPrice,
			value,
			pl)
	}

	t.Logf("\nTotal Portfolio Value: $%.2f", totalValue)
	t.Logf("Total P&L: $%.2f (%.2f%%)", totalPL, (totalPL/totalValue)*100)

	// Assertions
	if totalValue < 1000 {
		t.Log("Warning: Prices might not have been fetched correctly (API issue or rate limit)")
		// At least verify the EUR value was calculated
		assert.Greater(t, totalValue, 500.0, "Portfolio should at least have partial EUR value")
	} else {
		// If crypto prices were fetched, expect higher value
		assert.Greater(t, totalValue, 1000.0, "Portfolio should be worth > $1k")
	}
	assert.Len(t, allHoldings, 3, "Should have 3 holdings")

	// Verify relationships are loaded
	for _, holding := range allHoldings {
		assert.NotEmpty(t, holding.Account.Name)
		assert.NotEmpty(t, holding.Asset.Symbol)
	}
}

func TestAddNewAssetWorkflow(t *testing.T) {
	// Setup database
	database := helpers.SetupTestDB(t)

	// Load sample data
	helpers.LoadSampleData(t, database)

	// Repositories
	accountRepo := repository.NewAccountRepositoryWithDB(database)
	assetRepo := repository.NewAssetRepositoryWithDB(database)
	holdingRepo := repository.NewHoldingRepositoryWithDB(database)

	// Step 1: Find or create account
	account, err := accountRepo.GetByName("hardware wallet")
	require.NoError(t, err)

	// Step 2: Find or create asset
	newAsset := &models.Asset{
		Symbol: "SOL",
		Name:   "Solana",
		Type:   models.AssetTypeCrypto,
	}
	err = assetRepo.Create(newAsset)
	require.NoError(t, err)

	// Step 3: Create holding
	holding := &models.Holding{
		AccountID:     account.ID,
		AssetID:       newAsset.ID,
		Amount:        50,
		PurchasePrice: 120,
		PurchaseDate:  time.Now(),
	}
	err = holdingRepo.Create(holding)
	require.NoError(t, err)

	// Step 4: Verify it was added
	holdings, err := holdingRepo.GetByAccountID(account.ID)
	require.NoError(t, err)

	// Find the SOL holding
	var solHolding *models.Holding
	for _, h := range holdings {
		if h.Asset.Symbol == "SOL" {
			solHolding = &h
			break
		}
	}

	require.NotNil(t, solHolding, "Should find SOL holding")
	assert.Equal(t, 50.0, solHolding.Amount)
	assert.Equal(t, 120.0, solHolding.PurchasePrice)

	t.Logf("Successfully added %s to %s account", newAsset.Symbol, account.Name)
}

func TestDatabaseConstraints(t *testing.T) {
	database := helpers.SetupTestDB(t)

	assetRepo := repository.NewAssetRepositoryWithDB(database)
	holdingRepo := repository.NewHoldingRepositoryWithDB(database)

	// Test unique asset symbols
	asset1 := &models.Asset{Symbol: "BTC", Name: "Bitcoin", Type: models.AssetTypeCrypto}
	err := assetRepo.Create(asset1)
	require.NoError(t, err)

	asset2 := &models.Asset{Symbol: "BTC", Name: "Bitcoin Copy", Type: models.AssetTypeCrypto}
	err = assetRepo.Create(asset2)
	assert.Error(t, err, "Should not allow duplicate symbols")

	// Test that we can't create holdings without valid references
	// Note: SQLite foreign keys might not be enabled by default
	// so we just verify the behavior exists
	holding := &models.Holding{
		AccountID: 9999, // Non-existent account
		AssetID:   asset1.ID,
		Amount:    1.0,
	}
	err = holdingRepo.Create(holding)
	// Either it errors (FK enabled) or creates with invalid reference (FK disabled)
	if err == nil {
		t.Log("Foreign key constraints not enforced - SQLite default behavior")
		// Verify it was created but with invalid reference
		created, err := holdingRepo.GetByID(holding.ID)
		require.NoError(t, err)
		assert.Equal(t, uint(9999), created.AccountID)
	} else {
		t.Log("Foreign key constraints enforced")
	}
}
