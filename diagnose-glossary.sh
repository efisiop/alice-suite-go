#!/bin/bash

echo "=========================================="
echo "Glossary API Diagnostic Script"
echo "=========================================="
echo ""

# Step 1: Check if database file exists
echo "Step 1: Checking if database file exists..."
DB_PATH="data/alice-suite.db"
if [ -f "$DB_PATH" ]; then
    echo "✅ Database file exists: $DB_PATH"
    echo "   File size: $(ls -lh "$DB_PATH" | awk '{print $5}')"
else
    echo "❌ Database file NOT found: $DB_PATH"
    echo "   You need to run migrations first!"
    echo ""
    echo "   Run: go run cmd/migrate/main.go"
    exit 1
fi

echo ""
echo "Step 2: Checking if sqlite3 is installed..."
if command -v sqlite3 &> /dev/null; then
    echo "✅ sqlite3 is installed"
else
    echo "⚠️  sqlite3 not found - install it to run database checks"
    echo "   macOS: brew install sqlite3"
    echo "   Linux: sudo apt-get install sqlite3"
    echo ""
    echo "   Continuing without sqlite3 checks..."
    exit 0
fi

echo ""
echo "Step 3: Checking if alice_glossary table exists..."
TABLE_EXISTS=$(sqlite3 "$DB_PATH" "SELECT name FROM sqlite_master WHERE type='table' AND name='alice_glossary';" 2>/dev/null)
if [ -n "$TABLE_EXISTS" ]; then
    echo "✅ Table 'alice_glossary' exists"
else
    echo "❌ Table 'alice_glossary' does NOT exist!"
    echo "   You need to run migrations!"
    echo ""
    echo "   Run: go run cmd/migrate/main.go"
    exit 1
fi

echo ""
echo "Step 4: Counting glossary terms in database..."
COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM alice_glossary WHERE book_id='alice-in-wonderland';" 2>/dev/null)
if [ -n "$COUNT" ]; then
    echo "✅ Found $COUNT glossary terms for 'alice-in-wonderland'"
    if [ "$COUNT" -eq 0 ]; then
        echo "   ⚠️  WARNING: No glossary terms found!"
        echo "   You may need to seed the database:"
        echo "   Run: go run cmd/seed/main.go"
    fi
else
    echo "❌ Error counting glossary terms"
fi

echo ""
echo "Step 5: Sample glossary terms (first 5)..."
sqlite3 "$DB_PATH" "SELECT term, definition FROM alice_glossary WHERE book_id='alice-in-wonderland' LIMIT 5;" 2>/dev/null | while IFS='|' read -r term definition; do
    echo "   - $term: ${definition:0:50}..."
done

echo ""
echo "=========================================="
echo "Diagnostic Complete!"
echo "=========================================="
echo ""
echo "If all checks passed, the issue might be:"
echo "1. Server not restarted after code changes"
echo "2. Database connection not initialized in server"
echo "3. Check server logs for actual error messages"
echo ""
echo "Next steps:"
echo "1. Make sure your server is running"
echo "2. Check server terminal for error messages"
echo "3. Try the API endpoint: http://localhost:8080/rest/v1/alice_glossary?book_id=alice-in-wonderland"
