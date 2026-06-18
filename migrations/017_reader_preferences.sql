-- Migration 017: Reader preferences
-- Stores reader output/help language separately from auth profile data.

CREATE TABLE IF NOT EXISTS reader_preferences (
  user_id TEXT PRIMARY KEY,
  preferred_language_code TEXT NOT NULL DEFAULT 'en',
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT NOT NULL DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_reader_preferences_language ON reader_preferences(preferred_language_code);

INSERT INTO reader_preferences (user_id, preferred_language_code)
SELECT id, 'en'
FROM users
WHERE role = 'reader'
  AND id NOT IN (SELECT user_id FROM reader_preferences);
