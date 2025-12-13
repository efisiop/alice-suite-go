# Database Architecture Execution Plan by Claude
## Reader/Consultant Application - SQLite Implementation

---

## Executive Summary

This document provides a complete database schema design for a dual-interface application:
- **Reader Interface**: 100-1000+ concurrent readers with isolated accounts
- **Consultant Interface**: 10-20 consultants monitoring all reader activity in real-time

**Technology Stack**: Go + SQLite  
**Scalability Path**: Designed for easy PostgreSQL migration

---

## Table of Contents

1. [Database Schema Overview](#database-schema-overview)
2. [Table Definitions](#table-definitions)
3. [Indexes for Performance](#indexes-for-performance)
4. [Go Implementation Guidelines](#go-implementation-guidelines)
5. [Security Implementation](#security-implementation)
6. [Real-Time Monitoring Strategy](#real-time-monitoring-strategy)
7. [Scalability Considerations](#scalability-considerations)
8. [Migration Path to PostgreSQL](#migration-path-to-postgresql)

---

## Database Schema Overview

```
┌─────────────────┐
│     users       │ (Readers + Consultants)
└────────┬────────┘
         │
         ├──────────────────┬────────────────────┬─────────────────
         │                  │                    │
┌────────▼────────┐ ┌───────▼──────────┐ ┌──────▼─────────────┐
│ reader_profiles │ │ reader_sessions  │ │ reader_activities  │
└─────────────────┘ └──────────────────┘ └────────────────────┘
                                              
         ┌────────────────────────┐
         │ consultant_access_log  │
         └────────────────────────┘
```

---

## Table Definitions

### 1. Users Table (Core Authentication)

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_role VARCHAR(20) NOT NULL CHECK(user_role IN ('reader', 'consultant', 'admin')),
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    email_verified BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
);

-- Trigger to auto-update updated_at
CREATE TRIGGER update_users_timestamp 
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

**Rationale**: 
- Single table for all users simplifies authentication
- `user_role` enables role-based access control
- `is_active` allows soft account suspension
- `deleted_at` for soft deletes (audit trail)

---

### 2. Reader Profiles Table (Extended Reader Data)

```sql
CREATE TABLE reader_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER UNIQUE NOT NULL,
    reading_level VARCHAR(50),
    preferred_language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    last_login DATETIME,
    total_reading_time_minutes INTEGER DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    avatar_url VARCHAR(500),
    preferences_json TEXT, -- Store flexible preferences as JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER update_reader_profiles_timestamp 
AFTER UPDATE ON reader_profiles
FOR EACH ROW
BEGIN
    UPDATE reader_profiles SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

**Rationale**:
- Separates reader-specific data from core authentication
- `preferences_json` allows flexible feature expansion
- Aggregate fields (`total_reading_time_minutes`) for quick dashboard stats

---

### 3. Reader Sessions Table (Active Session Tracking)

```sql
CREATE TABLE reader_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_heartbeat DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME DEFAULT NULL,
    is_active BOOLEAN DEFAULT 1,
    device_type VARCHAR(50), -- 'mobile', 'tablet', 'desktop'
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**Rationale**:
- `last_heartbeat` enables real-time "who's online" tracking
- Consultants query sessions with `is_active = 1 AND last_heartbeat > datetime('now', '-5 minutes')`
- `ended_at` distinguishes active from historical sessions
- Session tokens for secure authentication

---

### 4. Reader Activities Table (Interaction Logging)

```sql
CREATE TABLE reader_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_id INTEGER,
    activity_type VARCHAR(50) NOT NULL, -- 'page_read', 'bookmark_added', 'annotation', 'quiz_completed', etc.
    content_id VARCHAR(255), -- Reference to what was read (book ID, chapter ID, etc.)
    content_title VARCHAR(500),
    progress_percentage INTEGER, -- Reading progress (0-100)
    duration_seconds INTEGER, -- Time spent on activity
    metadata_json TEXT, -- Flexible JSON for activity-specific data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES reader_sessions(id) ON DELETE SET NULL
);
```

**Rationale**:
- Captures all reader interactions for consultant monitoring
- `activity_type` allows flexible event tracking
- `metadata_json` for activity-specific details without schema changes
- `session_id` links activities to sessions
- High-volume table: needs aggressive indexing

---

### 5. Consultant Access Log Table (Audit Trail)

```sql
CREATE TABLE consultant_access_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    consultant_user_id INTEGER NOT NULL,
    action_type VARCHAR(100) NOT NULL, -- 'viewed_reader_profile', 'exported_report', 'modified_settings', etc.
    target_user_id INTEGER, -- Which reader was accessed (if applicable)
    action_details TEXT,
    ip_address VARCHAR(45),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (consultant_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (target_user_id) REFERENCES users(id) ON DELETE SET NULL
);
```

**Rationale**:
- Audit trail for compliance and security
- Tracks which consultant accessed which reader's data
- Enables accountability and debugging

---

## Indexes for Performance

```sql
-- Users table indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(user_role);
CREATE INDEX idx_users_active ON users(is_active, deleted_at);

-- Reader profiles indexes
CREATE INDEX idx_reader_profiles_user_id ON reader_profiles(user_id);
CREATE INDEX idx_reader_profiles_last_login ON reader_profiles(last_login);

-- Reader sessions indexes (CRITICAL for real-time monitoring)
CREATE INDEX idx_sessions_user_id ON reader_sessions(user_id);
CREATE INDEX idx_sessions_active ON reader_sessions(is_active, last_heartbeat);
CREATE INDEX idx_sessions_token ON reader_sessions(session_token);

-- Reader activities indexes (MOST CRITICAL - high volume table)
CREATE INDEX idx_activities_user_id ON reader_activities(user_id);
CREATE INDEX idx_activities_session_id ON reader_activities(session_id);
CREATE INDEX idx_activities_type ON reader_activities(activity_type);
CREATE INDEX idx_activities_created_at ON reader_activities(created_at);
CREATE INDEX idx_activities_content_id ON reader_activities(content_id);
-- Composite index for consultant dashboard queries
CREATE INDEX idx_activities_user_type_created ON reader_activities(user_id, activity_type, created_at);

-- Consultant access log indexes
CREATE INDEX idx_consultant_log_consultant_id ON consultant_access_log(consultant_user_id);
CREATE INDEX idx_consultant_log_target_id ON consultant_access_log(target_user_id);
CREATE INDEX idx_consultant_log_created_at ON consultant_access_log(created_at);
```

**Rationale**:
- Indexes on foreign keys speed up joins
- Composite indexes optimize common query patterns
- `last_heartbeat` index crucial for "active users" queries
- `created_at` indexes enable efficient time-range queries

---

## Go Implementation Guidelines

### 1. Database Connection Setup

```go
package database

import (
    "database/sql"
    "log"
    "sync"
    
    _ "github.com/mattn/go-sqlite3"
)

var (
    db   *sql.DB
    once sync.Once
)

// InitDB initializes the database connection with connection pooling
func InitDB(dbPath string) (*sql.DB, error) {
    var err error
    once.Do(func() {
        db, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_timeout=5000")
        if err != nil {
            log.Fatal(err)
            return
        }
        
        // SQLite connection pool settings (CRITICAL for concurrency)
        db.SetMaxOpenConns(25) // Limit concurrent connections
        db.SetMaxIdleConns(5)
        db.SetConnMaxLifetime(0) // Reuse connections indefinitely
        
        // Enable Write-Ahead Logging for better concurrent reads
        _, err = db.Exec("PRAGMA journal_mode=WAL;")
        if err != nil {
            log.Fatal(err)
        }
        
        // Foreign key constraints
        _, err = db.Exec("PRAGMA foreign_keys=ON;")
        if err != nil {
            log.Fatal(err)
        }
    })
    
    return db, err
}
```

**Key Settings**:
- **WAL mode**: Allows concurrent reads while writing (ESSENTIAL)
- **Connection pooling**: Prevents database lock contention
- **Timeout**: 5s prevents indefinite locks

---

### 2. User Model Example

```go
package models

import (
    "database/sql"
    "time"
)

type User struct {
    ID            int64      `json:"id"`
    Email         string     `json:"email"`
    PasswordHash  string     `json:"-"` // Never serialize password
    UserRole      string     `json:"user_role"`
    FullName      string     `json:"full_name"`
    IsActive      bool       `json:"is_active"`
    EmailVerified bool       `json:"email_verified"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
    DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, userID int64) (*User, error) {
    var user User
    err := db.QueryRow(`
        SELECT id, email, password_hash, user_role, full_name, 
               is_active, email_verified, created_at, updated_at, deleted_at
        FROM users 
        WHERE id = ? AND deleted_at IS NULL
    `, userID).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.UserRole, 
        &user.FullName, &user.IsActive, &user.EmailVerified, 
        &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
    )
    
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Always filter by deleted_at IS NULL for soft deletes
```

---

### 3. Session Management Example

```go
package models

import (
    "database/sql"
    "time"
)

type ReaderSession struct {
    ID            int64     `json:"id"`
    UserID        int64     `json:"user_id"`
    SessionToken  string    `json:"session_token"`
    IPAddress     string    `json:"ip_address"`
    UserAgent     string    `json:"user_agent"`
    StartedAt     time.Time `json:"started_at"`
    LastHeartbeat time.Time `json:"last_heartbeat"`
    EndedAt       *time.Time `json:"ended_at,omitempty"`
    IsActive      bool      `json:"is_active"`
    DeviceType    string    `json:"device_type"`
}

// UpdateHeartbeat keeps session alive
func UpdateHeartbeat(db *sql.DB, sessionToken string) error {
    _, err := db.Exec(`
        UPDATE reader_sessions 
        SET last_heartbeat = CURRENT_TIMESTAMP 
        WHERE session_token = ? AND is_active = 1
    `, sessionToken)
    return err
}

// GetActiveSessions returns all active sessions (for consultant dashboard)
func GetActiveSessions(db *sql.DB, minutesThreshold int) ([]ReaderSession, error) {
    rows, err := db.Query(`
        SELECT s.id, s.user_id, s.session_token, s.ip_address, 
               s.started_at, s.last_heartbeat, s.device_type, u.full_name
        FROM reader_sessions s
        JOIN users u ON s.user_id = u.id
        WHERE s.is_active = 1 
        AND s.last_heartbeat > datetime('now', '-' || ? || ' minutes')
        AND u.user_role = 'reader'
        ORDER BY s.last_heartbeat DESC
    `, minutesThreshold)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var sessions []ReaderSession
    for rows.Next() {
        var s ReaderSession
        var userName string
        err := rows.Scan(&s.ID, &s.UserID, &s.SessionToken, &s.IPAddress,
                        &s.StartedAt, &s.LastHeartbeat, &s.DeviceType, &userName)
        if err != nil {
            continue
        }
        sessions = append(sessions, s)
    }
    
    return sessions, nil
}
```

---

### 4. Activity Logging Example

```go
package models

import (
    "database/sql"
    "encoding/json"
    "time"
)

type ReaderActivity struct {
    ID                 int64           `json:"id"`
    UserID             int64           `json:"user_id"`
    SessionID          *int64          `json:"session_id,omitempty"`
    ActivityType       string          `json:"activity_type"`
    ContentID          string          `json:"content_id"`
    ContentTitle       string          `json:"content_title"`
    ProgressPercentage *int            `json:"progress_percentage,omitempty"`
    DurationSeconds    *int            `json:"duration_seconds,omitempty"`
    Metadata           json.RawMessage `json:"metadata,omitempty"`
    CreatedAt          time.Time       `json:"created_at"`
}

// LogActivity records a reader activity
func LogActivity(db *sql.DB, activity *ReaderActivity) error {
    _, err := db.Exec(`
        INSERT INTO reader_activities 
        (user_id, session_id, activity_type, content_id, content_title, 
         progress_percentage, duration_seconds, metadata_json)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, activity.UserID, activity.SessionID, activity.ActivityType, 
       activity.ContentID, activity.ContentTitle, activity.ProgressPercentage,
       activity.DurationSeconds, activity.Metadata)
    
    return err
}

// GetReaderActivities retrieves activities for a specific reader
func GetReaderActivities(db *sql.DB, userID int64, limit int) ([]ReaderActivity, error) {
    rows, err := db.Query(`
        SELECT id, user_id, session_id, activity_type, content_id, 
               content_title, progress_percentage, duration_seconds, 
               metadata_json, created_at
        FROM reader_activities
        WHERE user_id = ?
        ORDER BY created_at DESC
        LIMIT ?
    `, userID, limit)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var activities []ReaderActivity
    for rows.Next() {
        var a ReaderActivity
        err := rows.Scan(&a.ID, &a.UserID, &a.SessionID, &a.ActivityType,
                        &a.ContentID, &a.ContentTitle, &a.ProgressPercentage,
                        &a.DurationSeconds, &a.Metadata, &a.CreatedAt)
        if err != nil {
            continue
        }
        activities = append(activities, a)
    }
    
    return activities, nil
}
```

---

## Security Implementation

### 1. Middleware for Role-Based Access Control

```go
package middleware

import (
    "net/http"
    "yourapp/models"
)

// RequireRole middleware checks user role
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract user from context (set during authentication)
            user, ok := r.Context().Value("user").(*models.User)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            // Check if user role is allowed
            allowed := false
            for _, role := range allowedRoles {
                if user.UserRole == role {
                    allowed = true
                    break
                }
            }
            
            if !allowed {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// EnsureOwnData middleware ensures readers only access their own data
func EnsureOwnData(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, ok := r.Context().Value("user").(*models.User)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Extract requested user ID from URL params
        requestedUserID := extractUserIDFromRequest(r)
        
        // Readers can only access their own data
        if user.UserRole == "reader" && user.ID != requestedUserID {
            http.Error(w, "Forbidden - Access denied", http.StatusForbidden)
            return
        }
        
        // Consultants can access all reader data
        next.ServeHTTP(w, r)
    })
}
```

### 2. Query Isolation Pattern

**Critical Rule**: ALWAYS filter by `user_id` for reader queries

```go
// CORRECT: Reader accessing their own activities
activities, err := GetReaderActivities(db, authenticatedUser.ID, 50)

// WRONG: Never allow readers to query all activities
// activities, err := db.Query("SELECT * FROM reader_activities") // DANGEROUS!
```

For consultants, explicitly check role before allowing broad queries:

```go
if user.UserRole != "consultant" {
    return nil, errors.New("unauthorized access")
}
allActivities, err := GetAllReaderActivities(db, limit)
```

---

## Real-Time Monitoring Strategy

### Consultant Dashboard Queries

#### 1. Active Readers Count

```sql
SELECT COUNT(DISTINCT s.user_id) as active_readers
FROM reader_sessions s
WHERE s.is_active = 1 
AND s.last_heartbeat > datetime('now', '-5 minutes');
```

#### 2. Recent Activity Stream (Last 30 minutes)

```sql
SELECT 
    u.full_name,
    a.activity_type,
    a.content_title,
    a.created_at
FROM reader_activities a
JOIN users u ON a.user_id = u.id
WHERE a.created_at > datetime('now', '-30 minutes')
ORDER BY a.created_at DESC
LIMIT 100;
```

#### 3. Reader Progress Summary

```sql
SELECT 
    u.id,
    u.full_name,
    rp.total_reading_time_minutes,
    rp.total_sessions,
    rp.last_login,
    COUNT(DISTINCT a.content_id) as unique_content_read
FROM users u
JOIN reader_profiles rp ON u.id = rp.user_id
LEFT JOIN reader_activities a ON u.id = a.user_id
WHERE u.user_role = 'reader' AND u.is_active = 1
GROUP BY u.id
ORDER BY rp.last_login DESC;
```

### Heartbeat Implementation

Readers should send heartbeat every 60-120 seconds:

```go
// Client-side: Send heartbeat via AJAX/WebSocket
setInterval(() => {
    fetch('/api/heartbeat', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + token }
    });
}, 60000); // Every 60 seconds
```

```go
// Server-side heartbeat endpoint
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
    sessionToken := extractTokenFromRequest(r)
    err := models.UpdateHeartbeat(db, sessionToken)
    if err != nil {
        http.Error(w, "Failed to update heartbeat", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

---

## Scalability Considerations

### SQLite Limitations (100-1000 Users)

| Concurrent Users | SQLite Performance | Recommendation |
|------------------|-------------------|----------------|
| 1-50             | ✅ Excellent      | SQLite is fine |
| 50-200           | ⚠️ Good with WAL  | Monitor performance, optimize queries |
| 200-500          | ⚠️ Marginal       | Consider PostgreSQL migration |
| 500-1000+        | ❌ Poor           | **Must migrate to PostgreSQL** |

### Performance Optimization Techniques

#### 1. Batch Inserts for Activities

Instead of inserting one activity at a time:

```go
// Use transactions for batch inserts
tx, err := db.Begin()
for _, activity := range activities {
    _, err = tx.Exec("INSERT INTO reader_activities (...) VALUES (...)", activity)
}
tx.Commit()
```

#### 2. Prepared Statements

```go
stmt, err := db.Prepare("INSERT INTO reader_activities (...) VALUES (?, ?, ...)")
defer stmt.Close()

for _, activity := range activities {
    _, err = stmt.Exec(activity.UserID, activity.ActivityType, ...)
}
```

#### 3. Read Replicas (Advanced)

For read-heavy consultant queries, consider:
- Separate database file for read-only consultant queries
- Periodic replication from main DB
- Complex setup for SQLite, easier with PostgreSQL

#### 4. Caching Layer

```go
import "github.com/patrickmn/go-cache"

// Cache active sessions for 1 minute
var activeSessionsCache = cache.New(1*time.Minute, 2*time.Minute)

func GetActiveSessionsCached(db *sql.DB) ([]ReaderSession, error) {
    if cached, found := activeSessionsCache.Get("active_sessions"); found {
        return cached.([]ReaderSession), nil
    }
    
    sessions, err := GetActiveSessions(db, 5)
    if err != nil {
        return nil, err
    }
    
    activeSessionsCache.Set("active_sessions", sessions, cache.DefaultExpiration)
    return sessions, nil
}
```

---

## Migration Path to PostgreSQL

When you outgrow SQLite, migrate to PostgreSQL with minimal code changes:

### Schema Translation

SQLite types → PostgreSQL types:

```sql
-- SQLite
INTEGER PRIMARY KEY AUTOINCREMENT → SERIAL PRIMARY KEY
VARCHAR(255) → VARCHAR(255) (same)
TEXT → TEXT (same)
DATETIME → TIMESTAMP WITH TIME ZONE
BOOLEAN → BOOLEAN (same)
```

### PostgreSQL-Specific Optimizations

```sql
-- Partial indexes (more efficient)
CREATE INDEX idx_active_sessions 
ON reader_sessions(user_id, last_heartbeat) 
WHERE is_active = true;

-- JSONB for metadata (better than TEXT)
ALTER TABLE reader_activities 
ALTER COLUMN metadata_json TYPE JSONB USING metadata_json::JSONB;

-- Full-text search capabilities
CREATE INDEX idx_activities_content_search 
ON reader_activities USING GIN(to_tsvector('english', content_title));
```

### Connection Pool Changes

```go
import _ "github.com/lib/pq"

db, err := sql.Open("postgres", 
    "host=localhost port=5432 user=youruser password=yourpass dbname=yourdb sslmode=disable")
db.SetMaxOpenConns(100) // Much higher concurrency support
db.SetMaxIdleConns(20)
```

---

## Database Initialization Script

Create `init_db.go`:

```go
package database

import (
    "database/sql"
    "log"
)

func InitializeSchema(db *sql.DB) error {
    schema := `
    -- Paste all CREATE TABLE statements here
    CREATE TABLE IF NOT EXISTS users (...);
    CREATE TABLE IF NOT EXISTS reader_profiles (...);
    -- ... all other tables
    
    -- Paste all CREATE INDEX statements
    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    -- ... all other indexes
    `
    
    _, err := db.Exec(schema)
    if err != nil {
        log.Fatal("Failed to initialize schema:", err)
        return err
    }
    
    log.Println("Database schema initialized successfully")
    return nil
}
```

---

## Testing Checklist

- [ ] **Concurrency test**: Simulate 100+ concurrent reader logins
- [ ] **Data isolation test**: Verify readers cannot access other readers' data
- [ ] **Session timeout test**: Verify inactive sessions are properly marked
- [ ] **Consultant access test**: Verify consultants can view all reader data
- [ ] **Performance test**: Measure query times under load
- [ ] **Migration test**: Test PostgreSQL migration path

---

## Monitoring and Maintenance

### 1. Database Health Queries

```sql
-- Check database size
SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();

-- Check table sizes
SELECT name, SUM("pgsize") as size 
FROM "dbstat" 
GROUP BY name 
ORDER BY size DESC;

-- Check active connections (requires custom tracking)
SELECT COUNT(*) FROM reader_sessions WHERE is_active = 1;
```

### 2. Cleanup Jobs (Run daily)

```sql
-- Mark stale sessions as inactive (no heartbeat for 1 hour)
UPDATE reader_sessions 
SET is_active = 0, ended_at = CURRENT_TIMESTAMP
WHERE is_active = 1 
AND last_heartbeat < datetime('now', '-1 hour');

-- Archive old activities (older than 1 year)
-- Consider moving to archive table or deleting based on retention policy
```

### 3. Vacuum Database (Weekly)

```sql
VACUUM; -- Reclaim space from deleted rows
ANALYZE; -- Update query planner statistics
```

---

## Emergency Scalability Fixes

If you hit SQLite limits before PostgreSQL migration:

1. **Enable WAL mode** (should already be on)
2. **Increase timeout**: `?_timeout=30000` in connection string
3. **Reduce connection pool**: Lower `MaxOpenConns` to 10-15
4. **Cache aggressively**: Cache consultant dashboard queries for 30-60s
5. **Offload analytics**: Move heavy reporting to separate read-only DB copy
6. **Shard by user ID**: Advanced—split users across multiple SQLite files (complex)

---

## Final Recommendations

1. **Start with SQLite** for proof-of-concept (0-100 users)
2. **Monitor performance** closely as you approach 100-200 concurrent users
3. **Plan PostgreSQL migration** when you hit 200+ concurrent users
4. **Implement caching early** to reduce database load
5. **Use connection pooling** from day one
6. **Test with realistic load** before production launch

---

## Questions or Issues?

If you encounter any issues implementing this architecture:

1. Check SQLite error logs for lock contention
2. Monitor query performance with `EXPLAIN QUERY PLAN`
3. Profile your Go application with `pprof`
4. Consider reaching out for PostgreSQL migration assistance

---

**Document Version**: 1.0  
**Last Updated**: December 2025  
**Author**: Claude (Anthropic AI Assistant)