-- ============================================================
-- Migration 006: Add Sessions Table and Activity Tracking
-- Purpose: Replace in-memory sessions with database-backed sessions
--          Add comprehensive activity tracking for consultants
-- Created: 2025-01-20
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
-- Note: SQLite doesn't support IF NOT EXISTS for ALTER TABLE ADD COLUMN
-- We'll handle this in application code or check if column exists first
-- For now, we'll add it and handle errors gracefully

-- Check if column exists by attempting to add it (will fail silently if exists)
-- In production, you might want to check schema_version or use a more robust method

