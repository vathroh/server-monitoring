package repository

import (
	"time"

	"github.com/velocity/server-monitoring/backend/internal/domain"
	"gorm.io/gorm"
)

type MetricRepository interface {
	Create(metric *domain.Metric) error
	FindByServerID(serverID uint, limit int) ([]domain.Metric, error)
	GetTrend(serverID uint, from time.Time) ([]domain.Metric, error)
	DeleteOlderThan(threshold time.Time) error
}

type metricRepository struct {
	db *gorm.DB
}

func NewMetricRepository(db *gorm.DB) MetricRepository {
	return &metricRepository{db: db}
}

func (r *metricRepository) Create(metric *domain.Metric) error {
	return r.db.Create(metric).Error
}

func (r *metricRepository) FindByServerID(serverID uint, limit int) ([]domain.Metric, error) {
	var metrics []domain.Metric
	err := r.db.Where("server_id = ?", serverID).Order("created_at desc").Limit(limit).Find(&metrics).Error
	return metrics, err
}

func (r *metricRepository) GetTrend(serverID uint, from time.Time) ([]domain.Metric, error) {
	var metrics []domain.Metric
	err := r.db.Where("server_id = ? AND created_at >= ?", serverID, from).Order("created_at asc").Find(&metrics).Error
	return metrics, err
}

func (r *metricRepository) DeleteOlderThan(threshold time.Time) error {
	return r.db.Where("created_at < ?", threshold).Delete(&domain.Metric{}).Error
}
