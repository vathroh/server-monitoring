package repository

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"gorm.io/gorm"
)

type ServerRepository interface {
	Create(server *domain.Server) error
	FindAll(offset, limit int) ([]domain.Server, int64, error)
	FindByID(id uint) (*domain.Server, error)
	FindByAPIKey(apiKey string) (*domain.Server, error)
	GetDashboardSummary() (total, online, warning, offline int64, err error)
	Update(server *domain.Server) error
	UpdateStatusBatch(ids []uint, status string) error
	Delete(id uint) error
}

type serverRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{db: db}
}

func (r *serverRepository) Create(server *domain.Server) error {
	return r.db.Create(server).Error
}

func (r *serverRepository) FindAll(offset, limit int) ([]domain.Server, int64, error) {
	var servers []domain.Server
	var total int64

	err := r.db.Model(&domain.Server{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&servers).Error
	return servers, total, err
}

func (r *serverRepository) FindByID(id uint) (*domain.Server, error) {
	var server domain.Server
	err := r.db.First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *serverRepository) FindByAPIKey(apiKey string) (*domain.Server, error) {
	var server domain.Server
	err := r.db.Where("api_key = ?", apiKey).First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *serverRepository) GetDashboardSummary() (int64, int64, int64, int64, error) {
	var total, online, warning, offline int64

	err := r.db.Model(&domain.Server{}).Count(&total).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	r.db.Model(&domain.Server{}).Where("status = ?", "ONLINE").Count(&online)
	r.db.Model(&domain.Server{}).Where("status = ?", "WARNING").Count(&warning)
	r.db.Model(&domain.Server{}).Where("status = ?", "OFFLINE").Count(&offline)

	return total, online, warning, offline, nil
}

func (r *serverRepository) Update(server *domain.Server) error {
	return r.db.Save(server).Error
}

func (r *serverRepository) UpdateStatusBatch(ids []uint, status string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&domain.Server{}).Where("id IN ?", ids).Update("status", status).Error
}
func (r *serverRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Server{}, id).Error
}
