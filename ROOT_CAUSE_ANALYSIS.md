# Root Cause Analysis: Render.com vs Localhost Discrepancies

## Executive Summary

All discrepancies stem from **Render.com's ephemeral filesystem** which resets the database on every deploy/restart, causing loss of all runtime data. Additionally, there are code-level issues with glossary lookup case sensitivity.

## Issue-by-Issue Analysis

### Issue A: Online Status Dot Not Showing on Render

**Symptom:** Reader who is online doesn't show colored dot on Render, but works on localhost.

**Root Cause:**
- Online status depends on `sessions` table (`GetOnlineReaderIDs()` queries active sessions)
- When Render.com resets database, all sessions are lost
- Even if user logs in after reset, if database resets again before consultant checks, no sessions exist
- Sessions are created on login (`HandleLogin` calls `database.CreateSession`), but if database resets, they're gone

**Code Location:**
- `internal/database/sessions.go:184-210` - `GetOnlineReaderIDs()`
- `internal/handlers/auth.go:70` - Session creation on login
- `internal/templates/consultant/readers.html:684-722` - Online status display

**Fix Required:**
- Ensure database persists on Render (upgrade to persistent disk OR accept ephemeral nature)
- OR: Use in-memory session tracking as fallback (not ideal for multi-instance)

---

### Issue B: Activity Colors (Gray vs Contrasted)

**Symptom:** Recent activities show in gray on Render, but with contrasted colors on localhost.

**Root Cause:**
- Activity colors depend on `activity_logs` table
- CSS classes: `.activity-recent` (within 1 hour) = colored, `.activity-history` (older) = gray
- When database resets, all activity logs are lost
- Even if new activities are logged, there's no historical context
- The JavaScript function `updateActivityStyling()` compares timestamps - if database resets, all activities appear "recent" but might not have proper styling applied

**Code Location:**
- `internal/templates/consultant/dashboard.html:1392-1419` - `updateActivityStyling()`
- `internal/database/activity.go` - Activity logging
- CSS: `.activity-recent` vs `.activity-history` classes

**Fix Required:**
- Same as Issue A - database persistence
- OR: Ensure activity logging works correctly even after reset (should work, but styling might be off)

---

### Issue C: Help Requests Not Persisting

**Symptom:** Help requests from previous sessions are missing on Render, but persist on localhost.

**Root Cause:**
- Help requests stored in `help_requests` table
- This is **runtime data** (created by users during app usage)
- When Render.com resets database, all help requests are lost
- Migrations only seed **static data** (books, chapters, sections, glossary), not runtime data

**Code Location:**
- `internal/database/queries.go:692-700` - `CreateHelpRequest()`
- `internal/database/queries.go:702-756` - `GetHelpRequests()`
- `migrations/001_initial_schema.sql` - Table creation (no seed data for help_requests)

**Fix Required:**
- Database persistence (same as Issue A)
- **Cannot be fixed with migrations** - help requests are user-generated runtime data

---

### Issue D: New Readers Not Persisting

**Symptom:** New reader accounts created are lost when app restarts on Render, but persist on localhost.

**Root Cause:**
- New users stored in `users` table
- `init-users` only creates **test users** (reader@example.com, consultant@example.com, efisio@efisio.com)
- When users register via `/auth/v1/signup`, they're inserted into `users` table
- When Render.com resets database, all user-created accounts are lost
- Only the seeded test users from `init-users` remain

**Code Location:**
- `pkg/auth/auth.go:32-66` - `Register()` function
- `internal/database/queries.go:16-23` - `CreateUser()`
- `cmd/init-users/main.go` - Only seeds test users

**Fix Required:**
- Database persistence (same as Issue A)
- **Cannot be fixed with migrations** - user accounts are runtime data

---

### Issue E: Glossary Term "waistcoat-pocket" Not Found

**Symptom:** On Render, looking up "waistcoat-pocket" in Page 2 - Section 1 shows "Word not found in glossary", but on localhost it shows the definition.

**Root Cause:**
- The term **DOES exist** in migration 002: `('gloss-2', 'alice-in-wonderland', 'waistcoat-pocket', ...)`
- **TWO different lookup paths:**
  1. **RPC handler** (`internal/handlers/rpc.go:72`) uses: `WHERE term = ?` (case-sensitive)
  2. **DictionaryService** (`internal/services/dictionary_service.go:27`) normalizes to lowercase first
