package repository

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"gorm.io/gorm"
)

type NotificationLogRepository interface {
	Create(log *domain.NotificationLog) error
}

type notificationLogRepository struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) NotificationLogRepository {
	return &notificationLogRepository{db: db}
}

func (r *notificationLogRepository) Create(log *domain.NotificationLog) error {
	return r.db.Create(log).Error
}
