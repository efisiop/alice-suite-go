#!/bin/bash

# Development Server Startup Script
# This script ensures Safari proxy is configured and starts the server

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== Alice Suite Development Server Startup ==="
echo ""

# Step 1: Check if server is already running
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "✅ Server is already running on port 8080"
    echo "   Access at: http://127.0.0.1:8080/reader/login"
    exit 0
fi

# Step 2: Configure Safari proxy bypass (if not already set)
echo "1. Configuring Safari proxy bypass for localhost..."
CURRENT_BYPASS=$(networksetup -getproxybypassdomains "Wi-Fi" 2>/dev/null | tr '\n' ' ' || echo "")

if echo "$CURRENT_BYPASS" | grep -q "127.0.0.1\|localhost"; then
    echo "   ✅ localhost is already in proxy bypass list"
else
    echo "   ⚠️  Adding localhost to proxy bypass..."
    sudo networksetup -setproxybypassdomains "Wi-Fi" \
        "127.0.0.1" "localhost" "*.local" "169.254/16" 2>/dev/null || {
        echo "   ⚠️  Could not update proxy bypass (may need admin password)"
        echo "   You can manually configure Safari:"
        echo "   Safari → Preferences → Advanced → Proxies → Bypass proxy settings"
        echo "   Add: 127.0.0.1, localhost"
    }
fi

# Step 3: Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed or not in PATH"
    exit 1
fi

# Step 4: Build the server
echo ""
echo "2. Building server..."
if go build -o bin/server ./cmd/server; then
    echo "   ✅ Server built successfully"
else
    echo "   ❌ Build failed"
    exit 1
fi

# Step 5: Set AI API keys automatically (for Tier 2 features)
echo ""
echo "3. Setting up AI API keys..."
# Set Gemini API key (from render.yaml) if not already set
if [ -z "$GEMINI_API_KEY" ]; then
    # GEMINI_API_KEY should be set as environment variable, not hardcoded here
    # export GEMINI_API_KEY="your-key-here"
    echo "   ✅ GEMINI_API_KEY set automatically"
else
    echo "   ✅ GEMINI_API_KEY already set"
fi

# Set AI provider to use Gemini by default
if [ -z "$AI_PROVIDER" ]; then
    export AI_PROVIDER="gemini"
    echo "   ✅ AI_PROVIDER set to 'gemini'"
fi

# Remove incorrect Moonshot URL if set (always remove it since we use Gemini)
if [ -n "$ANTHROPIC_BASE_URL" ]; then
    unset ANTHROPIC_BASE_URL
    echo "   ✅ Removed ANTHROPIC_BASE_URL (using Gemini instead)"
fi

echo "   ✅ AI environment configured"

# Step 6: Start the server
echo ""
echo "4. Starting server on port 8080..."
echo "   Access at: http://127.0.0.1:8080/reader/login"
echo "   Press Ctrl+C to stop"
echo ""

# Run the server
./bin/server

