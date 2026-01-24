# Fix Render Sections Issue

## Problem
On Render, pages are not showing multiple sections - only the first line/snippet is shown, unlike localhost which shows all sections correctly.

## Root Cause
The sections data might not be imported on Render, or the `fix-render` script isn't running properly.

## Diagnostic Steps

### Step 1: Check Render Logs
After deploying, check the Render logs for:
1. **Migration execution:**
   ```
   ✅ Migration 003_restructure_pages_and_sections.sql completed
   ```

2. **fix-render execution:**
   ```
   Verifying and fixing sections data...
   ✅ Page 1 now has 5 sections (expected 5+)
   ```

3. **Section count in logs:**
   ```
   Found X sections in database for page 1
   ```

### Step 2: Use Render Shell to Diagnose

1. **Open Render Shell:**
   - Go to Render Dashboard → Your Service → Shell

2. **Run diagnostic tool:**
   ```bash
   ./bin/diagnose-sections
   ```
   
   This will show:
   - Total sections in database
   - Sections per page
   - Page 1 details
   - Whether API query works

3. **Check sections manually:**
   ```bash
   sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections WHERE page_number = 1;"
   ```
   Should return: **5** (or more)

4. **Check if fix-render ran:**
   ```bash
   sqlite3 data/alice-suite.db "SELECT COUNT(*) FROM sections;"
   ```
   Should return: **77** (or more)

### Step 3: Manual Fix (if needed)

If sections are missing, run fix-render manually:

```bash
./bin/fix-render
```

This will:
- Check current state
- Import sections data if needed
- Verify page 1 has 5+ sections

### Step 4: Verify API Response

Test the API endpoint directly:

```bash
curl -X POST https://your-app.onrender.com/rest/v1/rpc/get_sections_for_page \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"book_id":"alice-in-wonderland","page_number":1}'
```

Should return JSON with `sections` array containing 5 items.

## Common Issues

### Issue 1: fix-render not running
**Symptom:** Logs show "fix-render binary not found"
**Fix:** Check that `bin/fix-render` exists in the build output

### Issue 2: Sections data not embedded
**Symptom:** fix-render runs but says "No sections data found"
**Fix:** Verify `cmd/fix-render/sections-data.sql` exists and is committed to git

### Issue 3: Sections imported but not returned
**Symptom:** Database has sections but API returns empty array
**Fix:** Check Render logs for query errors in `handleGetSectionsForPage`

### Issue 4: Database resets on deploy
**Symptom:** Sections work after fix-render but disappear after redeploy
**Fix:** This is expected on Render free tier (ephemeral filesystem). Ensure `start.sh` runs fix-render on every start.

## Expected Behavior

### Localhost (Working):
- Page 1 shows 5 section snippets in left column
- Each snippet is clickable
- Clicking a snippet shows full section content in middle column
- All sections have content

### Render (Should match):
- Same as localhost
- 5 section snippets visible
- All sections clickable
- Full content displayed when section selected

## Verification Checklist

After deploying to Render:

- [ ] Check Render logs for "fix-render" execution
- [ ] Run `./bin/diagnose-sections` in Render Shell
- [ ] Verify page 1 has 5+ sections: `SELECT COUNT(*) FROM sections WHERE page_number = 1;`
- [ ] Test API endpoint returns sections array
- [ ] Check browser console for errors
- [ ] Verify sections appear in left column
- [ ] Click each section to verify content displays

## Quick Fix Command

If sections are missing, run this in Render Shell:

```bash
./bin/fix-render
```

Then restart the service or wait for next request (sections should be available immediately).
