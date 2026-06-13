package service

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
)

type AlertService interface {
	GetAlerts(state string, serverID uint) ([]domain.Alert, error)
}

type alertService struct {
	repo repository.AlertRepository
}

func NewAlertService(repo repository.AlertRepository) AlertService {
	return &alertService{repo: repo}
}

func (s *alertService) GetAlerts(state string, serverID uint) ([]domain.Alert, error) {
	return s.repo.FindAll(state, serverID)
}
