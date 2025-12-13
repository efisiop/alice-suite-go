#!/bin/sh

# Initialize database if it doesn't exist
if [ ! -f "$DB_PATH" ]; then
    echo "Initializing database..."
    ./bin/migrate
    ./bin/init-users
    echo "Database initialized successfully"
fi

# Start the server
echo "Starting server..."
exec ./bin/server
