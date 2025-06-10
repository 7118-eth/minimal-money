package helpers

import (
	"testing"

	"github.com/bioharz/budget/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
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

	// Cleanup when test finishes
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	})

	return db
}

// TruncateAllTables clears all data from the database
func TruncateAllTables(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		tables := []interface{}{
			&models.Holding{},
			&models.PriceHistory{},
			&models.Asset{},
			&models.Account{},
			&models.PortfolioSnapshot{},
		}

		for _, table := range tables {
			if err := tx.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// WithTestDB runs a test function with a test database
func WithTestDB(t *testing.T, fn func(*gorm.DB)) {
	db := SetupTestDB(t)
	fn(db)
}