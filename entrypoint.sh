#!/bin/sh

echo "Waiting for postgres..."

until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  sleep 1
done

echo "PostgreSQL is ready. Running migrations..."

goose -dir ./migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable" up

echo "Starting the server..."
./weather