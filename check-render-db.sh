#!/bin/bash

# Script to check database structure on Render
# Run this in Render Shell to compare with localhost

echo "🔍 Checking Render Database Structure"
echo "============================================================"
echo ""

# Check if database exists
if [ ! -f "data/alice-suite.db" ]; then
    echo "❌ Database file not found at data/alice-suite.db"
    exit 1
fi

echo "✅ Database file exists"
echo ""

# Check sections table structure
echo "📋 Sections Table Structure:"
echo "------------------------------------------------------------"
sqlite3 data/alice-suite.db "SELECT sql FROM sqlite_master WHERE type='table' AND name='sections';"
echo ""

# Check sections table columns
echo "📊 Sections Table Columns:"
echo "------------------------------------------------------------"
sqlite3 data/alice-suite.db "PRAGMA table_info(sections);"
echo ""

# Check data counts
echo "📊 Data Counts:"
echo "------------------------------------------------------------"
echo "Total sections:"
sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections;"
echo ""
echo "Sections for page 1:"
sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections WHERE page_number = 1;"
echo ""
echo "Total pages:"
sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM pages;"
echo ""

# Check sections per page (first 10 pages)
echo "📄 Sections per page (first 10):"
echo "------------------------------------------------------------"
sqlite3 data/alice-suite.db "SELECT page_number, COUNT(*) as section_count FROM sections GROUP BY page_number ORDER BY page_number LIMIT 10;"
echo ""

# Check if sections_new table exists (should not exist if migration completed)
echo "🔍 Checking for sections_new table (should not exist):"
echo "------------------------------------------------------------"
if sqlite3 data/alice-suite.db "SELECT name FROM sqlite_master WHERE type='table' AND name='sections_new';" | grep -q sections_new; then
    echo "⚠️  WARNING: sections_new table still exists (migration may not have completed)"
else
    echo "✅ sections_new table does not exist (migration completed)"
fi
echo ""

# Summary
echo "📋 Summary:"
echo "------------------------------------------------------------"
TOTAL_SECTIONS=$(sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections;")
PAGE1_SECTIONS=$(sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections WHERE page_number = 1;")

echo "Total sections: $TOTAL_SECTIONS (expected: 77)"
echo "Page 1 sections: $PAGE1_SECTIONS (expected: 5)"
echo ""

if [ "$PAGE1_SECTIONS" -ge 5 ] && [ "$TOTAL_SECTIONS" -ge 70 ]; then
    echo "✅ Database structure looks correct!"
else
    echo "⚠️  WARNING: Database structure may be incorrect"
    echo "   Run: ./bin/fix-render"
fi
