#!/bin/bash

# Consultant Login Flow Diagnostic Test
# Tests all 24 steps of the consultant login process

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test credentials
EMAIL="consultant@example.com"
PASSWORD="consultant123"
BASE_URL="http://127.0.0.1:8080"

# Counters
PASSED=0
FAILED=0
TOTAL=0

# Function to test a step
test_step() {
    local step_num=$1
    local step_desc=$2
    local test_cmd=$3
    
    TOTAL=$((TOTAL + 1))
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}STEP $step_num: $step_desc${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    if eval "$test_cmd"; then
        echo -e "${GREEN}✓ PASSED${NC}"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# Function to check HTTP response
check_http() {
    local url=$1
    local expected_status=$2
    local method=${3:-GET}
    local data=$4
    local headers=$5
    
    # Use --noproxy to bypass corporate proxy
    if [ -z "$data" ]; then
        if [ -z "$headers" ]; then
            response=$(curl --noproxy "*" -s -w "\n%{http_code}" -X "$method" "$url" 2>&1)
        else
            response=$(curl --noproxy "*" -s -w "\n%{http_code}" -X "$method" -H "$headers" "$url" 2>&1)
        fi
    else
        if [ -z "$headers" ]; then
            response=$(curl --noproxy "*" -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url" 2>&1)
        else
            response=$(curl --noproxy "*" -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -H "$headers" -d "$data" "$url" 2>&1)
        fi
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_status" ]; then
        echo "HTTP Status: $http_code (Expected: $expected_status)"
        echo "Response: ${body:0:200}..."
        return 0
    else
        echo "HTTP Status: $http_code (Expected: $expected_status)"
        echo "Response: $body"
        return 1
    fi
}

echo -e "${YELLOW}═══════════════════════════════════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}  CONSULTANT LOGIN FLOW DIAGNOSTIC TEST${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════════════════════════════════════${NC}"
echo -e "Testing: $BASE_URL/consultant/login"
echo -e "Email: $EMAIL"
echo -e ""

# Initialize token variable
TOKEN=""

# STEP 1-2: Check if login page loads
test_step "1-2" "Server responds to GET /consultant/login" \
    "check_http '$BASE_URL/consultant/login' '200' 'GET'"

# STEP 3-5: Test form submission preparation (check if form exists in HTML)
test_step "3-5" "Login form exists in HTML" \
    "curl --noproxy '*' -s '$BASE_URL/consultant/login' | grep -q 'id=\"login-form\"'"

# STEP 6-7: Test POST to /auth/v1/token endpoint exists
test_step "6-7" "POST /auth/v1/token endpoint accepts requests" \
    "(check_http '$BASE_URL/auth/v1/token' '400' 'POST' '{}' || check_http '$BASE_URL/auth/v1/token' '401' 'POST' '{}')"

# STEP 8-9: Test database query and password verification
test_step "8-9" "Database query and password verification (Login with valid credentials)" \
    "response=\$(curl --noproxy '*' -s -w '\n%{http_code}' -X POST -H 'Content-Type: application/json' -d '{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}' '$BASE_URL/auth/v1/token' 2>&1) && \
     http_code=\$(echo \"\$response\" | tail -n1) && \
     body=\$(echo \"\$response\" | sed '\$d') && \
     if [ \"\$http_code\" = '200' ]; then \
       TOKEN=\$(echo \"\$body\" | grep -o '\"access_token\":\"[^\"]*\"' | cut -d'\"' -f4) && \
       echo \"Token received: \${TOKEN:0:50}...\" && \
       [ -n \"\$TOKEN\" ]; \
     else \
       echo \"HTTP \$http_code: \$body\" && \
       false; \
     fi"

# Extract token for subsequent tests
echo -e "\n${YELLOW}Extracting token for subsequent tests...${NC}"
TOKEN_RESPONSE=$(curl --noproxy "*" -s -X POST -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" "$BASE_URL/auth/v1/token" 2>&1)
TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}ERROR: Could not extract token. Cannot continue with remaining tests.${NC}"
    exit 1
fi

echo -e "${GREEN}Token extracted: ${TOKEN:0:50}...${NC}"

# STEP 10: Test role check (verify token contains consultant role)
test_step "10" "JWT token contains consultant role" \
    "payload=\$(echo \"$TOKEN\" | cut -d'.' -f2 | base64 -d 2>/dev/null || echo \"$TOKEN\" | cut -d'.' -f2 | base64 -D 2>/dev/null) && \
     echo \"\$payload\" | grep -q '\"role\":\"consultant\"' && \
     echo \"Token payload contains consultant role\""

# STEP 11: Test JWT token structure
test_step "11" "JWT token has valid structure (3 parts separated by dots)" \
    "parts=\$(echo \"$TOKEN\" | tr '.' ' ' | wc -w | tr -d ' ') && \
     [ \"\$parts\" = '3' ] && \
     echo \"Token has 3 parts (header.payload.signature)\""

# STEP 12: Test session creation (verify session exists)
test_step "12" "Session created (token is valid)" \
    "check_http '$BASE_URL/auth/v1/user' '200' 'GET' '' 'Authorization: Bearer $TOKEN'"

# STEP 13-14: Test response format
test_step "13-14" "Login response contains required fields" \
    "response=\$(curl --noproxy '*' -s -X POST -H 'Content-Type: application/json' -d '{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}' '$BASE_URL/auth/v1/token') && \
     echo \"\$response\" | grep -q '\"access_token\"' && \
     echo \"\$response\" | grep -q '\"token_type\"' && \
     echo \"\$response\" | grep -q '\"user\"' && \
     echo \"Response contains all required fields\""

# STEP 15: Test token storage capability (simulate localStorage)
test_step "15" "Token can be stored (format is valid for localStorage)" \
    "[ -n \"$TOKEN\" ] && [ \${#TOKEN} -gt 50 ] && echo \"Token length: \${#TOKEN} characters\""

# STEP 16: Test redirect target exists
test_step "16" "Dashboard endpoint exists (/consultant)" \
    "check_http '$BASE_URL/consultant' '302' 'GET' '' 'Authorization: Bearer $TOKEN' || \
     check_http '$BASE_URL/consultant' '200' 'GET' '' 'Authorization: Bearer $TOKEN'"

# STEP 17-18: Test dashboard access with token
test_step "17-18" "Dashboard access with valid token (GET /consultant with Authorization header)" \
    "response=\$(curl --noproxy '*' -s -w '\n%{http_code}' -X GET -H 'Authorization: Bearer $TOKEN' '$BASE_URL/consultant' 2>&1) && \
     http_code=\$(echo \"\$response\" | tail -n1) && \
     body=\$(echo \"\$response\" | sed '\$d') && \
     if [ \"\$http_code\" = '200' ]; then \
       echo \"Dashboard HTML received (length: \${#body} chars)\" && \
       echo \"\$body\" | grep -q 'Consultant Dashboard' && \
       true; \
     elif [ \"\$http_code\" = '302' ]; then \
       echo \"Redirect received (expected for some auth flows)\" && \
       true; \
     else \
       echo \"HTTP \$http_code\" && \
       false; \
     fi"

# STEP 19: Test dashboard HTML template loads
test_step "19" "Dashboard HTML template loads correctly" \
    "response=\$(curl --noproxy '*' -s -X GET -H 'Authorization: Bearer $TOKEN' '$BASE_URL/consultant' 2>&1) && \
     echo \"\$response\" | grep -q 'Consultant Dashboard' && \
     echo \"\$response\" | grep -q 'dashboard.html' || echo \"\$response\" | grep -q 'Logged In' && \
     echo \"Dashboard template loaded\""

# STEP 20-21: Test dashboard JavaScript initialization
test_step "20-21" "Dashboard JavaScript initialization code exists" \
    "response=\$(curl --noproxy '*' -s -X GET -H 'Authorization: Bearer $TOKEN' '$BASE_URL/consultant' 2>&1) && \
     (echo \"\$response\" | grep -q 'initializeDashboard' || \
      echo \"\$response\" | grep -q 'loadLoggedInReaders' || \
      echo \"\$response\" | grep -q 'loadActiveReaders') && \
     echo \"Dashboard JavaScript functions found\""

# STEP 22: Test API endpoints exist
test_step "22" "API endpoints exist and accept requests" \
    "check_http '$BASE_URL/api/consultant/logged-in-readers-count' '200' 'GET' '' 'Authorization: Bearer $TOKEN' || \
     check_http '$BASE_URL/api/consultant/logged-in-readers-count' '401' 'GET' '' && \
     check_http '$BASE_URL/api/consultant/active-readers-count' '200' 'GET' '' 'Authorization: Bearer $TOKEN' || \
     check_http '$BASE_URL/api/consultant/active-readers-count' '401' 'GET' ''"

# STEP 23: Test API endpoints return data
test_step "23" "API endpoints return valid JSON data" \
    "response=\$(curl --noproxy '*' -s -X GET -H 'Authorization: Bearer $TOKEN' '$BASE_URL/api/consultant/logged-in-readers-count' 2>&1) && \
     (echo \"\$response\" | grep -q '\"count\"' || echo \"\$response\" | grep -q 'count') && \
     echo \"API returns JSON with count field\""

# STEP 24: Test full dashboard functionality
test_step "24" "Dashboard displays all required elements" \
    "response=\$(curl --noproxy '*' -s -X GET -H 'Authorization: Bearer $TOKEN' '$BASE_URL/consultant' 2>&1) && \
     echo \"\$response\" | grep -q 'Logged In' && \
     echo \"\$response\" | grep -q 'Active Readers' && \
     echo \"\$response\" | grep -q 'Pending Requests' && \
     echo \"All dashboard stat cards found\""

# Summary
echo -e "\n${YELLOW}═══════════════════════════════════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}  TEST SUMMARY${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════════════════════════════════════${NC}"
echo -e "Total Steps Tested: $TOTAL"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}✓ ALL TESTS PASSED! The consultant login flow is working correctly.${NC}"
    exit 0
else
    echo -e "\n${RED}✗ SOME TESTS FAILED. Please review the errors above.${NC}"
    exit 1
fi

