#!/bin/bash

# Install the migrate tool if not already installed
if ! command -v migrate &> /dev/null; then
    echo "Installing migrate tool..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Run the migrations
migrate -database "$DATABASE_URL" -path migrations up

echo "Migrations complete."