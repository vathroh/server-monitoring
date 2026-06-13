package domain

import (
	"time"
	"gorm.io/gorm"
)

type Server struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null"`
	Hostname    string         `json:"hostname" gorm:"unique;not null"`
	IPAddress   string         `json:"ip_address" gorm:"not null"`
	Environment string         `json:"environment"`
	Status      string         `json:"status" gorm:"default:'OFFLINE'"`
	APIKey      string         `json:"api_key" gorm:"unique;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
