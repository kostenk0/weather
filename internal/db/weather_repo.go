package db

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"weather/internal/models"
)

type WeatherRepository struct {
	DB *sql.DB
}

func NewWeatherRepository(db *sql.DB) *WeatherRepository {
	return &WeatherRepository{DB: db}
}

func (r *WeatherRepository) GetCached(ctx context.Context, city string) (*models.Weather, error) {
	city = strings.TrimSpace(city)

	query := `SELECT *
	          FROM weather_cache
	          WHERE LOWER(city) = LOWER($1)`

	row := r.DB.QueryRowContext(ctx, query, city)

	var w models.Weather
	err := row.Scan(&w.City, &w.Temperature, &w.Humidity, &w.Description, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		log.Println("[DEBUG] No weather found in cache")
		return nil, nil
	} else if err != nil {
		log.Printf("[DEBUG] Scan error: %v", err)
		return nil, err
	}

	return &w, nil
}

func (r *WeatherRepository) Save(ctx context.Context, w *models.Weather) error {
	query := `
        INSERT INTO weather_cache (city, temperature, humidity, description, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (city) DO UPDATE SET 
            temperature = EXCLUDED.temperature,
            humidity = EXCLUDED.humidity,
            description = EXCLUDED.description,
            updated_at = EXCLUDED.updated_at
    `
	_, err := r.DB.ExecContext(ctx, query,
		w.City, w.Temperature, w.Humidity, w.Description, w.UpdatedAt,
	)
	return err
}
