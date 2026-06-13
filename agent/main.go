package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type MetricPayload struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Uptime      uint64  `json:"uptime"`
}

func main() {
	apiURL := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")

	if apiURL == "" || apiKey == "" {
		log.Fatal("API_URL and API_KEY environment variables are required")
	}

	log.Println("Starting Velocity Server Monitoring Agent...")
	log.Printf("Target API: %s\n", apiURL)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		collectAndSend(client, apiURL, apiKey)
		<-ticker.C
	}
}

func collectAndSend(client *http.Client, apiURL, apiKey string) {
	// 1. Collect Metrics
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

	payload := MetricPayload{
		CPUUsage:    cpuUsage,
		MemoryUsage: memUsage,
		DiskUsage:   diskUsage,
		Uptime:      uptime,
	}

	// 2. Send payload with retries
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v\n", err)
		return
	}

	maxRetries := 3
	backoff := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error creating request: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Attempt %d failed: %v\n", i+1, err)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Println("Metrics sent successfully.")
			return
		}

		log.Printf("Attempt %d failed with status code %d\n", i+1, resp.StatusCode)
		time.Sleep(backoff)
		backoff *= 2
	}

	log.Println("Failed to send metrics after max retries")
}
