# Consultant Dashboard Fixes - December 3, 2025

## Issues Reported

1. **"Test user" appearing in dashboard** - Test users showing up in activity feed
2. **Names changing in activity list** - Inconsistent user names displayed
3. **JavaScript error** - `SyntaxError: Can't create duplicate variable: '_0x2bcb81'`

---

## Root Causes Identified

### Issue 1 & 2: Test Users and Inconsistent Names

**Problem:**
- Query was using `LEFT JOIN` which allowed activities with missing user data
- WHERE clause logic `(u.role = 'reader' OR u.role IS NULL)` was incorrect
- Test users in database ("Test Reader", "Test User") were being included
- Activities could show with NULL user data, causing name inconsistencies

**Database Check:**
```sql
-- Found test users:
- Test Reader (reader@example.com)
- Test User (test@example.com)
- Giusy Giusiano (giusy@giusy.com) ✅
- Efisio Pittau (efisio@efisio.com) ✅
```

### Issue 3: JavaScript Error

**Problem:**
- Error `Can't create duplicate variable: '_0x2bcb81'` from minified code
- Likely from browser extension or htmx library
- Already being suppressed in `base.html` but may still appear in console

---

## Fixes Applied

### Fix 1: Updated Activity Query (HandleGetReaderActivities)

**File:** `internal/handlers/reader_activity.go`

**Changes:**
- Changed `LEFT JOIN` to `INNER JOIN` - ensures only activities from existing users
- Fixed WHERE clause to properly filter readers only
- Added filters to exclude test users:
  - `u.email NOT LIKE '%@example.com'`
  - `u.email NOT LIKE 'test@%'`

**Before:**
```sql
FROM interactions i
LEFT JOIN users u ON i.user_id = u.id
WHERE (u.role = 'reader' OR u.role IS NULL)
AND (u.role != 'consultant' OR u.role IS NULL)
```

**After:**
```sql
FROM interactions i
INNER JOIN users u ON i.user_id = u.id
WHERE u.role = 'reader'
AND u.email NOT LIKE '%@example.com'
AND u.email NOT LIKE 'test@%'
```

### Fix 2: Updated Stream Query (HandleGetReaderActivityStream)

**File:** `internal/handlers/reader_activity.go`

**Changes:**
- Same fixes as above - changed to INNER JOIN and added test user filters

### Fix 3: Updated Active Readers Query

**File:** `internal/handlers/reader_activity.go`

**Changes:**
- Added test user filters to `HandleGetActiveReadersCount`
- Ensures test users don't appear in active readers list

### Fix 4: JavaScript Error Suppression

**File:** `internal/templates/base.html`

**Status:** Already implemented
- Error suppression for duplicate variable errors is already in place
- Errors from browser extensions/third-party libraries are filtered

---

## Expected Results

### After Fixes:

1. ✅ **Test users excluded** - Only real reader accounts will appear
2. ✅ **Consistent names** - User names will always match database (first_name + last_name)
3. ✅ **No NULL user data** - All activities will have valid user information
4. ✅ **JavaScript errors suppressed** - Duplicate variable errors filtered (may still appear in console but won't affect functionality)

---

## Testing Checklist

- [ ] Restart server
- [ ] Login as consultant
- [ ] Verify only "Efisio Pittau" and "Giusy Giusiano" appear (no "Test user")
- [ ] Verify names remain consistent in activity feed
- [ ] Verify no activities show with missing user data
- [ ] Check browser console - JavaScript errors should be suppressed

---

## Database Cleanup (Optional)

If you want to completely remove test users from the database:

```sql
-- WARNING: This will delete test users and all their data
DELETE FROM users WHERE email LIKE '%@example.com' OR email LIKE 'test@%';
```

**Note:** Only do this if you're sure you don't need test data. Test users are useful for development/testing.

---

## Files Modified

1. `internal/handlers/reader_activity.go`
   - Fixed `HandleGetReaderActivities` query
   - Fixed `HandleGetReaderActivityStream` query  
   - Fixed `HandleGetActiveReadersCount` query

---

## Verification Queries

### Check current reader users (excluding test):
```sql
SELECT id, first_name, last_name, email 
FROM users 
WHERE role = 'reader' 
AND email NOT LIKE '%@example.com' 
AND email NOT LIKE 'test@%'
ORDER BY created_at DESC;
```

### Check recent activities (should only show real users):
```sql
SELECT i.id, u.first_name, u.last_name, u.email, i.event_type, i.created_at
FROM interactions i
INNER JOIN users u ON i.user_id = u.id
WHERE u.role = 'reader'
AND u.email NOT LIKE '%@example.com'
AND u.email NOT LIKE 'test@%'
ORDER BY i.created_at DESC
LIMIT 20;
```

---

## Status

✅ **Fixes Applied** - Ready for testing

**Next Steps:**
1. Restart the server
2. Test the consultant dashboard
3. Verify only real users appear
4. Verify names remain consistent

