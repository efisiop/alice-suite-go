-- ============================================================
-- Migration 012: Consultant Prompts (AI-style suggestions for readers)
-- Purpose: Let consultants create prompts that suggest the reader
--          get help on a specific page/section; shown to reader as AI-style suggestion
-- ============================================================

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS consultant_prompts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    book_id TEXT NOT NULL,
    page_number INTEGER NOT NULL,
    section_number INTEGER,
    prompt_text TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_consultant_prompts_user_book ON consultant_prompts(user_id, book_id);
CREATE INDEX IF NOT EXISTS idx_consultant_prompts_page ON consultant_prompts(user_id, book_id, page_number);
