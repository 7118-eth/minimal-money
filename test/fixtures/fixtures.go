package fixtures

import (
	"testing"
	"time"

	"github.com/bioharz/budget/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// AccountBuilder helps create test accounts
type AccountBuilder struct {
	account models.Account
}

func NewAccount() *AccountBuilder {
	return &AccountBuilder{
		account: models.Account{
			Name: "Test Account",
			Type: "wallet",
		},
	}
}

func (b *AccountBuilder) WithName(name string) *AccountBuilder {
	b.account.Name = name
	return b
}

func (b *AccountBuilder) WithType(accountType string) *AccountBuilder {
	b.account.Type = accountType
	return b
}

func (b *AccountBuilder) WithColor(color string) *AccountBuilder {
	b.account.Color = color
	return b
}

func (b *AccountBuilder) Build() models.Account {
	return b.account
}

func (b *AccountBuilder) Create(t *testing.T, db *gorm.DB) *models.Account {
	account := b.Build()
	require.NoError(t, db.Create(&account).Error)
	return &account
}

// AssetBuilder helps create test assets
type AssetBuilder struct {
	asset models.Asset
}

func NewAsset() *AssetBuilder {
	return &AssetBuilder{
		asset: models.Asset{
			Symbol: "BTC",
			Name:   "Bitcoin",
			Type:   models.AssetTypeCrypto,
		},
	}
}

func (b *AssetBuilder) WithSymbol(symbol string) *AssetBuilder {
	b.asset.Symbol = symbol
	return b
}

func (b *AssetBuilder) WithName(name string) *AssetBuilder {
	b.asset.Name = name
	return b
}

func (b *AssetBuilder) WithType(assetType models.AssetType) *AssetBuilder {
	b.asset.Type = assetType
	return b
}

func (b *AssetBuilder) Build() models.Asset {
	return b.asset
}

func (b *AssetBuilder) Create(t *testing.T, db *gorm.DB) *models.Asset {
	asset := b.Build()
	require.NoError(t, db.Create(&asset).Error)
	return &asset
}

// HoldingBuilder helps create test holdings
type HoldingBuilder struct {
	holding models.Holding
}

func NewHolding() *HoldingBuilder {
	return &HoldingBuilder{
		holding: models.Holding{
			Amount:       1.0,
			PurchaseDate: time.Now(),
		},
	}
}

func (b *HoldingBuilder) WithAccount(account *models.Account) *HoldingBuilder {
	b.holding.AccountID = account.ID
	b.holding.Account = *account
	return b
}

func (b *HoldingBuilder) WithAsset(asset *models.Asset) *HoldingBuilder {
	b.holding.AssetID = asset.ID
	b.holding.Asset = *asset
	return b
}

func (b *HoldingBuilder) WithAmount(amount float64) *HoldingBuilder {
	b.holding.Amount = amount
	return b
}

func (b *HoldingBuilder) WithPurchasePrice(price float64) *HoldingBuilder {
	b.holding.PurchasePrice = price
	return b
}

func (b *HoldingBuilder) WithPurchaseDate(date time.Time) *HoldingBuilder {
	b.holding.PurchaseDate = date
	return b
}

func (b *HoldingBuilder) Build() models.Holding {
	return b.holding
}

func (b *HoldingBuilder) Create(t *testing.T, db *gorm.DB) *models.Holding {
	holding := b.Build()
	require.NoError(t, db.Create(&holding).Error)
	return &holding
}

// Common test data sets
func CreateSamplePortfolio(t *testing.T, db *gorm.DB) {
	// Create accounts
	hardwareWallet := NewAccount().WithName("hardware wallet").WithType("wallet").Create(t, db)
	neobank := NewAccount().WithName("NeoBank").WithType("bank").Create(t, db)

	// Create assets
	btc := NewAsset().WithSymbol("BTC").WithName("Bitcoin").Create(t, db)
	eth := NewAsset().WithSymbol("ETH").WithName("Ethereum").Create(t, db)
	usd := NewAsset().WithSymbol("USD").WithName("US Dollar").WithType(models.AssetTypeFiat).Create(t, db)

	// Create holdings
	NewHolding().
		WithAccount(hardwareWallet).
		WithAsset(btc).
		WithAmount(0.5).
		WithPurchasePrice(40000).
		Create(t, db)

	NewHolding().
		WithAccount(hardwareWallet).
		WithAsset(eth).
		WithAmount(10).
		WithPurchasePrice(2000).
		Create(t, db)

	NewHolding().
		WithAccount(neobank).
		WithAsset(usd).
		WithAmount(1000).
		WithPurchasePrice(1).
		Create(t, db)
}
