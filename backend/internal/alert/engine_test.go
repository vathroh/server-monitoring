package alert_test

import (
	"testing"

	"github.com/velocity/server-monitoring/backend/internal/alert"
	"github.com/velocity/server-monitoring/backend/internal/domain"
)

type MockAlertRepository struct {
	openAlert *domain.Alert
	created   bool
	updated   bool
}

func (m *MockAlertRepository) FindOpenByServerAndRule(serverID uint, rule string) (*domain.Alert, error) {
	if m.openAlert != nil && m.openAlert.ServerID == serverID && m.openAlert.Rule == rule {
		return m.openAlert, nil
	}
	return nil, nil // nil simulates not found
}
func (m *MockAlertRepository) Create(alert *domain.Alert) error {
	m.created = true
	return nil
}
func (m *MockAlertRepository) Update(alert *domain.Alert) error {
	m.updated = true
	return nil
}
func (m *MockAlertRepository) FindAll(state string, serverID uint) ([]domain.Alert, error) { return nil, nil }
func (m *MockAlertRepository) CountOpen() (int64, error) { return 0, nil }

type MockNotifier struct {
	queued bool
}

func (m *MockNotifier) QueueAlert(alert *domain.Alert) {
	m.queued = true
}

func TestProcessMetrics_Critical(t *testing.T) {
	repo := &MockAlertRepository{}
	notifier := &MockNotifier{}
	engine := alert.NewEngine(repo, notifier)

	metric := &domain.Metric{
		ServerID:    1,
		CPUUsage:    99.0, // Critical
		MemoryUsage: 50.0, // Normal
		DiskUsage:   50.0, // Normal
	}

	engine.ProcessMetrics(metric)

	if !repo.created {
		t.Errorf("Expected alert to be created for critical CPU")
	}
	if !notifier.queued {
		t.Errorf("Expected alert to be queued to notifier")
	}
}

func TestProcessMetrics_Recovered(t *testing.T) {
	repo := &MockAlertRepository{
		openAlert: &domain.Alert{
			ServerID: 1,
			Rule:     "CPU",
			State:    "OPEN",
		},
	}
	notifier := &MockNotifier{}
	engine := alert.NewEngine(repo, notifier)

	metric := &domain.Metric{
		ServerID:    1,
		CPUUsage:    20.0, // Normal, should resolve the open alert
	}

	engine.ProcessMetrics(metric)

	if !repo.updated {
		t.Errorf("Expected open alert to be updated (resolved)")
	}
	if repo.openAlert.State != "RESOLVED" {
		t.Errorf("Expected alert state to be RESOLVED")
	}
	if !notifier.queued {
		t.Errorf("Expected resolved alert to be queued to notifier")
	}
}
