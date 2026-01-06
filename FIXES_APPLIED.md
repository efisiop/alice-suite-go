# Fixes Applied for Render.com Discrepancies

## Summary

**Root Cause:** Render.com's ephemeral filesystem resets the database on every deploy/restart, losing all runtime data (sessions, help requests, new users, activity logs).

**Code Bug Fixed:** Glossary lookup case sensitivity issue.

## Fixes Applied

### ✅ Fix 1: Glossary Lookup Case Sensitivity (Issue E)

**File:** `internal/handlers/rpc.go`

**Problem:** RPC handler used case-sensitive comparison, causing "waistcoat-pocket" lookups to fail if casing didn't match exactly.

**Solution:** Normalize term to lowercase and use `LOWER()` in SQL query for case-insensitive matching (consistent with `DictionaryService`).

**Code Change:**
```go
// Before: Case-sensitive
query := `SELECT term, definition, example FROM alice_glossary WHERE term = ? AND book_id = ? LIMIT 1`

// After: Case-insensitive
normalizedTerm := strings.ToLower(strings.TrimSpace(term))
query := `SELECT term, definition, example FROM alice_glossary WHERE LOWER(term) = ? AND book_id = ? LIMIT 1`
```

**Status:** ✅ Fixed

---

## Issues Requiring Database Persistence

The following issues **cannot be fixed with code changes alone** - they require database persistence on Render.com:

### Issue A: Online Status Dot
- **Cause:** Sessions table resets, so no online status available
- **Fix:** Use Render persistent disk or external database

### Issue B: Activity Colors
- **Cause:** Activity logs table resets, so no historical data for color comparison
- **Fix:** Use Render persistent disk or external database

### Issue C: Help Requests Not Persisting
- **Cause:** Help requests table resets, losing all user-generated requests
- **Fix:** Use Render persistent disk or external database

### Issue D: New Readers Not Persisting
- **Cause:** Users table resets, losing all user-created accounts
- **Fix:** Use Render persistent disk or external database

### Issue F: Help Requests/Messages Showing 0
- **Cause:** Same as Issue C - runtime data lost
- **Fix:** Use Render persistent disk or external database

---

## Recommended Solutions

### Option 1: Render Persistent Disk (Easiest)
1. Upgrade Render.com plan to include persistent disk
2. Mount persistent disk to `data/` directory
3. Database will persist across deploys
4. **All issues resolved**

### Option 2: External PostgreSQL Database
1. Create Render PostgreSQL database
2. Update code to use PostgreSQL driver
3. Migrate schema
4. **All issues resolved**

### Option 3: Accept Ephemeral Nature
1. Document that data is ephemeral
2. Users must re-register after each deploy
3. No historical data available
4. **Only Issue E fixed (glossary lookup)**

---

## Testing Checklist

After applying fixes:

- [ ] Test glossary lookup with different casings: "waistcoat-pocket", "Waistcoat-Pocket", "WAISTCOAT-POCKET"
- [ ] Verify term is found on Render.com
- [ ] Test online status (requires database persistence)
- [ ] Test activity colors (requires database persistence)
- [ ] Test help requests persistence (requires database persistence)
- [ ] Test new user registration persistence (requires database persistence)

---

## Next Steps

1. **Immediate:** Commit glossary lookup fix
2. **Short-term:** Decide on database persistence solution (persistent disk vs PostgreSQL)
3. **Long-term:** Consider data backup/export strategy if using ephemeral database

---

## Files Modified

- ✅ `internal/handlers/rpc.go` - Fixed glossary lookup case sensitivity

## Files Created

- ✅ `ROOT_CAUSE_ANALYSIS.md` - Detailed analysis of all issues
- ✅ `FIXES_APPLIED.md` - This file

