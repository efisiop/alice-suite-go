#!/bin/bash
# Complete database restructuring script

cd "$(dirname "$0")"

echo "🔄 Complete Database Restructuring"
echo "📖 Physical book structure: Pages -> Sections (60-65 words each)"
echo ""

# Backup existing database
if [ -f "data/alice-suite.db" ]; then
    echo "📦 Backing up existing database..."
    cp data/alice-suite.db data/alice-suite.db.backup
fi

# Remove old database
rm -f data/alice-suite.db

# Step 1: Create schema
echo "📄 Step 1: Creating database schema..."
sqlite3 data/alice-suite.db < migrations/001_initial_schema.sql

# Step 2: Load seed data
echo "🌱 Step 2: Loading seed data..."
go run cmd/seed/main.go

# Step 3: Run restructuring
echo "🔄 Step 3: Restructuring to page-based system..."
go run cmd/restructure/main.go

echo ""
echo "✅ Database restructuring complete!"
echo "📊 Check the output above for summary statistics."

