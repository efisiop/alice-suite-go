# Build stage - Use Debian-based image for better SQLite compatibility
FROM golang:1.23 AS builder

# Enable toolchain support
ENV GOTOOLCHAIN=auto

WORKDIR /app

# Install build dependencies (Debian/Ubuntu packages)
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/server ./cmd/server
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/migrate ./cmd/migrate
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/init-users ./cmd/init-users

# Production stage - Use Debian slim for smaller size
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

# Copy binaries from builder
COPY --from=builder /app/bin/server ./bin/server
COPY --from=builder /app/bin/migrate ./bin/migrate
COPY --from=builder /app/bin/init-users ./bin/init-users

# Copy required files
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/static ./internal/static
COPY --from=builder /app/internal/templates ./internal/templates

# Create data directory
RUN mkdir -p /data

# Set environment variables
ENV PORT=8080
ENV DB_PATH=/data/alice-suite.db
ENV ENV=production

# Expose port
EXPOSE 8080

# Start script
COPY --from=builder /app/start.sh ./start.sh
RUN chmod +x ./start.sh

CMD ["./start.sh"]
