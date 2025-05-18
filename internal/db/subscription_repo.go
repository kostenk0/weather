package db

import (
	"context"
	"database/sql"
	"weather/internal/models"
)

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *models.Subscription) error {
	query := `
        INSERT INTO subscriptions (email, city, frequency, confirmed, token)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
	return r.DB.QueryRowContext(ctx, query, s.Email, s.City, s.Frequency, s.Confirmed, s.Token).
		Scan(&s.ID, &s.CreatedAt)
}

func (r *SubscriptionRepository) GetByToken(ctx context.Context, token string) (*models.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE token = $1`
	row := r.DB.QueryRowContext(ctx, query, token)

	var s models.Subscription
	err := row.Scan(&s.ID, &s.Email, &s.City, &s.Frequency, &s.Confirmed, &s.Token, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) ConfirmByToken(ctx context.Context, token string) error {
	query := `UPDATE subscriptions SET confirmed = true WHERE token = $1 AND confirmed = false`
	result, err := r.DB.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows // або кастомна помилка
	}

	return nil
}

func (r *SubscriptionRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM subscriptions WHERE token = $1`
	result, err := r.DB.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
