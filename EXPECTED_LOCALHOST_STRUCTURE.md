# Expected Localhost Database Structure

This document shows what the database structure **should be** on both localhost and Render.

## ✅ Localhost Structure (Verified)

### Sections Table Structure
```sql
CREATE TABLE "sections" (
  id TEXT PRIMARY KEY,
  page_id TEXT NOT NULL,
  page_number INTEGER NOT NULL, -- Denormalized for quick lookup
  section_number INTEGER NOT NULL, -- Section number within the page (1, 2, or 3)
  content TEXT NOT NULL,
  word_count INTEGER, -- Word count for this section
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE,
  UNIQUE(page_id, section_number)
)
```

### Sections Table Columns
| Column | Type | Nullable | Primary Key |
|--------|------|----------|-------------|
| id | TEXT | NO | YES |
| page_id | TEXT | NO | NO |
| page_number | INTEGER | NO | NO |
| section_number | INTEGER | NO | NO |
| content | TEXT | NO | NO |
| word_count | INTEGER | YES | NO |
| created_at | TEXT | YES | NO |

### Data Counts (Localhost)
- **Total sections:** 77
- **Sections for page 1:** 5
- **Total pages:** 17
- **Total tables:** 21

### Sections per Page (First 10 Pages)
| Page | Sections |
|------|----------|
| 1 | 5 |
| 2 | 3 |
| 3 | 5 |
| 4 | 5 |
| 5 | 5 |
| 6 | 2 |
| 7 | 5 |
| 8 | 5 |
| 9 | 5 |
| 10 | 5 |

## ✅ Render Should Match

**Render database structure should be IDENTICAL to localhost.**

### Key Indicators of Correct Structure:
1. ✅ Sections table has `page_number` and `section_number` columns
2. ✅ Page 1 has 5 sections (not 1)
3. ✅ Total sections = 77 (not less)
4. ✅ No `sections_new` table (migration completed)

### Key Indicators of Wrong Structure:
1. ❌ Sections table missing `page_number` or `section_number`
2. ❌ Page 1 has only 1 section
3. ❌ Total sections < 70
4. ❌ `sections_new` table still exists

## How to Check on Render

### Option 1: Use the comparison tool
```bash
./bin/compare-db-structure
```

### Option 2: Use the check script
```bash
./check-render-db.sh
```

### Option 3: Manual SQL commands
```bash
# Check sections table structure
sqlite3 data/alice-suite.db "SELECT sql FROM sqlite_master WHERE type='table' AND name='sections';"

# Check sections count
sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections WHERE page_number = 1;"

# Check sections per page
sqlite3 data/alice-suite.db "SELECT page_number, COUNT(*) FROM sections GROUP BY page_number ORDER BY page_number LIMIT 10;"
```

## If Render Doesn't Match

1. **Check Render logs** for migration errors
2. **Run fix-render:**
   ```bash
   ./bin/fix-render
   ```
3. **Check if migrations ran:**
   ```bash
   ./bin/migrate
   ```
4. **Restart the service** to ensure start.sh runs
