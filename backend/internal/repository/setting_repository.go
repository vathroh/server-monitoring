package repository

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingRepository interface {
	GetAll() ([]domain.Setting, error)
	GetByKey(key string) (*domain.Setting, error)
	Upsert(setting *domain.Setting) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) GetAll() ([]domain.Setting, error) {
	var settings []domain.Setting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *settingRepository) GetByKey(key string) (*domain.Setting, error) {
	var settings []domain.Setting
	result := r.db.Where("key = ?", key).Limit(1).Find(&settings)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(settings) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &settings[0], nil
}

func (r *settingRepository) Upsert(setting *domain.Setting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(setting).Error
}
