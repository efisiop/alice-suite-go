#!/bin/sh

# Always run migrations to ensure database schema is up to date
echo "Running database migrations..."
export DB_PATH="${DB_PATH:-data/alice-suite.db}"
mkdir -p "$(dirname "$DB_PATH")"
./bin/migrate

# Always run init-users to ensure all users exist (it checks and only creates if missing)
echo "Ensuring users are initialized..."
./bin/init-users

# Run fix-render to ensure sections data is correct (idempotent - safe to run multiple times)
echo "Checking and fixing sections data if needed..."
./bin/fix-render

# Start the server
echo "Starting server..."
exec ./bin/server
