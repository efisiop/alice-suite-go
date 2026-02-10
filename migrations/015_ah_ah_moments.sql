-- ============================================================
-- Migration 015: Ah Ah Moments (reader successes/realizations shared with others)
-- Purpose: Readers can record "ah ah moments" (successes, realizations) and
--          share them with other readers of the same book.
-- ============================================================

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS ah_ah_moments (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    book_id TEXT NOT NULL,
    content TEXT NOT NULL,
    page_number INTEGER,
    section_number INTEGER,
    shared INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ah_ah_moments_user_book ON ah_ah_moments(user_id, book_id);
CREATE INDEX IF NOT EXISTS idx_ah_ah_moments_book_shared ON ah_ah_moments(book_id, shared);
CREATE INDEX IF NOT EXISTS idx_ah_ah_moments_created ON ah_ah_moments(created_at DESC);
