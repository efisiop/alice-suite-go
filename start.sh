#!/bin/sh

# Initialize database if it doesn't exist
if [ ! -f "$DB_PATH" ]; then
    echo "Initializing database..."
    ./bin/migrate
    echo "Database initialized successfully"
fi

# Always run init-users to ensure all users exist (it checks and only creates if missing)
echo "Ensuring users are initialized..."
./bin/init-users

# Start the server
echo "Starting server..."
exec ./bin/server
