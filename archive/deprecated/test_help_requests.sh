#!/bin/bash

echo "=== Testing help_requests API directly ==="

# Create a test user
echo "1. Testing with sample data..."
curl -X POST http://localhost:8080/api/help/request \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test_user_123",
    "book_id": "alice-in-wonderland", 
    "content": "I need help understanding this chapter",
    "context": "Chapter 1, first few pages",
    "section_id": null
  }' \
  -v 2>&1 | head -20

echo -e "\n\n=== Test completed ==="
