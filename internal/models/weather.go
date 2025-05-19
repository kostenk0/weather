package models

import "time"

type Weather struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
