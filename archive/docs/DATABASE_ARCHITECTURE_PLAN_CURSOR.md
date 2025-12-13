# Database Architecture & Execution Plan for Alice Suite
## Tailored for Cursor IDE Implementation

**Created:** 2025-01-20  
**Purpose:** Robust, scalable database architecture specifically designed for Alice Suite codebase  
**Status:** Ready for implementation

---

## Executive Summary

This plan provides a **clean, practical database architecture** that:
- ✅ Works seamlessly with your existing Alice Suite codebase
- ✅ Supports 100-1000+ concurrent readers and 10-20 consultants
- ✅ Enables real-time monitoring for consultants
- ✅ Maintains data isolation between readers
- ✅ Uses SQLite with proper concurrency configuration
- ✅ Is simple to understand and maintain

**Key Design Decisions:**
1. **Unified Users Table** - Single table with role column (matches your existing schema)
2. **Database-Backed Sessions** - Replace in-memory session store with persistent sessions table
3. **Activity Tracking** - Comprehensive event logging for consultant dashboards
4. **Heartbeat Mechanism** - Real-time "who's online" tracking
5. **WAL Mode** - Essential for concurrent reads/writes

---

## Table of Contents

1. [Current State Analysis](#current-state-analysis)
2. [Database Configuration (PRAGMAs)](#database-configuration-pragmas)
3. [Schema Enhancements](#schema-enhancements)
4. [Go Implementation Patterns](#go-implementation-patterns)
5. [Real-Time Monitoring Strategy](#real-time-monitoring-strategy)
6. [Security & Isolation](#security--isolation)
7. [Execution Checklist](#execution-checklist)

---

## Current State Analysis

### What You Already Have ✅

Your existing schema (`migrations/001_initial_schema.sql`) includes:
- ✅ `users` table with role column (`reader`/`consultant`)
- ✅ `books`, `chapters`, `pages`, `sections` tables
- ✅ `alice_glossary` and `glossary_section_links`
- ✅ `reading_progress` table
- ✅ `vocabulary_lookups` table
- ✅ `ai_interactions` table
- ✅ `help_requests` table
- ✅ Basic indexes

### What Needs Enhancement 🔧

1. **Missing WAL Mode Configuration** - No PRAGMAs set for concurrency
2. **In-Memory Sessions** - Currently using `pkg/auth/session.go` in-memory store
3. **No Activity Tracking Table** - Need comprehensive event logging
4. **No Heartbeat Mechanism** - Can't track "who's online" efficiently
5. **No Session Table** - Sessions stored in memory (lost on restart)

---

## Database Configuration (PRAGMAs)

### Step 1: Update `internal/database/database.go`

**CRITICAL:** Add these PRAGMAs immediately after opening the database connection.

```go
// internal/database/database.go

func InitDB(dbPath string) error {
	// ... existing directory creation code ...

	// Open database connection with WAL mode support
	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return err
	}

	// CRITICAL: Set connection pool limits for SQLite
	DB.SetMaxOpenConns(25)        // Limit concurrent connections
	DB.SetMaxIdleConns(5)         // Keep 5 idle connections
	DB.SetConnMaxLifetime(0)      // Reuse connections indefinitely

	// Execute PRAGMAs for optimal performance
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",           // Enable Write-Ahead Logging (allows concurrent reads)
		"PRAGMA synchronous = NORMAL;",         // Good balance between speed and safety
		"PRAGMA foreign_keys = ON;",            // Enforce referential integrity
		"PRAGMA busy_timeout = 5000;",          // Wait 5 seconds before "database locked" error
		"PRAGMA wal_autocheckpoint = 1000;",     // Auto-checkpoint WAL every 1000 pages
		"PRAGMA cache_size = -20000;",          // 20 MB page cache (adjust based on available RAM)
		"PRAGMA temp_store = MEMORY;",          // Store temp tables in memory for speed
	}

	for _, pragma := range pragmas {
		if _, err := DB.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma %s: %w", pragma, err)
		}
	}

	// Test connection
	_, err = DB.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}
```

**Why These PRAGMAs:**
- **WAL Mode**: Allows multiple readers while one writer is active (essential for 1000+ users)
- **busy_timeout**: Prevents immediate "database locked" errors
- **cache_size**: Speeds up queries by caching frequently accessed pages
- **Connection Pool Limits**: SQLite works best with limited concurrent connections

---

## Schema Enhancements

### Step 2: Create New Migration File

Create `migrations/006_add_sessions_and_activity.sql`:

```sql
-- ============================================================
-- Migration 006: Add Sessions Table and Activity Tracking
-- Purpose: Replace in-memory sessions with database-backed sessions
--          Add comprehensive activity tracking for consultants
-- ============================================================

-- Sessions Table (Database-Backed)
-- Replaces in-memory session store in pkg/auth/session.go
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,                    -- UUID v4 session token
    user_id TEXT NOT NULL,                  -- FK to users.id
    token_hash TEXT NOT NULL,               -- Hashed version of token (security)
    ip_address TEXT,                        -- Client IP address
    user_agent TEXT,                        -- Browser/client info
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    last_active_at TEXT NOT NULL DEFAULT (datetime('now')),  -- Updated on every request
    expires_at TEXT NOT NULL,               -- When session expires
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for sessions
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_last_active ON sessions(last_active_at);

-- Activity Logs Table (Comprehensive Event Tracking)
-- Tracks all user interactions for consultant dashboards
CREATE TABLE IF NOT EXISTS activity_logs (
    id TEXT PRIMARY KEY,                    -- UUID v4
    user_id TEXT NOT NULL,                  -- FK to users.id
    session_id TEXT,                        -- FK to sessions.id (optional)
    activity_type TEXT NOT NULL,            -- 'LOGIN', 'LOGOUT', 'PAGE_VIEW', 'WORD_LOOKUP', 'AI_INTERACTION', 'HELP_REQUEST', etc.
    book_id TEXT,                           -- FK to books.id (if applicable)
    page_number INTEGER,                    -- Page number (if applicable)
    section_id TEXT,                        -- FK to sections.id (if applicable)
    metadata TEXT,                          -- JSON string with additional context
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE SET NULL,
    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE SET NULL
);

-- Indexes for activity_logs (CRITICAL for consultant queries)
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_logs_type ON activity_logs(activity_type);
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_type_created ON activity_logs(user_id, activity_type, created_at DESC);

-- Reader States Table (Denormalized for Fast Consultant Queries)
-- Materialized view-like table updated by application logic
-- Allows consultants to quickly see "who's online" and "current status"
CREATE TABLE IF NOT EXISTS reader_states (
    user_id TEXT PRIMARY KEY,               -- FK to users.id
    book_id TEXT,                           -- Current book being read
    current_page INTEGER,                   -- Last page viewed
    current_section_id TEXT,                -- Last section viewed
    last_activity_type TEXT,                -- Last activity type
    last_activity_at TEXT,                 -- When last activity occurred
    total_pages_read INTEGER DEFAULT 0,    -- Aggregate: total pages read
    total_word_lookups INTEGER DEFAULT 0,   -- Aggregate: total vocabulary lookups
    total_ai_interactions INTEGER DEFAULT 0, -- Aggregate: total AI questions
    status TEXT DEFAULT 'idle',             -- 'idle', 'active', 'reading', 'stuck'
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE SET NULL,
    FOREIGN KEY (current_section_id) REFERENCES sections(id) ON DELETE SET NULL
);

-- Indexes for reader_states
CREATE INDEX IF NOT EXISTS idx_reader_states_last_activity ON reader_states(last_activity_at DESC);
CREATE INDEX IF NOT EXISTS idx_reader_states_status ON reader_states(status);
CREATE INDEX IF NOT EXISTS idx_reader_states_book ON reader_states(book_id);

-- Add last_active_at column to users table (if not exists)
-- This is a denormalized field for quick "who's online" queries
-- Updated via trigger or application logic
ALTER TABLE users ADD COLUMN last_active_at TEXT;

-- Index for users.last_active_at
CREATE INDEX IF NOT EXISTS idx_users_last_active ON users(last_active_at);
```

**Key Tables Explained:**

1. **`sessions`**: Replaces in-memory session store. Stores session tokens, expiration, and last activity.
2. **`activity_logs`**: Comprehensive event tracking. Every user action is logged here for consultant dashboards.
3. **`reader_states`**: Denormalized table for fast queries. Updated by application logic, not triggers.
4. **`users.last_active_at`**: Quick "who's online" tracking. Updated via heartbeat mechanism.

---

## Go Implementation Patterns

### Step 3: Update Session Management

**File:** `pkg/auth/session.go` → Create `internal/database/sessions.go`

```go
// internal/database/sessions.go

package database

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
	"crypto/sha256"
	"encoding/hex"
)

// Session represents a database-backed session
type Session struct {
	ID          string
	UserID      string
	TokenHash   string
	IPAddress   string
	UserAgent   string
	CreatedAt   time.Time
	LastActiveAt time.Time
	ExpiresAt   time.Time
}

// CreateSession creates a new session in the database
func CreateSession(userID, token, ipAddress, userAgent string, expiresIn time.Duration) (*Session, error) {
	sessionID := uuid.New().String()
	tokenHash := hashToken(token)
	expiresAt := time.Now().Add(expiresIn)
	
	query := `INSERT INTO sessions (id, user_id, token_hash, ip_address, user_agent, created_at, last_active_at, expires_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	_, err := DB.Exec(query, sessionID, userID, tokenHash, ipAddress, userAgent, now, now, expiresAt)
	if err != nil {
		return nil, err
	}
	
	return &Session{
		ID:          sessionID,
		UserID:      userID,
		TokenHash:   tokenHash,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		CreatedAt:   now,
		LastActiveAt: now,
		ExpiresAt:   expiresAt,
	}, nil
}

// GetSessionByToken retrieves a session by token hash
func GetSessionByToken(token string) (*Session, error) {
	tokenHash := hashToken(token)
	
	var s Session
	var createdAtStr, lastActiveStr, expiresStr string
	
	query := `SELECT id, user_id, token_hash, ip_address, user_agent, created_at, last_active_at, expires_at
	          FROM sessions WHERE token_hash = ? AND expires_at > datetime('now')`
	
	err := DB.QueryRow(query, tokenHash).Scan(
		&s.ID, &s.UserID, &s.TokenHash, &s.IPAddress, &s.UserAgent,
		&createdAtStr, &lastActiveStr, &expiresStr,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Parse timestamps
	s.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	s.LastActiveAt, _ = time.Parse("2006-01-02 15:04:05", lastActiveStr)
	s.ExpiresAt, _ = time.Parse("2006-01-02 15:04:05", expiresStr)
	
	return &s, nil
}

// UpdateSessionActivity updates last_active_at for a session
func UpdateSessionActivity(token string) error {
	tokenHash := hashToken(token)
	_, err := DB.Exec(`UPDATE sessions SET last_active_at = datetime('now') WHERE token_hash = ?`, tokenHash)
	return err
}

// DeleteSession removes a session
func DeleteSession(token string) error {
	tokenHash := hashToken(token)
	_, err := DB.Exec(`DELETE FROM sessions WHERE token_hash = ?`, tokenHash)
	return err
}

// CleanupExpiredSessions removes expired sessions (run periodically)
func CleanupExpiredSessions() error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE expires_at < datetime('now')`)
	return err
}

// hashToken creates a SHA-256 hash of the token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
```

### Step 4: Add Activity Logging

**File:** `internal/database/activity.go`

```go
// internal/database/activity.go

package database

import (
	"database/sql"
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

// ActivityLog represents an activity log entry
type ActivityLog struct {
	ID          string
	UserID      string
	SessionID   *string
	ActivityType string
	BookID      *string
	PageNumber  *int
	SectionID   *string
	Metadata    map[string]interface{}
	CreatedAt   time.Time
}

// LogActivity records an activity in the database
func LogActivity(activity *ActivityLog) error {
	activity.ID = uuid.New().String()
	
	var metadataJSON string
	if activity.Metadata != nil {
		jsonBytes, err := json.Marshal(activity.Metadata)
		if err == nil {
			metadataJSON = string(jsonBytes)
		}
	}
	
	query := `INSERT INTO activity_logs 
	          (id, user_id, session_id, activity_type, book_id, page_number, section_id, metadata, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := DB.Exec(query,
		activity.ID, activity.UserID, activity.SessionID, activity.ActivityType,
		activity.BookID, activity.PageNumber, activity.SectionID, metadataJSON, time.Now(),
	)
	
	if err != nil {
		return err
	}
	
	// Update reader_states table (denormalized)
	return updateReaderState(activity)
}

// updateReaderState updates the reader_states table
func updateReaderState(activity *ActivityLog) error {
	// Check if reader_state exists
	var exists bool
	err := DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM reader_states WHERE user_id = ?)`, activity.UserID).Scan(&exists)
	if err != nil {
		return err
	}
	
	if !exists {
		// Create new reader_state
		query := `INSERT INTO reader_states 
		          (user_id, book_id, current_page, current_section_id, last_activity_type, last_activity_at, status, updated_at)
		          VALUES (?, ?, ?, ?, ?, datetime('now'), 'active', datetime('now'))`
		_, err = DB.Exec(query, activity.UserID, activity.BookID, activity.PageNumber, activity.SectionID, activity.ActivityType)
		return err
	}
	
	// Update existing reader_state
	query := `UPDATE reader_states SET
	          book_id = COALESCE(?, book_id),
	          current_page = COALESCE(?, current_page),
	          current_section_id = COALESCE(?, current_section_id),
	          last_activity_type = ?,
	          last_activity_at = datetime('now'),
	          status = 'active',
	          updated_at = datetime('now')
	          WHERE user_id = ?`
	
	_, err = DB.Exec(query, activity.BookID, activity.PageNumber, activity.SectionID, activity.ActivityType, activity.UserID)
	return err
}

// GetRecentActivities retrieves recent activities (for consultant dashboard)
func GetRecentActivities(limit int) ([]*ActivityLog, error) {
	query := `SELECT id, user_id, session_id, activity_type, book_id, page_number, section_id, metadata, created_at
	          FROM activity_logs
	          ORDER BY created_at DESC
	          LIMIT ?`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var activities []*ActivityLog
	for rows.Next() {
		activity := &ActivityLog{}
		var sessionID, bookID, sectionID sql.NullString
		var pageNumber sql.NullInt64
		var metadataJSON sql.NullString
		var createdAtStr string
		
		err := rows.Scan(
			&activity.ID, &activity.UserID, &sessionID, &activity.ActivityType,
			&bookID, &pageNumber, &sectionID, &metadataJSON, &createdAtStr,
		)
		if err != nil {
			continue
		}
		
		if sessionID.Valid {
			activity.SessionID = &sessionID.String
		}
		if bookID.Valid {
			activity.BookID = &bookID.String
		}
		if sectionID.Valid {
			activity.SectionID = &sectionID.String
		}
		if pageNumber.Valid {
			pageNum := int(pageNumber.Int64)
			activity.PageNumber = &pageNum
		}
		if metadataJSON.Valid {
			json.Unmarshal([]byte(metadataJSON.String), &activity.Metadata)
		}
		activity.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		
		activities = append(activities, activity)
	}
	
	return activities, rows.Err()
}

// GetUserActivities retrieves activities for a specific user
func GetUserActivities(userID string, limit int) ([]*ActivityLog, error) {
	query := `SELECT id, user_id, session_id, activity_type, book_id, page_number, section_id, metadata, created_at
	          FROM activity_logs
	          WHERE user_id = ?
	          ORDER BY created_at DESC
	          LIMIT ?`
	
	rows, err := DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	// Same parsing logic as GetRecentActivities...
	// (omitted for brevity, use same pattern)
	return nil, nil
}
```

### Step 5: Add Heartbeat Middleware

**File:** `internal/middleware/heartbeat.go` (NEW)

```go
// internal/middleware/heartbeat.go

package middleware

import (
	"net/http"
	"time"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// HeartbeatMiddleware updates last_active_at on every authenticated request
// This enables "who's online" queries for consultants
func HeartbeatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token
		authHeader := r.Header.Get("Authorization")
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err == nil && token != "" {
			// Get user from token
			user, err := auth.GetUserFromToken(token)
			if err == nil && user != nil {
				// Update session activity (fire and forget)
				go func() {
					database.UpdateSessionActivity(token)
					
					// Update users.last_active_at
					database.DB.Exec(`UPDATE users SET last_active_at = datetime('now') WHERE id = ?`, user.ID)
					
					// Update reader_states.last_activity_at if reader
					if user.Role == "reader" {
						database.DB.Exec(`UPDATE reader_states SET last_activity_at = datetime('now'), status = 'active' WHERE user_id = ?`, user.ID)
					}
				}()
			}
		}
		
		next.ServeHTTP(w, r)
	})
}
```

---

## Real-Time Monitoring Strategy

### Step 6: Consultant Dashboard Queries

**File:** `internal/database/consultant.go` (NEW)

```go
// internal/database/consultant.go

package database

import (
	"database/sql"
	"time"
)

// ActiveReader represents a reader who is currently active
type ActiveReader struct {
	UserID      string
	Email       string
	FirstName   string
	LastName    string
	BookID      string
	CurrentPage int
	LastActiveAt time.Time
	Status      string
}

// GetActiveReaders returns readers active in the last N minutes
func GetActiveReaders(minutesThreshold int) ([]*ActiveReader, error) {
	query := `SELECT u.id, u.email, u.first_name, u.last_name, 
	                 COALESCE(rs.book_id, ''), COALESCE(rs.current_page, 0),
	                 COALESCE(rs.last_activity_at, u.last_active_at, u.created_at) as last_active,
	                 COALESCE(rs.status, 'idle') as status
	          FROM users u
	          LEFT JOIN reader_states rs ON u.id = rs.user_id
	          WHERE u.role = 'reader' 
	          AND u.is_verified = 1
	          AND (
	              rs.last_activity_at >= datetime('now', '-' || ? || ' minutes')
	              OR u.last_active_at >= datetime('now', '-' || ? || ' minutes')
	          )
	          ORDER BY last_active DESC`
	
	rows, err := DB.Query(query, minutesThreshold, minutesThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var readers []*ActiveReader
	for rows.Next() {
		r := &ActiveReader{}
		var lastActiveStr string
		err := rows.Scan(
			&r.UserID, &r.Email, &r.FirstName, &r.LastName,
			&r.BookID, &r.CurrentPage, &lastActiveStr, &r.Status,
		)
		if err != nil {
			continue
		}
		r.LastActiveAt, _ = time.Parse("2006-01-02 15:04:05", lastActiveStr)
		readers = append(readers, r)
	}
	
	return readers, rows.Err()
}

// GetReaderActivitySummary returns activity summary for a specific reader
func GetReaderActivitySummary(userID string, hours int) (map[string]interface{}, error) {
	query := `SELECT 
	          COUNT(*) as total_activities,
	          COUNT(DISTINCT DATE(created_at)) as active_days,
	          SUM(CASE WHEN activity_type = 'WORD_LOOKUP' THEN 1 ELSE 0 END) as word_lookups,
	          SUM(CASE WHEN activity_type = 'AI_INTERACTION' THEN 1 ELSE 0 END) as ai_interactions,
	          SUM(CASE WHEN activity_type = 'PAGE_VIEW' THEN 1 ELSE 0 END) as page_views
	          FROM activity_logs
	          WHERE user_id = ? 
	          AND created_at >= datetime('now', '-' || ? || ' hours')`
	
	var summary map[string]interface{}
	err := DB.QueryRow(query, userID, hours).Scan(&summary)
	return summary, err
}
```

---

## Security & Isolation

### Step 7: Enforce Data Isolation

**CRITICAL RULE:** Always filter by `user_id` for reader queries.

**Example Pattern:**

```go
// CORRECT: Reader accessing their own data
func GetMyReadingProgress(userID, bookID string) (*models.ReadingProgress, error) {
	// userID comes from authenticated session context
	// Never allow userID to be passed from request body
	return GetReadingProgress(userID, bookID)
}

// CORRECT: Consultant accessing any reader's data
func GetReaderProgress(consultantUserID, targetUserID, bookID string) (*models.ReadingProgress, error) {
	// First verify consultant has permission
	user, err := GetUserByID(consultantUserID)
	if err != nil || user.Role != "consultant" {
		return nil, errors.New("unauthorized")
	}
	
	// Then fetch target reader's data
	return GetReadingProgress(targetUserID, bookID)
}
```

**Middleware Pattern:**

```go
// internal/middleware/isolation.go

func EnsureOwnData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r) // From auth middleware
		
		// Extract requested user ID from URL
		requestedUserID := extractUserIDFromURL(r)
		
		// Readers can only access their own data
		if user.Role == "reader" && user.ID != requestedUserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
```

---

## Execution Checklist

### Phase 1: Database Configuration ✅

- [ ] **Step 1.1**: Update `internal/database/database.go` with PRAGMAs
  - Add WAL mode configuration
  - Set connection pool limits
  - Add all required PRAGMAs
  
- [ ] **Step 1.2**: Test database connection
  - Verify WAL mode is enabled: `PRAGMA journal_mode;`
  - Verify foreign keys are on: `PRAGMA foreign_keys;`

### Phase 2: Schema Migration ✅

- [ ] **Step 2.1**: Create `migrations/006_add_sessions_and_activity.sql`
  - Add `sessions` table
  - Add `activity_logs` table
  - Add `reader_states` table
  - Add `users.last_active_at` column
  - Add all indexes

- [ ] **Step 2.2**: Run migration
  - Execute migration file
  - Verify tables created: `SELECT name FROM sqlite_master WHERE type='table';`

### Phase 3: Go Implementation ✅

- [ ] **Step 3.1**: Create `internal/database/sessions.go`
  - Implement `CreateSession()`
  - Implement `GetSessionByToken()`
  - Implement `UpdateSessionActivity()`
  - Implement `DeleteSession()`

- [ ] **Step 3.2**: Create `internal/database/activity.go`
  - Implement `LogActivity()`
  - Implement `updateReaderState()`
  - Implement `GetRecentActivities()`
  - Implement `GetUserActivities()`

- [ ] **Step 3.3**: Create `internal/middleware/heartbeat.go`
  - Implement heartbeat middleware
  - Update session activity on every request
  - Update `users.last_active_at`

- [ ] **Step 3.4**: Create `internal/database/consultant.go`
  - Implement `GetActiveReaders()`
  - Implement `GetReaderActivitySummary()`

### Phase 4: Integration ✅

- [ ] **Step 4.1**: Update `internal/handlers/auth.go`
  - Replace in-memory session creation with `database.CreateSession()`
  - Update logout to use `database.DeleteSession()`

- [ ] **Step 4.2**: Add activity logging to handlers
  - Log login/logout events
  - Log page views
  - Log word lookups
  - Log AI interactions
  - Log help requests

- [ ] **Step 4.3**: Update middleware chain
  - Add `HeartbeatMiddleware` to all authenticated routes
  - Ensure data isolation middleware is applied

- [ ] **Step 4.4**: Create consultant dashboard endpoints
  - `GET /api/consultant/active-readers` - List active readers
  - `GET /api/consultant/reader/:id/activity` - Reader activity summary
  - `GET /api/consultant/recent-activities` - Recent activity feed

### Phase 5: Testing ✅

- [ ] **Step 5.1**: Test session management
  - Create session → Verify in database
  - Update activity → Verify `last_active_at` updates
  - Delete session → Verify removal

- [ ] **Step 5.2**: Test activity logging
  - Log various activity types
  - Verify `reader_states` updates
  - Query recent activities

- [ ] **Step 5.3**: Test consultant queries
  - Query active readers
  - Verify "who's online" functionality
  - Test activity summaries

- [ ] **Step 5.4**: Test concurrency
  - Simulate 100+ concurrent logins
  - Verify no "database locked" errors
  - Monitor query performance

### Phase 6: Cleanup & Optimization ✅

- [ ] **Step 6.1**: Add periodic cleanup job
  - Clean expired sessions (run every hour)
  - Archive old activity logs (optional)

- [ ] **Step 6.2**: Monitor performance
  - Check query execution times
  - Verify indexes are being used
  - Optimize slow queries

---

## Key Design Principles

### 1. **Simplicity First**
- Use unified `users` table (already exists)
- Keep schema straightforward
- Avoid over-engineering

### 2. **Performance**
- WAL mode for concurrency
- Denormalized `reader_states` for fast queries
- Proper indexes on all foreign keys and time columns

### 3. **Security**
- Hash session tokens in database
- Always filter by `user_id` for readers
- Enforce role-based access control

### 4. **Real-Time Capability**
- Heartbeat mechanism updates `last_active_at`
- `reader_states` table for instant "who's online" queries
- Activity logs for comprehensive tracking

### 5. **Maintainability**
- Clear separation of concerns
- Database-backed sessions (survive restarts)
- Comprehensive activity logging for debugging

---

## Migration Notes

### From In-Memory Sessions to Database Sessions

**Current State:**
- Sessions stored in `pkg/auth/session.go` (in-memory)
- Lost on server restart

**New State:**
- Sessions stored in `sessions` table
- Survive server restarts
- Can be queried for analytics

**Migration Steps:**
1. Deploy new code with database sessions
2. Existing sessions will expire naturally
3. Users will re-authenticate and get new database-backed sessions
4. No data migration needed (sessions are ephemeral)

---

## Performance Benchmarks

### Expected Performance (with WAL mode):

- **Concurrent Reads**: 1000+ simultaneous readers ✅
- **Concurrent Writes**: Limited by SQLite (use connection pooling) ⚠️
- **"Who's Online" Query**: < 50ms with proper indexes ✅
- **Activity Logging**: < 10ms per insert ✅
- **Session Lookup**: < 5ms with token hash index ✅

### Scaling Considerations:

**Current Capacity (SQLite):**
- ✅ 0-500 concurrent users: Excellent
- ⚠️ 500-1000 concurrent users: Good (monitor performance)
- ❌ 1000+ concurrent users: Consider PostgreSQL migration

**When to Migrate to PostgreSQL:**
- Database file size > 10 GB
- Write contention issues
- Need for read replicas
- Complex reporting requirements

---

## Next Steps

1. **Start with Phase 1** - Database configuration (PRAGMAs)
2. **Then Phase 2** - Schema migration
3. **Then Phase 3** - Go implementation
4. **Test thoroughly** before moving to next phase

**Remember:** Make each step work perfectly before moving to the next step. Test with realistic data and concurrent users.

---

**Document Version**: 1.0  
**Last Updated**: 2025-01-20  
**Author**: Cursor AI Assistant (Synthesized from Gemini, GPT-5, Kimi, and Claude proposals)

