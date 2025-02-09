#!/bin/bash

# chmod +x ./vendor/github.com/golang-migrate/migrate/v4/cmd/migrate
# ./vendor/github.com/golang-migrate/migrate/v4/migrate.go -database "$DATABASE_URL" -path migrations up
migrate -database "$DATABASE_URL" -path migrations up

echo "Migrations complete."