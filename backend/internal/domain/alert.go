package domain

import "time"

type Alert struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	ServerID   uint       `json:"server_id" gorm:"index;not null"`
	Rule       string     `json:"rule" gorm:"index;not null"` // CPU, MEMORY, DISK, OFFLINE
	Severity   string     `json:"severity" gorm:"not null"`   // WARNING, CRITICAL
	State      string     `json:"state" gorm:"index;not null"`// OPEN, RESOLVED
	Message    string     `json:"message" gorm:"not null"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at"`

	Server *Server `json:"server,omitempty" gorm:"foreignKey:ServerID"`
}
