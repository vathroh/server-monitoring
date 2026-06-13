package service

import (
	"time"

	"github.com/velocity/server-monitoring/backend/internal/alert"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
)

type MetricService interface {
	SaveMetric(metric *domain.Metric) error
	GetLatestMetrics(serverID uint, limit int) ([]domain.Metric, error)
	GetMetricsTrend(serverID uint, timeRange string) ([]domain.Metric, error)
}

type metricService struct {
	repo        repository.MetricRepository
	alertEngine alert.Engine
}

func NewMetricService(repo repository.MetricRepository, engine alert.Engine) MetricService {
	return &metricService{repo: repo, alertEngine: engine}
}

func (s *metricService) SaveMetric(metric *domain.Metric) error {
	err := s.repo.Create(metric)
	if err == nil {
		s.alertEngine.ProcessMetrics(metric)
	}
	return err
}

func (s *metricService) GetLatestMetrics(serverID uint, limit int) ([]domain.Metric, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.FindByServerID(serverID, limit)
}

func (s *metricService) GetMetricsTrend(serverID uint, timeRange string) ([]domain.Metric, error) {
	now := time.Now()
	var from time.Time

	switch timeRange {
	case "24h":
		from = now.Add(-24 * time.Hour)
	case "7d":
		from = now.Add(-7 * 24 * time.Hour)
	case "1h":
		fallthrough
	default:
		from = now.Add(-1 * time.Hour)
	}

	metrics, err := s.repo.GetTrend(serverID, from)
	if err != nil {
		return nil, err
	}

	// Downsample to max 500 points for frontend performance
	if len(metrics) > 500 {
		step := len(metrics) / 500
		var downsampled []domain.Metric
		for i := 0; i < len(metrics); i += step {
			downsampled = append(downsampled, metrics[i])
		}
		return downsampled, nil
	}

	return metrics, nil
}
