#!/bin/sh

migrate -path ./db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

if [ "$ENV" = "production" ]; then
    echo "Starting backend service using build binary"
    /app/main
else
    echo "Starting backend service using hot reload"
    air
fi  