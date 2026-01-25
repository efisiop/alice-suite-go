# Compare Database Structure: Localhost vs Render

## Purpose
This tool helps verify that the database structure on Render matches localhost, which is crucial for ensuring sections display correctly.

## How to Use

### On Localhost:
```bash
./bin/compare-db-structure
```

### On Render:
1. Open Render Shell (Dashboard → Your Service → Shell)
2. Run:
```bash
./bin/compare-db-structure
```

## What It Checks

### 1. Table Structure
- Lists all tables in the database
- Shows column definitions for each table
- Verifies key fields exist

### 2. Sections Table Analysis
- Checks if `sections` table exists
- Verifies it has the correct structure:
  - `page_number` (INTEGER)
  - `section_number` (INTEGER)
  - `page_id` (TEXT)
  - `content` (TEXT)
- Detects if old or new structure is in use

### 3. Data Counts
- Total sections in database
- Sections for page 1 (should be 5+)
- Total pages

### 4. Summary
- Overall health check
- Recommendations if issues found

## Expected Output (Localhost)

```
📊 Found 21 tables
✅ Sections table exists
✅ Pages table exists
   Sections: 77
   Sections for page 1: 5
   Pages: 17
   ✅ Page 1 has correct number of sections (5+)
```

## Expected Output (Render - Should Match)

The Render output should be **identical** to localhost:
- Same table count (21)
- Same sections table structure
- Same data counts (77 sections, 5 for page 1)

## If Structures Don't Match

### Issue 1: Different Table Count
**Symptom:** Render has fewer tables than localhost
**Fix:** Check that all migrations ran successfully in Render logs

### Issue 2: Sections Table Missing Fields
**Symptom:** Sections table doesn't have `page_number` or `section_number`
**Fix:** 
1. Check migration 003_restructure_pages_and_sections.sql ran
2. Run: `./bin/migrate` in Render Shell

### Issue 3: Wrong Data Counts
**Symptom:** Page 1 has less than 5 sections
**Fix:** Run `./bin/fix-render` in Render Shell

### Issue 4: sections_new Table Exists
**Symptom:** Both `sections` and `sections_new` tables exist
**Fix:** This is usually OK - migration creates `sections_new` then renames it. If both exist, the migration might not have completed. Check migration logs.

## Quick Comparison

To quickly compare localhost vs Render:

1. **Run on localhost:**
   ```bash
   ./bin/compare-db-structure > localhost-structure.txt
   ```

2. **Run on Render (in Shell):**
   ```bash
   ./bin/compare-db-structure > render-structure.txt
   ```

3. **Compare the files** (or just compare the Summary sections)

## Key Indicators

✅ **Structure is correct if:**
- Sections table has `page_number` and `section_number`
- Page 1 has 5+ sections
- Total sections = 77

❌ **Structure is wrong if:**
- Sections table missing `page_number`
- Page 1 has only 1 section
- Total sections < 70

## Related Tools

- `./bin/diagnose-sections` - Detailed sections analysis
- `./bin/fix-render` - Fix sections data if missing
- `./bin/migrate` - Run migrations manually
