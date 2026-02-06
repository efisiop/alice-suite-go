-- ============================================================
-- Migration 013: Add Administrator Role
-- Purpose: Allow users with role 'administrator' for the Admin dashboard.
--          Optional admin_level for global vs regional (e.g. 'global', 'regional').
-- ============================================================

-- SQLite does not allow altering CHECK constraints; recreate users table with new role.
PRAGMA foreign_keys = OFF;

CREATE TABLE users_new (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  role TEXT CHECK (role IN ('reader', 'consultant', 'administrator')) DEFAULT 'reader',
  is_verified INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  admin_level TEXT
);

INSERT INTO users_new (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
SELECT id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

PRAGMA foreign_keys = ON;
