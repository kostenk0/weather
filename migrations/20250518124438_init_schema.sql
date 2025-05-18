-- +goose Up

CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               email TEXT NOT NULL,
                               city TEXT NOT NULL,
                               frequency TEXT NOT NULL CHECK (frequency IN ('hourly', 'daily')),
                               confirmed BOOLEAN NOT NULL DEFAULT false,
                               token TEXT NOT NULL UNIQUE,
                               created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_unique_email_city ON subscriptions(email, city);

CREATE TABLE weather_cache (
                               city TEXT PRIMARY KEY,
                               temperature NUMERIC,
                               humidity NUMERIC,
                               description TEXT,
                               updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS weather_cache;
DROP TABLE IF EXISTS subscriptions;