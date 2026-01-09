-- Migration 009: Create dictionary cache table for external API lookups
-- This caches definitions from external dictionary APIs to reduce API calls
-- and improve performance

PRAGMA foreign_keys = ON;

-- Dictionary cache table stores definitions from external APIs (e.g., dictionaryapi.dev)
CREATE TABLE IF NOT EXISTS dictionary_cache (
  id TEXT PRIMARY KEY,
  word TEXT NOT NULL,              -- Normalized word (lowercase, trimmed)
  definition TEXT NOT NULL,        -- Main definition
  example TEXT,                    -- Example sentence (if available)
  phonetic TEXT,                   -- Phonetic pronunciation (if available)
  part_of_speech TEXT,             -- Part of speech (noun, verb, etc.)
  source_api TEXT,                 -- Source API name (e.g., 'dictionaryapi.dev')
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  UNIQUE(word)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_dictionary_cache_word ON dictionary_cache(word);
CREATE INDEX IF NOT EXISTS idx_dictionary_cache_created_at ON dictionary_cache(created_at);
