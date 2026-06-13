package agent

import (
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/service"
)

type SelfMonitor struct {
	metricSvc service.MetricService
	serverSvc service.ServerService
	done      chan bool
}

func NewSelfMonitor(metricSvc service.MetricService, serverSvc service.ServerService) *SelfMonitor {
	return &SelfMonitor{
		metricSvc: metricSvc,
		serverSvc: serverSvc,
		done:      make(chan bool),
	}
}

func (s *SelfMonitor) Start() {
	go s.run()
}

func (s *SelfMonitor) Stop() {
	s.done <- true
}

func (s *SelfMonitor) run() {
	// Register the local server
	serverID := s.ensureLocalServer()

	log.Printf("Self-Monitor started for Local Server (ID: %d)", serverID)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial metrics push
	s.collectAndSave(serverID)

	for {
		select {
		case <-ticker.C:
			s.collectAndSave(serverID)
		case <-s.done:
			log.Println("Self-Monitor stopped")
			return
		}
	}
}

func (s *SelfMonitor) ensureLocalServer() uint {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Local Server"
	}

	servers, err := s.serverSvc.GetServers(1, 1000)
	if err == nil {
		for _, srv := range servers.Data {
			if srv.Hostname == hostname || srv.Name == "Local Server" {
				return srv.ID
			}
		}
	}

	// Create if not exists
	newServer := &domain.Server{
		Name:        "Local Server",
		Hostname:    hostname,
		IPAddress:   "127.0.0.1",
		Environment: "PRODUCTION",
		Description: "Auto-registered local server",
	}

	err = s.serverSvc.CreateServer(newServer)
	if err != nil {
		log.Printf("Failed to auto-register local server: %v", err)
		// Fallback ID if creation fails, though metrics might drop
		return 0
	}

	return newServer.ID
}

func (s *SelfMonitor) collectAndSave(serverID uint) {
	if serverID == 0 {
		return
	}

	cpuPercents, err := cpu.Percent(0, false)
	cpuUsage := 0.0
	if err == nil && len(cpuPercents) > 0 {
		cpuUsage = cpuPercents[0]
	}

	memInfo, err := mem.VirtualMemory()
	memUsage := 0.0
	if err == nil {
		memUsage = memInfo.UsedPercent
	}

	diskInfo, err := disk.Usage("/")
	diskUsage := 0.0
	if err == nil {
		diskUsage = diskInfo.UsedPercent
	}

	hostInfo, err := host.Info()
	uptime := uint64(0)
	if err == nil {
		uptime = hostInfo.Uptime
	}

	metric := &domain.Metric{
		ServerID:    serverID,
		CPUUsage:    cpuUsage,
		MemoryUsage: memUsage,
		DiskUsage:   diskUsage,
		Uptime:      uptime,
	}

	err = s.metricSvc.SaveMetric(metric)
	if err != nil {
		log.Printf("Self-Monitor failed to save metrics: %v", err)
	}
}
