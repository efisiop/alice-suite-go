-- Migration 003: Restructure to Page-Based System
-- Physical Book Structure: Pages -> Sections (60-65 words each)
-- Based on first edition: 1865 Macmillan & Co., London / 1866 D. Appleton & Co., New York

PRAGMA foreign_keys = ON;

-- Create pages table
CREATE TABLE IF NOT EXISTS pages (
  id TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  page_number INTEGER NOT NULL,
  chapter_id TEXT, -- Chapter that starts on this page (can be NULL if chapter continues)
  chapter_title TEXT, -- Chapter title if it appears on this page
  content TEXT, -- Full page content (for reference)
  word_count INTEGER, -- Approximate word count for this page
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE SET NULL,
  UNIQUE(book_id, page_number)
);

-- Create index for page lookups
CREATE INDEX IF NOT EXISTS idx_pages_book_page ON pages(book_id, page_number);
CREATE INDEX IF NOT EXISTS idx_pages_chapter ON pages(chapter_id);

-- Modify sections table to reference pages instead of chapters
-- First, create new sections table structure
CREATE TABLE IF NOT EXISTS sections_new (
  id TEXT PRIMARY KEY,
  page_id TEXT NOT NULL,
  page_number INTEGER NOT NULL, -- Denormalized for quick lookup
  section_number INTEGER NOT NULL, -- Section number within the page (1, 2, or 3)
  content TEXT NOT NULL,
  word_count INTEGER, -- Word count for this section
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE,
  UNIQUE(page_id, section_number)
);

-- Create index for section lookups
CREATE INDEX IF NOT EXISTS idx_sections_page ON sections_new(page_id);
CREATE INDEX IF NOT EXISTS idx_sections_page_number ON sections_new(page_number, section_number);

-- Note: Old sections table will be dropped after data migration
-- The migration script will handle copying data from old structure to new structure

