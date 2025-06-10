package models

import (
	"time"
	"gorm.io/gorm"
)

type AssetType string

const (
	AssetTypeCrypto AssetType = "crypto"
	AssetTypeFiat   AssetType = "fiat"
	AssetTypeStock  AssetType = "stock"
	AssetTypeOther  AssetType = "other"
)

type Asset struct {
	ID        uint           `gorm:"primaryKey"`
	Symbol    string         `gorm:"uniqueIndex;not null"`
	Name      string         `gorm:"not null"`
	Type      AssetType      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Holding struct {
	ID        uint           `gorm:"primaryKey"`
	AssetID   uint           `gorm:"not null"`
	Asset     Asset          `gorm:"foreignKey:AssetID"`
	Amount    float64        `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type PriceHistory struct {
	ID        uint      `gorm:"primaryKey"`
	AssetID   uint      `gorm:"not null;index"`
	Asset     Asset     `gorm:"foreignKey:AssetID"`
	PriceUSD  float64   `gorm:"not null"`
	Timestamp time.Time `gorm:"not null;index"`
}

type PortfolioSnapshot struct {
	ID           uint                   `gorm:"primaryKey"`
	TotalValueUSD float64               `gorm:"not null"`
	Details      map[string]interface{} `gorm:"serializer:json"`
	Timestamp    time.Time              `gorm:"not null;index"`
}