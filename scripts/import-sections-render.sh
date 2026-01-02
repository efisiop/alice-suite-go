#!/bin/bash
# Script to import sections data to Render.com database
# This script exports sections from localhost and provides instructions for Render.com

echo "📦 Exporting sections data from localhost..."
echo ""

# Export sections data
sqlite3 data/alice-suite.db <<EOF > sections-data.sql
.mode insert sections
SELECT * FROM sections ORDER BY page_number, section_number;
EOF

echo "✅ Sections data exported to sections-data.sql"
echo ""
echo "📋 Next steps to import to Render.com:"
echo ""
echo "1. Copy sections-data.sql to your Render.com instance"
echo "2. Connect to Render.com database (via SSH or Render Shell)"
echo "3. Run the import command:"
echo "   sqlite3 data/alice-suite.db < sections-data.sql"
echo ""
echo "OR use the diagnostic script first:"
echo "   go run cmd/fix-sections/main.go"
echo ""

