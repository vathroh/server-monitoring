package repository

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(log *domain.AuditLog) error
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *domain.AuditLog) error {
	return r.db.Create(log).Error
}
