package domain

import "time"

type Metric struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	ServerID    uint      `json:"server_id" gorm:"index;not null"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
	Uptime      uint64    `json:"uptime"`
	CreatedAt   time.Time `json:"created_at"`
}
