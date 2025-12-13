#!/bin/bash

# Test script to verify cookie-based login works correctly

set -e

BASE_URL="http://127.0.0.1:8080"
EMAIL="consultant@example.com"
PASSWORD="consultant123"

echo "Testing Consultant Login with Cookie..."
echo ""

# Step 1: Login and get token
echo "Step 1: Logging in..."
LOGIN_RESPONSE=$(curl --noproxy "*" -s -c /tmp/cookies.txt -X POST \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" \
  "$BASE_URL/auth/v1/token")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "ERROR: Failed to get token from login response"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✓ Token received: ${TOKEN:0:50}..."
echo ""

# Step 2: Check if cookie was set
echo "Step 2: Checking if cookie was set..."
if grep -q "auth_token" /tmp/cookies.txt; then
    echo "✓ Cookie found in response"
    grep "auth_token" /tmp/cookies.txt
else
    echo "✗ ERROR: Cookie not found in response"
    exit 1
fi
echo ""

# Step 3: Try to access dashboard with cookie
echo "Step 3: Accessing dashboard with cookie..."
DASHBOARD_RESPONSE=$(curl --noproxy "*" -s -b /tmp/cookies.txt -w "\n%{http_code}" \
  "$BASE_URL/consultant" 2>&1)

HTTP_CODE=$(echo "$DASHBOARD_RESPONSE" | tail -n1)
BODY=$(echo "$DASHBOARD_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    echo "✓ Dashboard accessed successfully (HTTP $HTTP_CODE)"
    if echo "$BODY" | grep -q "Consultant Dashboard"; then
        echo "✓ Dashboard HTML contains expected content"
    else
        echo "⚠ Warning: Dashboard HTML might not contain expected content"
    fi
else
    echo "✗ ERROR: Failed to access dashboard (HTTP $HTTP_CODE)"
    echo "Response: ${BODY:0:200}..."
    exit 1
fi
echo ""

# Step 4: Try to access dashboard without cookie (should fail)
echo "Step 4: Testing access without cookie (should redirect)..."
NO_COOKIE_RESPONSE=$(curl --noproxy "*" -s -w "\n%{http_code}" \
  "$BASE_URL/consultant" 2>&1)

NO_COOKIE_CODE=$(echo "$NO_COOKIE_RESPONSE" | tail -n1)

if [ "$NO_COOKIE_CODE" = "302" ] || [ "$NO_COOKIE_CODE" = "200" ]; then
    if echo "$NO_COOKIE_RESPONSE" | grep -q "consultant/login"; then
        echo "✓ Correctly redirects to login when no cookie present"
    else
        echo "⚠ Warning: Might not redirect correctly (HTTP $NO_COOKIE_CODE)"
    fi
else
    echo "⚠ Unexpected response code: $NO_COOKIE_CODE"
fi
echo ""

echo "=========================================="
echo "Cookie login test completed!"
echo "=========================================="

# Cleanup
rm -f /tmp/cookies.txt