- The RPC handler doesn't normalize case, so if the lookup uses different casing, it fails
- SQLite's `=` operator is case-sensitive for TEXT by default
- The term in database is lowercase: `'waistcoat-pocket'`
- If lookup passes `'Waistcoat-Pocket'` or `'WAISTCOAT-POCKET'`, it won't match

**Code Location:**
- `internal/handlers/rpc.go:60-91` - `handleGetDefinitionWithContext()` - **BUG: No case normalization**
- `internal/services/dictionary_service.go:25-48` - `LookupWord()` - **CORRECT: Normalizes to lowercase**
- `migrations/002_seed_first_3_chapters.sql:189` - Term definition

**Fix Required:**
- Normalize term to lowercase in RPC handler (like DictionaryService does)
- OR: Use case-insensitive comparison: `WHERE LOWER(term) = LOWER(?)`

---

### Issue F: Help Requests/Messages Showing 0

**Symptom:** Reader's "My Help Requests" and "My Messages" show 0 on Render, but show previous data on localhost.

**Root Cause:**
- Same as Issue C - help requests are runtime data lost on database reset
- Messages likely stored in similar runtime tables
- When database resets, all user-generated data is lost

**Code Location:**
- Same as Issue C
- Reader's "My Page" queries help requests for that user

**Fix Required:**
- Same as Issue C - database persistence

---

## Common Root Cause: Database Persistence

**The Fundamental Problem:**
Render.com's **free tier uses ephemeral filesystem** - the database file (`data/alice-suite.db`) is stored in the container's filesystem, which is **reset on every deploy or service restart**.

**What Gets Lost:**
- ✅ **Static data** (books, chapters, sections, glossary) - **RESTORED** by migrations on every start
- ❌ **Runtime data** (sessions, help requests, new users, activity logs) - **LOST FOREVER** on reset

**Why Localhost Works:**
- Local database file persists on disk
- No automatic resets
- All data accumulates over time

---

## Solutions

### Solution 1: Use Render Persistent Disk (Recommended)
**Pros:**
- Database persists across deploys
- All issues resolved
- No code changes needed

**Cons:**
- Requires paid plan (not free tier)
- May have additional costs

**Implementation:**
1. Upgrade Render.com plan to include persistent disk
2. Mount persistent disk to `data/` directory
3. Database will persist across deploys

---

### Solution 2: Accept Ephemeral Nature + Fix Code Issues
**Pros:**
- Works on free tier
- Fixes code bugs (glossary lookup)

**Cons:**
- Runtime data still lost on reset
- Users must re-register
- No historical data

**Implementation:**
1. Fix glossary lookup case sensitivity (Issue E)
2. Document that data is ephemeral
3. Consider in-memory alternatives for critical features

---

### Solution 3: Use External Database (PostgreSQL)
**Pros:**
- Full persistence
- Scalable
- Professional solution

**Cons:**
- Requires database setup
- Code changes needed (SQLite → PostgreSQL)
- Additional service to manage

**Implementation:**
1. Create Render PostgreSQL database
2. Update code to use PostgreSQL driver
3. Migrate schema and data

---

## Immediate Fixes (Code-Level)

### Fix 1: Glossary Lookup Case Sensitivity

**File:** `internal/handlers/rpc.go`

**Current Code (BUG):**
```go
query := `SELECT term, definition, example FROM alice_glossary WHERE term = ? AND book_id = ? LIMIT 1`
```

**Fixed Code:**
```go
// Normalize term to lowercase (like DictionaryService does)
normalizedTerm := strings.ToLower(strings.TrimSpace(term))
query := `SELECT term, definition, example FROM alice_glossary WHERE LOWER(term) = ? AND book_id = ? LIMIT 1`
err := database.DB.QueryRow(query, normalizedTerm, bookID).Scan(&term, &definition, &example)
```

---

## Summary Table

| Issue | Root Cause | Data Type | Fix Type | Priority |
|-------|------------|-----------|----------|----------|
| A: Online dot | Sessions lost | Runtime | Persistence | High |
| B: Activity colors | Activity logs lost | Runtime | Persistence | Medium |
| C: Help requests | Help requests lost | Runtime | Persistence | High |
| D: New readers | Users lost | Runtime | Persistence | High |
| E: Glossary lookup | Case sensitivity bug | Code bug | Code fix | High |
| F: Messages 0 | Help requests lost | Runtime | Persistence | High |

**Key Insight:** Issues A, B, C, D, F all share the same root cause - **database persistence**. Only Issue E is a code bug that can be fixed independently.

