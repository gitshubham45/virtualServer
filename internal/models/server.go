package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID           string         `gorm:"primaryKey;type:uuid" json:"id"`
	ServerNumber int64          `json:"serverNumber" gorm:"autoIncrement"`
	BillingRate  float64        `json:"billingRate"`
	Status       string         `json:"status"`
	Region       string         `json:"region"`
	Type         string         `json:"type"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
