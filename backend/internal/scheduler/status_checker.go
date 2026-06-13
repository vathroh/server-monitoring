package scheduler

import (
	"log"
	"time"

	"github.com/velocity/server-monitoring/backend/internal/alert"
	"github.com/velocity/server-monitoring/backend/internal/repository"
)

type StatusChecker interface {
	Start()
	Stop()
}

type statusChecker struct {
	serverRepo  repository.ServerRepository
	metricRepo  repository.MetricRepository
	alertEngine alert.Engine
	ticker      *time.Ticker
	quit        chan struct{}
}

func NewStatusChecker(serverRepo repository.ServerRepository, metricRepo repository.MetricRepository, engine alert.Engine) StatusChecker {
	return &statusChecker{
		serverRepo:  serverRepo,
		metricRepo:  metricRepo,
		alertEngine: engine,
		quit:        make(chan struct{}),
	}
}

func (s *statusChecker) Start() {
	s.ticker = time.NewTicker(1 * time.Minute)
	cleanupTicker := time.NewTicker(1 * time.Hour)
	go func() {
		log.Println("Starting Status Checker Scheduler...")
		// Run once immediately
		s.checkStatuses()
		s.cleanupOldMetrics()
		for {
			select {
			case <-s.ticker.C:
				s.checkStatuses()
			case <-cleanupTicker.C:
				s.cleanupOldMetrics()
			case <-s.quit:
				s.ticker.Stop()
				cleanupTicker.Stop()
				return
			}
		}
	}()
}

func (s *statusChecker) Stop() {
	close(s.quit)
}

func (s *statusChecker) checkStatuses() {
	// 1. Fetch all servers
	// In a real huge app, we'd paginate this. For 1-2 servers it's fine.
	servers, _, err := s.serverRepo.FindAll(0, 1000)
	if err != nil {
		log.Println("StatusChecker Error fetching servers:", err)
		return
	}

	var onlineIDs []uint
	var warningIDs []uint
	var offlineIDs []uint

	for _, srv := range servers {
		// 2. Find latest metric
		metrics, err := s.metricRepo.FindByServerID(srv.ID, 1)
		if err != nil || len(metrics) == 0 {
			offlineIDs = append(offlineIDs, srv.ID)
			continue
		}

		latest := metrics[0].CreatedAt
		timeSince := time.Since(latest)

		if timeSince >= 5*time.Minute {
			offlineIDs = append(offlineIDs, srv.ID)
			s.alertEngine.ProcessServerStatus(srv.ID, "OFFLINE")
		} else if timeSince >= 2*time.Minute {
			warningIDs = append(warningIDs, srv.ID)
			s.alertEngine.ProcessServerStatus(srv.ID, "WARNING")
		} else {
			onlineIDs = append(onlineIDs, srv.ID)
			s.alertEngine.ProcessServerStatus(srv.ID, "ONLINE")
		}
	}

	// 3. Batch Update Statuses
	if err := s.serverRepo.UpdateStatusBatch(onlineIDs, "ONLINE"); err != nil {
		log.Println("StatusChecker Error updating ONLINE:", err)
	}
	if err := s.serverRepo.UpdateStatusBatch(warningIDs, "WARNING"); err != nil {
		log.Println("StatusChecker Error updating WARNING:", err)
	}
	if err := s.serverRepo.UpdateStatusBatch(offlineIDs, "OFFLINE"); err != nil {
		log.Println("StatusChecker Error updating OFFLINE:", err)
	}
}

func (s *statusChecker) cleanupOldMetrics() {
	threshold := time.Now().Add(-7 * 24 * time.Hour)
	if err := s.metricRepo.DeleteOlderThan(threshold); err != nil {
		log.Println("StatusChecker Error cleaning up old metrics:", err)
	} else {
		log.Println("Successfully cleaned up metrics older than 7 days")
	}
}
