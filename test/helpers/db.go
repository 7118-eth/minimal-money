package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bioharz/budget/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates a real SQLite database file for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	// Create a unique test database for each test
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	dbPath := filepath.Join("test", "testdata", fmt.Sprintf("test_%s.db", testName))
	
	// Ensure directory exists
	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	require.NoError(t, err)
	
	// Remove any existing test database
	os.Remove(dbPath)
	
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(
		&models.Account{},
		&models.Asset{},
		&models.Holding{},
		&models.PriceHistory{},
		&models.PortfolioSnapshot{},
	)
	require.NoError(t, err)

	// Log database path for debugging
	t.Logf("Using test database: %s", dbPath)

	// Cleanup when test finishes
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
		
		// Keep database if TEST_KEEP_DB is set (for debugging)
		if os.Getenv("TEST_KEEP_DB") != "1" {
			os.Remove(dbPath)
		} else {
			t.Logf("Keeping test database for debugging: %s", dbPath)
		}
	})

	return db
}

// CleanTestData removes all test databases
func CleanTestData() error {
	testDataDir := filepath.Join("test", "testdata")
	files, err := filepath.Glob(filepath.Join(testDataDir, "test_*.db"))
	if err != nil {
		return err
	}
	
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	
	return nil
}

// LoadSampleData creates realistic test data
func LoadSampleData(t *testing.T, db *gorm.DB) {
	// Create accounts
	accounts := []models.Account{
		{Name: "hardware wallet", Type: "wallet", Color: "#FF5733"},
		{Name: "NeoBank", Type: "bank", Color: "#33FF57"},
		{Name: "Binance", Type: "exchange", Color: "#3357FF"},
	}
	
	for i := range accounts {
		require.NoError(t, db.Create(&accounts[i]).Error)
	}
	
	// Create common assets
	assets := []models.Asset{
		{Symbol: "BTC", Name: "Bitcoin", Type: models.AssetTypeCrypto},
		{Symbol: "ETH", Name: "Ethereum", Type: models.AssetTypeCrypto},
		{Symbol: "USD", Name: "US Dollar", Type: models.AssetTypeFiat},
		{Symbol: "EUR", Name: "Euro", Type: models.AssetTypeFiat},
	}
	
	for i := range assets {
		require.NoError(t, db.Create(&assets[i]).Error)
	}
	
	t.Logf("Created %d accounts and %d assets", len(accounts), len(assets))
}