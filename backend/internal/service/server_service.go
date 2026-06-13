package service

import (
	"github.com/google/uuid"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
)

type PaginatedServers struct {
	Data  []domain.Server `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

	type DashboardSummary struct {
		TotalServers   int64 `json:"total_servers"`
		OnlineServers  int64 `json:"online_servers"`
		WarningServers int64 `json:"warning_servers"`
		OfflineServers int64 `json:"offline_servers"`
		ActiveAlerts   int64 `json:"active_alerts"`
	}
	
	type ServerService interface {
		CreateServer(server *domain.Server) error
		GetServers(page, limit int) (*PaginatedServers, error)
		GetServerByID(id uint) (*domain.Server, error)
		GetServerByAPIKey(apiKey string) (*domain.Server, error)
		GetDashboardSummary() (*DashboardSummary, error)
		UpdateServer(server *domain.Server) error
		DeleteServer(id uint) error
	}

type serverService struct {
	repo      repository.ServerRepository
	alertRepo repository.AlertRepository
}

func NewServerService(repo repository.ServerRepository, alertRepo repository.AlertRepository) ServerService {
	return &serverService{repo: repo, alertRepo: alertRepo}
}

func (s *serverService) CreateServer(server *domain.Server) error {
	server.APIKey = uuid.New().String()
	return s.repo.Create(server)
}

func (s *serverService) GetServers(page, limit int) (*PaginatedServers, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	servers, total, err := s.repo.FindAll(offset, limit)
	if err != nil {
		return nil, err
	}

	return &PaginatedServers{
		Data:  servers,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *serverService) GetServerByID(id uint) (*domain.Server, error) {
	return s.repo.FindByID(id)
}

func (s *serverService) GetServerByAPIKey(apiKey string) (*domain.Server, error) {
	return s.repo.FindByAPIKey(apiKey)
}

func (s *serverService) GetDashboardSummary() (*DashboardSummary, error) {
	total, online, warning, offline, err := s.repo.GetDashboardSummary()
	if err != nil {
		return nil, err
	}
	activeAlerts, _ := s.alertRepo.CountOpen()

	return &DashboardSummary{
		TotalServers:   total,
		OnlineServers:  online,
		WarningServers: warning,
		OfflineServers: offline,
		ActiveAlerts:   activeAlerts,
	}, nil
}

func (s *serverService) UpdateServer(server *domain.Server) error {
	return s.repo.Update(server)
}

func (s *serverService) DeleteServer(id uint) error {
	return s.repo.Delete(id)
}
