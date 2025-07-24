package models

import "time"

type Server struct{
	ID string `json:"id"`
	Billing_rate float64 `json:"billing_rate"`
	Status string `json:"status"`
	Region string `json:"region"`
	Type string `json:"type"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}