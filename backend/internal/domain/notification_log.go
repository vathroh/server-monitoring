package domain

import "time"

type NotificationLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Provider  string    `json:"provider" gorm:"index"` // TELEGRAM, WHATSAPP
	AlertID   uint      `json:"alert_id" gorm:"index"`
	Status    string    `json:"status"` // SUCCESS, FAILED
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
}
