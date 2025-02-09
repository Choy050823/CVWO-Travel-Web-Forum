#!/bin/bash

chmod +x ./vendor/github.com/golang-migrate/migrate/v4/cmd/migrate
./vendor/github.com/golang-migrate/migrate/v4/cmd/migrate -database "$DATABASE_URL" -path migrations up

echo "Migrations complete."