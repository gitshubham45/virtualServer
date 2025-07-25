package models

import (
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Billing_rate float64 `json:"billing_rate"`
	Status       string  `json:"status"`
	Region       string  `json:"region"`
	Type         string  `json:"type"`
}
