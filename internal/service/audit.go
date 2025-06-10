package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bioharz/budget/internal/db"
	"github.com/bioharz/budget/internal/models"
	"github.com/bioharz/budget/internal/repository"
	"gorm.io/gorm"
)

type AuditService struct {
	auditRepo *repository.AuditLogRepository
}

func NewAuditService() *AuditService {
	return &AuditService{
		auditRepo: repository.NewAuditLogRepository(db.DB),
	}
}

func NewAuditServiceWithDB(database *gorm.DB) *AuditService {
	return &AuditService{
		auditRepo: repository.NewAuditLogRepository(database),
	}
}

func (s *AuditService) LogHoldingCreate(holding *models.Holding) error {
	newValue, err := json.Marshal(map[string]interface{}{
		"account_id":     holding.AccountID,
		"asset_id":       holding.AssetID,
		"amount":         holding.Amount,
		"purchase_price": holding.PurchasePrice,
		"purchase_date":  holding.PurchaseDate,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal holding: %w", err)
	}

	log := &models.AuditLog{
		Action:     models.AuditActionCreate,
		EntityType: models.AuditEntityHolding,
		EntityID:   holding.ID,
		OldValue:   "",
		NewValue:   string(newValue),
		CreatedAt:  time.Now(),
	}

	return s.auditRepo.Create(log)
}

func (s *AuditService) LogHoldingUpdate(oldHolding, newHolding *models.Holding) error {
	oldValue, err := json.Marshal(map[string]interface{}{
		"account_id":     oldHolding.AccountID,
		"asset_id":       oldHolding.AssetID,
		"amount":         oldHolding.Amount,
		"purchase_price": oldHolding.PurchasePrice,
		"purchase_date":  oldHolding.PurchaseDate,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal old holding: %w", err)
	}

	newValue, err := json.Marshal(map[string]interface{}{
		"account_id":     newHolding.AccountID,
		"asset_id":       newHolding.AssetID,
		"amount":         newHolding.Amount,
		"purchase_price": newHolding.PurchasePrice,
		"purchase_date":  newHolding.PurchaseDate,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal new holding: %w", err)
	}

	log := &models.AuditLog{
		Action:     models.AuditActionUpdate,
		EntityType: models.AuditEntityHolding,
		EntityID:   newHolding.ID,
		OldValue:   string(oldValue),
		NewValue:   string(newValue),
		CreatedAt:  time.Now(),
	}

	return s.auditRepo.Create(log)
}

func (s *AuditService) LogHoldingDelete(holding *models.Holding) error {
	oldValue, err := json.Marshal(map[string]interface{}{
		"account_id":     holding.AccountID,
		"asset_id":       holding.AssetID,
		"amount":         holding.Amount,
		"purchase_price": holding.PurchasePrice,
		"purchase_date":  holding.PurchaseDate,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal holding: %w", err)
	}

	log := &models.AuditLog{
		Action:     models.AuditActionDelete,
		EntityType: models.AuditEntityHolding,
		EntityID:   holding.ID,
		OldValue:   string(oldValue),
		NewValue:   "",
		CreatedAt:  time.Now(),
	}

	return s.auditRepo.Create(log)
}

func (s *AuditService) GetAllLogs(limit int) ([]models.AuditLog, error) {
	return s.auditRepo.GetAll(limit)
}

func (s *AuditService) GetHoldingLogs(holdingID uint) ([]models.AuditLog, error) {
	return s.auditRepo.GetByEntity(models.AuditEntityHolding, holdingID)
}