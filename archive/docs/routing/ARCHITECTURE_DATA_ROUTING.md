# Architecture: Data Routing and User Hierarchy Rules

**Version:** 1.0  
**Last Updated:** 2025-01-23  
**Status:** ACTIVE - Must Follow

---

## Table of Contents

1. [Overview](#overview)
2. [User Hierarchy](#user-hierarchy)
3. [Data Routing Rules](#data-routing-rules)
4. [Data Isolation](#data-isolation)
5. [API Endpoint Rules](#api-endpoint-rules)
6. [Database Schema Rules](#database-schema-rules)
7. [Implementation Guidelines](#implementation-guidelines)

---

## Overview

This document codifies the fundamental rules for data routing, user hierarchy, and data isolation in the Alice Suite system. These rules ensure:

- **Clear separation** between reader and consultant interfaces
- **Proper data isolation** - each reader's data is kept separate
- **Correct data flow** from UI to database
- **No data conflicts** between readers and consultants

---

## User Hierarchy

### Role Definitions

#### Reader (`role = "reader"`)
- **Purpose:** End users who read books and interact with content
- **Data Access:** Can only access their own data
- **UI Path:** `/reader/*`
- **Sign Up:** Public signup at `/reader/register`
- **Verification:** Must verify book code after signup
- **Data Scope:** 
  - Own reading progress
  - Own interactions
  - Own help requests
  - Own vocabulary lookups
  - Own AI interactions

#### Consultant (`role = "consultant"`)
- **Purpose:** Monitor and assist readers
- **Data Access:** Can view ALL reader data (read-only monitoring)
- **UI Path:** `/consultant/*`
- **Sign Up:** Admin/manual creation only (no public signup)
- **Verification:** No verification required
- **Data Scope:**
  - View all reader activities
  - View all reader progress
  - View all help requests
  - View all reader interactions
  - Respond to help requests
  - Send prompts to readers

### Hierarchy Rules

```
CONSULTANT (Higher Level)
    ↓ (can view all)
READER (Lower Level)
    ↓ (can only view own)
OWN DATA
```

**Key Principle:** Consultants can see everything, readers can only see their own data.

---

## Data Routing Rules

### Rule 1: UI to API Routing

#### Reader UI (`/reader/*`)
```
Reader UI → API Endpoint → Database Query
    ↓           ↓              ↓
/reader/*   /api/*        WHERE user_id = current_user.id
                         AND role = 'reader'
```

**Example Flow:**
```
/reader/interaction → /api/books → SELECT * FROM books
/reader/interaction → /api/activity/track → INSERT INTO interactions WHERE user_id = ?
```

#### Consultant UI (`/consultant/*`)
```
Consultant UI → API Endpoint → Database Query
    ↓              ↓                ↓
/consultant/*  /api/consultant/*  WHERE role = 'reader'
                                    (all readers)
```

**Example Flow:**
```
/consultant → /api/consultant/reader-activities → SELECT * FROM interactions WHERE user.role = 'reader'
/consultant → /api/consultant/active-readers-count → SELECT DISTINCT user_id FROM interactions WHERE user.role = 'reader'
```

### Rule 2: API Endpoint Naming Convention

#### Reader Endpoints
- **Pattern:** `/api/*` or `/rest/v1/*`
- **Authentication:** Requires reader token
- **Data Filter:** Automatically filters by `user_id` from token
- **Examples:**
  - `/api/books` - Get books (public)
  - `/api/activity/track` - Track activity (requires user_id)
  - `/rest/v1/reading_progress?user_id=...` - Get own progress

#### Consultant Endpoints
- **Pattern:** `/api/consultant/*`
- **Authentication:** Requires consultant token
- **Data Filter:** Returns ALL reader data (no user_id filter)
- **Middleware:** Protected by `middleware.RequireConsultant`
- **Examples:**
  - `/api/consultant/reader-activities` - Get all reader activities
  - `/api/consultant/active-readers-count` - Count active readers
  - `/api/consultant/logged-in-readers-count` - Count logged-in readers

### Rule 3: Data Flow Direction

```
┌─────────────────────────────────────────────────────────────┐
│                    READER DATA FLOW                         │
└─────────────────────────────────────────────────────────────┘

Reader UI (/reader/*)
    │
    ├─→ Track Activity → /api/activity/track
    │                      │
    │                      └─→ INSERT INTO interactions (user_id = reader.id)
    │
    ├─→ Get Progress → /rest/v1/reading_progress?user_id=...
    │                    │
    │                    └─→ SELECT * FROM reading_progress WHERE user_id = ?
    │
    └─→ Create Help Request → /api/help
                                │
                                └─→ INSERT INTO help_requests (user_id = reader.id)


┌─────────────────────────────────────────────────────────────┐
│                 CONSULTANT DATA FLOW                        │
└─────────────────────────────────────────────────────────────┘

Consultant UI (/consultant/*)
    │
    ├─→ View Activities → /api/consultant/reader-activities
    │                       │
    │                       └─→ SELECT * FROM interactions 
    │                           JOIN users ON interactions.user_id = users.id
    │                           WHERE users.role = 'reader'
    │
    ├─→ View Active Readers → /api/consultant/active-readers-count
    │                          │
    │                          └─→ SELECT DISTINCT user_id FROM interactions
    │                              WHERE user.role = 'reader' AND ...
    │
    └─→ Respond to Help → /rest/v1/help_requests (PATCH)
                            │
                            └─→ UPDATE help_requests 
                                SET consultant_id = ?, status = 'assigned'
                                WHERE id = ?
```

---

## Data Isolation

### Rule 4: Reader Data Isolation

**CRITICAL RULE:** Each reader's data MUST be isolated by `user_id`.

#### Database Queries for Readers

**✅ CORRECT:**
```sql
-- Reader accessing own progress
SELECT * FROM reading_progress WHERE user_id = ?

-- Reader accessing own interactions
SELECT * FROM interactions WHERE user_id = ?

-- Reader accessing own help requests
SELECT * FROM help_requests WHERE user_id = ?
```

**❌ WRONG:**
```sql
-- NEVER do this - exposes other readers' data
SELECT * FROM reading_progress

-- NEVER do this - missing user_id filter
SELECT * FROM interactions WHERE book_id = ?
```

#### API Handler Implementation for Readers

**✅ CORRECT:**
```go
// Extract user_id from JWT token
claims, err := auth.ValidateJWT(token)
userID := claims.UserID

// Query with user_id filter
query := `SELECT * FROM interactions WHERE user_id = ?`
rows, err := database.DB.Query(query, userID)
```

**❌ WRONG:**
```go
// NEVER trust user_id from request body
userID := req.UserID // ❌ INSECURE

// NEVER query without user_id filter
query := `SELECT * FROM interactions` // ❌ EXPOSES ALL DATA
```

### Rule 5: Consultant Data Access

**CRITICAL RULE:** Consultants can view ALL reader data, but MUST filter by `role = 'reader'`.

#### Database Queries for Consultants

**✅ CORRECT:**
```sql
-- Consultant viewing all reader activities
SELECT i.*, u.first_name, u.last_name, u.email
FROM interactions i
JOIN users u ON i.user_id = u.id
WHERE u.role = 'reader'
ORDER BY i.created_at DESC

-- Consultant viewing active readers
SELECT DISTINCT u.id, u.first_name, u.last_name, u.email
FROM users u
INNER JOIN interactions i ON i.user_id = u.id
WHERE u.role = 'reader'
AND i.created_at >= ?
```

**❌ WRONG:**
```sql
-- NEVER include consultant data in reader queries
SELECT * FROM interactions -- ❌ Includes consultant activities

-- NEVER forget role filter
SELECT * FROM users -- ❌ Includes consultants
```

#### API Handler Implementation for Consultants

**✅ CORRECT:**
```go
// Verify consultant role
claims, err := auth.ValidateJWT(token)
if claims.Role != "consultant" {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
}

// Query with role filter
query := `
    SELECT i.*, u.first_name, u.last_name, u.email
    FROM interactions i
    JOIN users u ON i.user_id = u.id
    WHERE u.role = 'reader'
`
rows, err := database.DB.Query(query)
```

**❌ WRONG:**
```go
// NEVER skip role verification
// NEVER query without role filter
query := `SELECT * FROM interactions` // ❌ Includes consultant data
```

---

## API Endpoint Rules

### Rule 6: Endpoint Protection

#### Reader Endpoints
```go
// Pattern: Extract user_id from token, filter queries
func HandleReaderEndpoint(w http.ResponseWriter, r *http.Request) {
    // 1. Extract token
    token := extractToken(r)
    
    // 2. Validate token and get user_id
    claims, err := auth.ValidateJWT(token)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // 3. Use claims.UserID in all queries
    userID := claims.UserID
    
    // 4. Query with user_id filter
    query := `SELECT * FROM table WHERE user_id = ?`
    database.DB.Query(query, userID)
}
```

#### Consultant Endpoints
```go
// Pattern: Use RequireConsultant middleware, filter by role
func HandleConsultantEndpoint(w http.ResponseWriter, r *http.Request) {
    // 1. Middleware already verified consultant role
    
    // 2. Query with role filter (no user_id filter)
    query := `
        SELECT i.*, u.first_name, u.last_name, u.email
        FROM interactions i
        JOIN users u ON i.user_id = u.id
        WHERE u.role = 'reader'
    `
    database.DB.Query(query)
}
```

### Rule 7: Route Registration

#### Reader Routes
```go
// In SetupReaderRoutes()
mux.HandleFunc("/reader/*", HandleReaderPage)
mux.HandleFunc("/api/activity/track", HandleTrackActivity) // Uses token user_id
mux.HandleFunc("/rest/v1/reading_progress", HandleReadingProgress) // Filters by user_id
```

#### Consultant Routes
```go
// In SetupAPIRoutes()
mux.Handle("/api/consultant/reader-activities", 
    middleware.RequireConsultant(http.HandlerFunc(HandleGetReaderActivities)))
mux.Handle("/api/consultant/active-readers-count", 
    middleware.RequireConsultant(http.HandlerFunc(HandleGetActiveReadersCount)))
```

---

## Database Schema Rules

### Rule 8: User ID Foreign Keys

**CRITICAL RULE:** All user-specific tables MUST have `user_id` column.

#### Required Tables with `user_id`
- ✅ `interactions` - `user_id` (required)
- ✅ `reading_progress` - `user_id` (required)
- ✅ `help_requests` - `user_id` (required)
- ✅ `vocabulary_lookups` - `user_id` (required)
- ✅ `ai_interactions` - `user_id` (required)
- ✅ `reading_stats` - `user_id` (required)

#### Query Pattern
```sql
-- Always join with users table to get role
SELECT i.*, u.role, u.first_name, u.last_name, u.email
FROM interactions i
JOIN users u ON i.user_id = u.id
WHERE u.role = 'reader'  -- For consultant queries
   OR i.user_id = ?      -- For reader queries
```

### Rule 9: Role-Based Filtering

**CRITICAL RULE:** Always filter by `role` when querying user data.

#### For Reader Queries
```sql
-- Filter by user_id (from token)
WHERE user_id = ? AND role = 'reader'
```

#### For Consultant Queries
```sql
-- Filter by role only (no user_id)
WHERE role = 'reader'
```

---

## Implementation Guidelines

### Rule 10: Code Organization

#### File Structure
```
internal/
├── handlers/
│   ├── reader.go          # Reader UI handlers (/reader/*)
│   ├── consultant.go       # Consultant UI handlers (/consultant/*)
│   ├── api.go              # Public API handlers (/api/*, /rest/v1/*)
│   ├── reader_activity.go  # Consultant reader activity endpoints
│   └── activity.go         # Activity tracking (used by readers)
├── middleware/
│   └── auth.go            # RequireAuth, RequireConsultant, RequireReader
└── database/
    └── queries.go         # Database query helpers
```

#### Handler Naming Convention
- `HandleReader*` - Reader-specific handlers
- `HandleConsultant*` - Consultant-specific handlers
- `HandleGetReaderActivities` - Consultant viewing reader data
- `HandleTrackActivity` - Reader tracking own activity

### Rule 11: Middleware Usage

#### Reader Endpoints
```go
// Use RequireAuth (extracts user_id from token)
mux.Handle("/api/activity/track", middleware.RequireAuth(http.HandlerFunc(HandleTrackActivity)))
```

#### Consultant Endpoints
```go
// Use RequireConsultant (verifies role = 'consultant')
mux.Handle("/api/consultant/reader-activities", 
    middleware.RequireConsultant(http.HandlerFunc(HandleGetReaderActivities)))
```

### Rule 12: Error Handling

#### Reader Endpoints
```go
// If user_id doesn't match, return 403 Forbidden
if claims.UserID != requestedUserID {
    http.Error(w, "Forbidden: Cannot access other user's data", http.StatusForbidden)
    return
}
```

#### Consultant Endpoints
```go
// If not consultant, return 403 Forbidden
if claims.Role != "consultant" {
    http.Error(w, "Forbidden: Consultant access required", http.StatusForbidden)
    return
}
```

---

## Data Flow Examples

### Example 1: Reader Tracks Activity

```
1. Reader clicks "Next Page" in /reader/interaction
   ↓
2. JavaScript calls: POST /api/activity/track
   Headers: Authorization: Bearer <reader_token>
   Body: { event_type: "PAGE_SYNC", page_number: 5, ... }
   ↓
3. Server extracts user_id from token (NOT from body)
   user_id = claims.UserID  // From JWT token
   ↓
4. Server inserts into database:
   INSERT INTO interactions (user_id, event_type, ...)
   VALUES (?, 'PAGE_SYNC', ...)
   WHERE user_id = <from_token>
   ↓
5. Server broadcasts to consultants via SSE
   ↓
6. Consultant dashboard receives update in real-time
```

### Example 2: Consultant Views Reader Activities

```
1. Consultant opens /consultant dashboard
   ↓
2. JavaScript calls: GET /api/consultant/reader-activities
   Headers: Authorization: Bearer <consultant_token>
   ↓
3. Server verifies consultant role:
   if claims.Role != "consultant" → 403 Forbidden
   ↓
4. Server queries database:
   SELECT i.*, u.first_name, u.last_name, u.email
   FROM interactions i
   JOIN users u ON i.user_id = u.id
   WHERE u.role = 'reader'  -- Only readers, not consultants
   ORDER BY i.created_at DESC
   ↓
5. Server returns all reader activities
   ↓
6. Consultant dashboard displays activities
```

### Example 3: Reader Views Own Progress

```
1. Reader opens /reader/dashboard
   ↓
2. JavaScript calls: GET /rest/v1/reading_progress?user_id=...
   Headers: Authorization: Bearer <reader_token>
   ↓
3. Server extracts user_id from token:
   user_id = claims.UserID
   ↓
4. Server ignores user_id from query parameter (security)
   ↓
5. Server queries database:
   SELECT * FROM reading_progress
   WHERE user_id = <from_token>  -- Only own data
   ↓
6. Server returns reader's own progress only
```

---

## Checklist for New Features

When adding new features, ensure:

- [ ] **Reader endpoints** extract `user_id` from token (never from request body)
- [ ] **Reader endpoints** filter queries by `user_id = claims.UserID`
- [ ] **Consultant endpoints** use `middleware.RequireConsultant`
- [ ] **Consultant endpoints** filter queries by `role = 'reader'`
- [ ] **Database tables** have `user_id` foreign key
- [ ] **Database queries** join with `users` table to get role
- [ ] **API routes** follow naming convention (`/api/consultant/*` for consultant)
- [ ] **Error handling** returns 403 Forbidden for unauthorized access
- [ ] **Data isolation** is enforced at database query level
- [ ] **No data leakage** between readers

---

## Violations and Fixes

### Common Violations

#### ❌ Violation 1: Reader endpoint without user_id filter
```go
// WRONG
query := `SELECT * FROM interactions WHERE book_id = ?`
```

**Fix:**
```go
// CORRECT
query := `SELECT * FROM interactions WHERE user_id = ? AND book_id = ?`
```

#### ❌ Violation 2: Consultant endpoint without role filter
```go
// WRONG
query := `SELECT * FROM interactions`
```

**Fix:**
```go
// CORRECT
query := `
    SELECT i.*, u.first_name, u.last_name, u.email
    FROM interactions i
    JOIN users u ON i.user_id = u.id
    WHERE u.role = 'reader'
`
```

#### ❌ Violation 3: Trusting user_id from request body
```go
// WRONG
userID := req.UserID
```

**Fix:**
```go
// CORRECT
claims, _ := auth.ValidateJWT(token)
userID := claims.UserID
```

---

## Summary

### Key Principles

1. **Readers** can only access their own data (filtered by `user_id` from token)
2. **Consultants** can access all reader data (filtered by `role = 'reader'`)
3. **Never trust** `user_id` from request body - always use token
4. **Always filter** by `role` when querying user data
5. **Always join** with `users` table to get role information
6. **Use middleware** to enforce role-based access control

### Enforcement

These rules MUST be followed in:
- All API handlers
- All database queries
- All route registrations
- All middleware implementations

**Violations of these rules are security issues and must be fixed immediately.**

---

**Document Status:** ACTIVE - All developers must follow these rules.

