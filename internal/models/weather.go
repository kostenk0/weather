package models

import "time"

type Weather struct {
	City        string    `db:"city"`
	Temperature float64   `db:"temperature"`
	Humidity    float64   `db:"humidity"`
	Description string    `db:"description"`
	UpdatedAt   time.Time `db:"updated_at"`
}
