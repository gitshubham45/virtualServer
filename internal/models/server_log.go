package models

import "time"

type ServerLog struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	ServerID  string    `gorm:"index" json:"serverId"`
	EventType string    `json:"eventType"`
	Message   string    `json:"message"`
	OldStatus string    `json:"oldStatus"`
	NewStatus string    `json:"newStatus"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
