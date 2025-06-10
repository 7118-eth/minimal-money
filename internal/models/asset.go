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

type Account struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	Type      string         `gorm:"not null"`
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

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
	ID            uint           `gorm:"primaryKey"`
	AccountID     uint           `gorm:"not null"`
	Account       Account        `gorm:"foreignKey:AccountID"`
	AssetID       uint           `gorm:"not null"`
	Asset         Asset          `gorm:"foreignKey:AssetID"`
	Amount        float64        `gorm:"not null"`
	PurchasePrice float64
	PurchaseDate  time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type AuditLogAction string

const (
	AuditActionCreate AuditLogAction = "CREATE"
	AuditActionUpdate AuditLogAction = "UPDATE"
	AuditActionDelete AuditLogAction = "DELETE"
)

type AuditLogEntityType string

const (
	AuditEntityHolding AuditLogEntityType = "HOLDING"
	AuditEntityAsset   AuditLogEntityType = "ASSET"
	AuditEntityAccount AuditLogEntityType = "ACCOUNT"
)

type AuditLog struct {
	ID         uint               `gorm:"primaryKey"`
	Action     AuditLogAction     `gorm:"not null"`
	EntityType AuditLogEntityType `gorm:"not null"`
	EntityID   uint               `gorm:"not null"`
	OldValue   string             `gorm:"type:text"` // JSON representation
	NewValue   string             `gorm:"type:text"` // JSON representation
	UserNote   string
	CreatedAt  time.Time          `gorm:"not null;index"`
}

type PortfolioSnapshot struct {
	ID           uint                   `gorm:"primaryKey"`
	TotalValueUSD float64               `gorm:"not null"`
	Details      map[string]interface{} `gorm:"serializer:json"`
	Timestamp    time.Time              `gorm:"not null;index"`
}