#!/bin/sh

echo "Running migrations..."

goose -dir ./migrations postgres "$DATABASE_URL" up

echo "Starting the server..."
./weather