package repository

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	CreateToken(token *domain.PasswordResetToken) error
	FindToken(tokenStr string) (*domain.PasswordResetToken, error)
	DeleteToken(id uint) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) CreateToken(token *domain.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *passwordResetRepository) FindToken(tokenStr string) (*domain.PasswordResetToken, error) {
	var token domain.PasswordResetToken
	err := r.db.Where("token = ?", tokenStr).First(&token).Error
	return &token, err
}

func (r *passwordResetRepository) DeleteToken(id uint) error {
	return r.db.Delete(&domain.PasswordResetToken{}, id).Error
}
