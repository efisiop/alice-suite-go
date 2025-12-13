#!/bin/bash

# Test Safari Connection Script
# This script helps test if the server is accessible and working

echo "=== Testing Server Connection ==="
echo ""

# Test 1: Check if server is running on port 8080
echo "1. Checking if server is running on port 8080..."
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "   ✅ Server is running on port 8080"
else
    echo "   ❌ Server is NOT running on port 8080"
    echo "   Please start the server with: go run ./cmd/server"
    exit 1
fi

echo ""

# Test 2: Test health endpoint (bypass proxy for localhost)
echo "2. Testing health endpoint..."
HEALTH_RESPONSE=$(curl --noproxy "*" -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/health)
if [ "$HEALTH_RESPONSE" = "200" ]; then
    echo "   ✅ Health endpoint responds (HTTP $HEALTH_RESPONSE)"
else
    echo "   ❌ Health endpoint failed (HTTP $HEALTH_RESPONSE)"
fi

echo ""

# Test 3: Test login page accessibility (bypass proxy for localhost)
echo "3. Testing login page accessibility..."
LOGIN_RESPONSE=$(curl --noproxy "*" -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/reader/login)
if [ "$LOGIN_RESPONSE" = "200" ]; then
    echo "   ✅ Login page is accessible (HTTP $LOGIN_RESPONSE)"
else
    echo "   ❌ Login page failed (HTTP $LOGIN_RESPONSE)"
fi

echo ""

# Test 4: Check if we can get the HTML content (bypass proxy for localhost)
echo "4. Checking HTML content..."
HTML_CONTENT=$(curl --noproxy "*" -s http://127.0.0.1:8080/reader/login | head -20)
if echo "$HTML_CONTENT" | grep -q "Login"; then
    echo "   ✅ Login page contains expected content"
else
    echo "   ⚠️  Login page content may be incomplete"
fi

echo ""

# Test 5: Test with localhost (for comparison, bypass proxy)
echo "5. Testing with localhost (for comparison)..."
LOCALHOST_RESPONSE=$(curl --noproxy "*" -s -o /dev/null -w "%{http_code}" http://localhost:8080/reader/login)
if [ "$LOCALHOST_RESPONSE" = "200" ]; then
    echo "   ✅ localhost also works (HTTP $LOCALHOST_RESPONSE)"
else
    echo "   ⚠️  localhost returned HTTP $LOCALHOST_RESPONSE"
fi

echo ""
echo "=== Test Complete ==="
echo ""
echo "⚠️  IMPORTANT: You have a proxy configured that may interfere with Safari!"
echo "   Proxy: $http_proxy"
echo ""
echo "To fix Safari access:"
echo "  1. Safari → Preferences → Advanced → Proxies"
echo "  2. Uncheck 'Web Proxy (HTTP)' OR"
echo "  3. Add '127.0.0.1' and 'localhost' to 'Bypass proxy settings for these Hosts & Domains'"
echo ""
echo "To test in Safari:"
echo "  1. Open Safari"
echo "  2. Go to: http://127.0.0.1:8080/reader/login"
echo "  3. Open Developer Console (Develop → Show JavaScript Console)"
echo "  4. Check for any errors"
echo ""
echo "To clear Safari cache/cookies:"
echo "  Safari → Preferences → Privacy → Manage Website Data → Remove All"
echo "  Or: Safari → Clear History → All History"
echo ""
echo "Quick test command (bypasses proxy):"
echo "  curl --noproxy '*' http://127.0.0.1:8080/health"

