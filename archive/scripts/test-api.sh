#!/bin/bash

# Alice Suite Go - API Testing Script

BASE_URL="http://localhost:8080"
TOKEN=""

echo "Testing Alice Suite Go API"
echo "=========================="
echo ""

# Test health check
echo "1. Testing health check..."
curl -s "${BASE_URL}/health" | jq '.' || echo "Failed"
echo ""

# Test registration
echo "2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/v1/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123",
    "first_name": "Test",
    "last_name": "User"
  }')
echo "$REGISTER_RESPONSE" | jq '.' || echo "$REGISTER_RESPONSE"
echo ""

# Test login
echo "3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/v1/token" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123"
  }')
echo "$LOGIN_RESPONSE" | jq '.' || echo "$LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token' 2>/dev/null)
if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "Failed to get token"
    exit 1
fi
echo "Token: ${TOKEN:0:20}..."
echo ""

# Test get user
echo "4. Testing get user..."
curl -s "${BASE_URL}/auth/v1/user" \
  -H "Authorization: Bearer ${TOKEN}" | jq '.' || echo "Failed"
echo ""

# Test get books
echo "5. Testing get books..."
curl -s "${BASE_URL}/rest/v1/books" \
  -H "Authorization: Bearer ${TOKEN}" | jq '.' || echo "Failed"
echo ""

# Test help request creation
echo "6. Testing help request creation..."
HELP_RESPONSE=$(curl -s -X POST "${BASE_URL}/rest/v1/help_requests?select=*" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "user_id": "test-user-id",
    "book_id": "alice-in-wonderland",
    "content": "I need help with this passage",
    "status": "pending"
  }')
echo "$HELP_RESPONSE" | jq '.' || echo "$HELP_RESPONSE"
echo ""

# Test get help requests
echo "7. Testing get help requests..."
curl -s "${BASE_URL}/rest/v1/help_requests?status=eq.pending" \
  -H "Authorization: Bearer ${TOKEN}" | jq '.' || echo "Failed"
echo ""

# Test RPC function
echo "8. Testing RPC: get_definition_with_context..."
curl -s -X POST "${BASE_URL}/rest/v1/rpc/get_definition_with_context" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "term": "wonderland",
    "book_id": "alice-in-wonderland"
  }' | jq '.' || echo "Failed"
echo ""

echo "API testing complete!"

