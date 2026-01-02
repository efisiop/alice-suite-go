# Database Diagnostic: Sections Issue on Render.com

## Problem
- **Localhost**: Page 1 has 5 sections, all snippets display correctly
- **Render.com**: Page 1 shows only 1 section/snippet

## Root Cause
The database structure on Render.com is likely different from localhost. Migration 003 creates the new structure but doesn't automatically migrate the data.

## Localhost Database Structure (WORKING)
- `sections` table has NEW structure: `page_id`, `page_number`, `section_number`, `content`
- Page 1 has 5 sections (section_number 1-5)
- Data was populated correctly

## Render.com Database (LIKELY ISSUE)
- `sections` table may still have OLD structure: `chapter_id`, `start_page`, `end_page`, `number`
- OR `sections_new` table exists but is empty
- OR data wasn't properly migrated/seeded

## Diagnostic Queries

### 1. Check Sections Table Structure on Render.com
Run this SQL on Render.com database:

```sql
SELECT sql FROM sqlite_master WHERE type='table' AND name='sections';
```

**Expected (NEW structure - GOOD):**
```sql
CREATE TABLE "sections" (
  id TEXT PRIMARY KEY,
  page_id TEXT NOT NULL,
  page_number INTEGER NOT NULL,
  section_number INTEGER NOT NULL,
  content TEXT NOT NULL,
  ...
)
```

**If you see OLD structure (BAD):**
```sql
CREATE TABLE "sections" (
  id TEXT PRIMARY KEY,
  chapter_id TEXT NOT NULL,
  start_page INTEGER NOT NULL,
  end_page INTEGER NOT NULL,
  number INTEGER NOT NULL,
  ...
)
```

### 2. Check How Many Sections Exist for Page 1
```sql
SELECT COUNT(*) as section_count, page_number 
FROM sections 
WHERE page_number = 1;
```

**Expected:** Should return 5 (or more) sections for page_number = 1

### 3. Check if sections_new Table Exists
```sql
SELECT name FROM sqlite_master 
WHERE type='table' AND name LIKE '%section%';
```

**Expected:** Should see `sections` table with new structure

### 4. Check Section Data for Page 1
```sql
SELECT page_number, section_number, 
       SUBSTR(content, 1, 50) as content_preview,
       LENGTH(content) as content_length
FROM sections 
WHERE page_number = 1 
ORDER BY section_number;
```

**Expected:** Should see multiple rows (5 sections) with section_number 1, 2, 3, 4, 5

## Solutions

### Option A: If sections table has OLD structure
You need to migrate the data. The migration 003 creates `sections_new` but doesn't populate it.

### Option B: If sections table is EMPTY or has wrong data
You need to seed the data. Check if seed data exists or needs to be imported.

### Option C: If sections_new exists but sections is old
You may need to:
1. Drop old sections table
2. Rename sections_new to sections
3. Populate with data

## Next Steps
1. Run diagnostic queries on Render.com database
2. Compare results with localhost (which works)
3. Determine which solution (A, B, or C) applies
4. Run appropriate migration/fix

