package repository

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"gorm.io/gorm"
)

type AlertRepository interface {
	Create(alert *domain.Alert) error
	Update(alert *domain.Alert) error
	FindOpenByServerAndRule(serverID uint, rule string) (*domain.Alert, error)
	FindAll(state string, serverID uint) ([]domain.Alert, error)
	CountOpen() (int64, error)
}

type alertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) Create(alert *domain.Alert) error {
	return r.db.Create(alert).Error
}

func (r *alertRepository) Update(alert *domain.Alert) error {
	return r.db.Save(alert).Error
}

func (r *alertRepository) FindOpenByServerAndRule(serverID uint, rule string) (*domain.Alert, error) {
	var alerts []domain.Alert
	result := r.db.Where("server_id = ? AND rule = ? AND state = ?", serverID, rule, "OPEN").Limit(1).Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(alerts) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &alerts[0], nil
}

func (r *alertRepository) FindAll(state string, serverID uint) ([]domain.Alert, error) {
	var alerts []domain.Alert
	query := r.db.Preload("Server").Order("created_at desc")
	if state != "" {
		query = query.Where("state = ?", state)
	}
	if serverID > 0 {
		query = query.Where("server_id = ?", serverID)
	}
	err := query.Find(&alerts).Error
	return alerts, err
}

func (r *alertRepository) CountOpen() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Alert{}).Where("state = ?", "OPEN").Count(&count).Error
	return count, err
}
