# Builder stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application (IMPORTANT!)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .  # -o main is crucial

# Install necessary packages for migrate
RUN apk add --no-cache wget tar gzip

# Set migrate version
ARG MIGRATE_VERSION=v4.18.2  # Or your desired version

# Construct the URL dynamically
RUN wget "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" -O migrate_linux_amd64.tar.gz

# Verify the download (important!)
RUN ls -l migrate_linux_amd64.tar.gz

# Extract and install
RUN gzip -d migrate_linux_amd64.tar.gz && tar -xf migrate_linux_amd64.tar && mv migrate /usr/local/bin/migrate

# Copy .env file (CRUCIAL):
# COPY ../.env /app/.env 

# Source the .env file to load environment variables
# RUN echo "source /app/.env" >> /root/.bashrc

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/.env /app/.env

# Copy necessary SSL certificates if your app requires them
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["./main"]