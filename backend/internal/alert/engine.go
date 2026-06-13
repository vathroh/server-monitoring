package alert

import (
	"fmt"
	"log"
	"time"

	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
)

type Engine interface {
	ProcessMetrics(metric *domain.Metric)
	ProcessServerStatus(serverID uint, status string)
}

type Notifier interface {
	QueueAlert(alert *domain.Alert)
}

type engine struct {
	repo     repository.AlertRepository
	notifier Notifier
}

func NewEngine(repo repository.AlertRepository, notifier Notifier) Engine {
	return &engine{repo: repo, notifier: notifier}
}

func (e *engine) ProcessMetrics(metric *domain.Metric) {
	e.evaluateRule(metric.ServerID, "CPU", metric.CPUUsage, 80.0, 95.0)
	e.evaluateRule(metric.ServerID, "MEMORY", metric.MemoryUsage, 80.0, 95.0)
	e.evaluateRule(metric.ServerID, "DISK", metric.DiskUsage, 85.0, 95.0)
}

func (e *engine) evaluateRule(serverID uint, rule string, value, warnThreshold, critThreshold float64) {
	var severity string
	var message string

	if value >= critThreshold {
		severity = "CRITICAL"
		message = fmt.Sprintf("%s usage is critical: %.2f%%", rule, value)
	} else if value >= warnThreshold {
		severity = "WARNING"
		message = fmt.Sprintf("%s usage is high: %.2f%%", rule, value)
	}

	e.handleAlertState(serverID, rule, severity, message, value < warnThreshold)
}

func (e *engine) ProcessServerStatus(serverID uint, status string) {
	isRecovered := status == "ONLINE" || status == "WARNING"
	var severity string
	var message string

	if status == "OFFLINE" {
		severity = "CRITICAL"
		message = "Server is OFFLINE (No heartbeat for 5+ minutes)"
	}

	e.handleAlertState(serverID, "OFFLINE", severity, message, isRecovered)
}

func (e *engine) handleAlertState(serverID uint, rule, severity, message string, isRecovered bool) {
	existing, err := e.repo.FindOpenByServerAndRule(serverID, rule)
	if err != nil && err.Error() != "record not found" {
		return
	}

	if isRecovered {
		if existing != nil {
			now := time.Now()
			existing.State = "RESOLVED"
			existing.ResolvedAt = &now
			e.repo.Update(existing)
			log.Printf("Alert Resolved: Server %d [%s]", serverID, rule)
			e.notifier.QueueAlert(existing)
		}
		return
	}

	if existing == nil {
		newAlert := &domain.Alert{
			ServerID: serverID,
			Rule:     rule,
			Severity: severity,
			State:    "OPEN",
			Message:  message,
		}
		e.repo.Create(newAlert)
		log.Printf("Alert Triggered: Server %d [%s] - %s", serverID, rule, severity)
		e.notifier.QueueAlert(newAlert)
	} else if existing.Severity != severity {
		existing.Severity = severity
		existing.Message = message
		e.repo.Update(existing)
		log.Printf("Alert Updated: Server %d [%s] - %s", serverID, rule, severity)
		e.notifier.QueueAlert(existing)
	}
}
