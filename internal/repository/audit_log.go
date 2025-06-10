package repository

import (
	"time"

	"github.com/bioharz/budget/internal/models"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) GetByEntity(entityType models.AuditLogEntityType, entityID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at desc").
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) GetAll(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) GetByDateRange(start, end time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at desc").
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) GetByAction(action models.AuditLogAction, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Where("action = ?", action).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&logs).Error
	return logs, err
}