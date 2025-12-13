#!/bin/bash

# Start Alice Suite Go Server
# This script starts the server and opens it in your browser

cd "$(dirname "$0")"

echo "🚀 Starting Alice Suite Go Server..."
echo ""

# Kill any existing server on port 8080
lsof -ti:8080 | xargs kill -9 2>/dev/null
sleep 1

# Start the server in background
echo "📡 Starting server on http://localhost:8080"
go run cmd/reader/main.go > /tmp/alice-server.log 2>&1 &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Check if server started successfully
if ps -p $SERVER_PID > /dev/null; then
    echo "✅ Server started successfully (PID: $SERVER_PID)"
    echo ""
    echo "📖 Open your browser to: http://localhost:8080"
    echo ""
    echo "📊 Server logs: tail -f /tmp/alice-server.log"
    echo "🛑 To stop: kill $SERVER_PID"
    echo ""
    
    # Try to open browser (macOS)
    if command -v open > /dev/null; then
        open http://localhost:8080
    fi
else
    echo "❌ Server failed to start. Check logs:"
    cat /tmp/alice-server.log
    exit 1
fi



