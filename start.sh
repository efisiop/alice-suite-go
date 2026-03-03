#!/bin/sh

# Load .env file if it exists (for local development)
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    set -a
    . ./.env
    set +a
fi

# Production (e.g. Render) requires PostgreSQL - SQLite data is ephemeral and lost on deploy
if [ "${ENV:-}" = "production" ] && [ -z "$DATABASE_URL" ]; then
    echo "ERROR: DATABASE_URL is required in production. Add a Render PostgreSQL database and link it."
    echo "See docs/POSTGRES_MIGRATION.md for setup."
    exit 1
fi

# Run migrations before server starts
echo "Running database migrations..."
if [ -n "$DATABASE_URL" ]; then
  echo "Using PostgreSQL (persistent storage)"
else
  export DB_PATH="${DB_PATH:-data/alice-suite.db}"
  mkdir -p "$(dirname "$DB_PATH")"
  echo "Using SQLite at $DB_PATH"
fi
./bin/migrate

# Always run init-users to ensure all users exist (it checks and only creates if missing)
echo "Ensuring users are initialized..."
./bin/init-users

# Run fix-render to ensure sections and data are correct (especially important for Render.com)
# This is safe to run multiple times - it checks and only fixes if needed
if [ -f "./bin/fix-render" ]; then
    echo "Verifying and fixing sections data..."
    if ./bin/fix-render; then
        echo "✅ Sections fix completed successfully"
    else
        echo "⚠️  Warning: fix-render exited with error, but continuing..."
        # Don't fail startup - sections might still work
    fi
else
    echo "⚠️  Warning: fix-render binary not found, skipping sections fix"
    echo "   This might cause sections to not display correctly on Render"
fi

# Optional: Run deployment verification (can be disabled for faster startup)
# Uncomment the next 3 lines to enable verification on every start
# if [ -f "./bin/verify-deployment" ]; then
#     echo "Running deployment verification..."
#     ./bin/verify-deployment || echo "⚠️  Verification found issues (non-fatal)"
# fi

# Start the server
echo "Starting server..."
exec ./bin/server
