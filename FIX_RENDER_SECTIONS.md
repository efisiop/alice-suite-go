# Fix Render.com Sections Issue

## Problem
Render.com shows only 1 section per page instead of multiple sections (5+ per page like localhost).

## Solution

### Step 1: Diagnose the Issue on Render.com

1. **Access Render.com Shell/SSH:**
   - Go to Render.com dashboard
   - Find your service (alice-suite-go)
   - Click "Shell" or use SSH to connect

2. **Run Diagnostic Script:**
   ```bash
   cd /opt/render/project/src  # or wherever your app is
   DB_PATH=data/alice-suite.db go run cmd/fix-sections/main.go
   ```

   This will show you:
   - Current table structure (old vs new)
   - How many sections exist
   - What needs to be fixed

### Step 2: Export Data from Localhost (Already Done)

The sections data has been exported to `scripts/sections-data.sql` (77 sections).

### Step 3: Import Data to Render.com

**Option A: Using Render Shell (Recommended)**

1. **Copy the SQL file to Render.com:**
   - Download `scripts/sections-data.sql` from your local machine
   - Upload it to Render.com (via Git, or copy-paste the content)

2. **Connect to Render.com database:**
   ```bash
   # In Render Shell
   cd /opt/render/project/src
   sqlite3 data/alice-suite.db < scripts/sections-data.sql
   ```

**Option B: Using SQLite Directly**

1. **Get database access:**
   - The database is at `data/alice-suite.db` in your Render.com service
   - Use Render Shell to access it

2. **Clear existing sections (if needed):**
   ```sql
   -- BE CAREFUL: This deletes all sections!
   -- Only do this if you're sure you want to replace all data
   DELETE FROM sections;
   ```

3. **Import the data:**
   ```bash
   sqlite3 data/alice-suite.db < scripts/sections-data.sql
   ```

**Option C: Manual SQL Import**

1. **Open the SQL file:** `scripts/sections-data.sql`
2. **Copy all INSERT statements**
3. **In Render Shell, run:**
   ```bash
   sqlite3 data/alice-suite.db
   ```
4. **Paste and execute the INSERT statements**

### Step 4: Verify the Fix

1. **Check the data:**
   ```sql
   -- Count sections for page 1 (should be 5)
   SELECT COUNT(*) FROM sections WHERE page_number = 1;
   
   -- List all sections for page 1
   SELECT page_number, section_number, 
          SUBSTR(content, 1, 50) as preview
   FROM sections 
   WHERE page_number = 1 
   ORDER BY section_number;
   ```

2. **Test the application:**
   - Go to Render.com reader app
   - Navigate to Reading page
   - Page 1 should now show 5 section snippets instead of 1

### Step 5: Restart the Service (if needed)

After importing data, you may need to restart the Render.com service:
- Go to Render.com dashboard
- Click "Manual Deploy" → "Clear build cache & deploy"

## Files Created

1. **`cmd/fix-sections/main.go`** - Diagnostic script to check database structure
2. **`scripts/sections-data.sql`** - Exported sections data (77 sections)
3. **`DATABASE_DIAGNOSTIC.md`** - Detailed diagnostic guide

## Quick Commands Reference

```bash
# Diagnose database
DB_PATH=data/alice-suite.db go run cmd/fix-sections/main.go

# Import sections data
sqlite3 data/alice-suite.db < scripts/sections-data.sql

# Check page 1 sections
sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections WHERE page_number = 1;"
```

## Notes

- The sections data includes 77 sections across multiple pages
- Page 1 should have 5 sections
- All sections reference pages that should exist in the `pages` table
- If you get foreign key errors, make sure the `pages` table has the corresponding page records

## Troubleshooting

**Error: "FOREIGN KEY constraint failed"**
- The `pages` table might be missing data
- Run migrations first: `DB_PATH=data/alice-suite.db ./bin/migrate`

**Error: "table sections does not exist"**
- Run migrations: `DB_PATH=data/alice-suite.db ./bin/migrate`

**Still seeing only 1 section after import**
- Clear browser cache (Ctrl+Shift+R / Cmd+Shift+R)
- Check database: `sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections WHERE page_number = 1;"`
- Verify the import: `sqlite3 data/alice-suite.db "SELECT * FROM sections WHERE page_number = 1;"`

