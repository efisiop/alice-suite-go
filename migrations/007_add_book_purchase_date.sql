-- Migration: Add book purchase date to reading_progress table
-- This allows consultants to track when a reader purchased/obtained a book

-- Add purchase_date column to reading_progress table
-- Note: SQLite doesn't support IF NOT EXISTS for ALTER TABLE ADD COLUMN
-- We'll handle this gracefully in application code if column already exists
ALTER TABLE reading_progress ADD COLUMN purchase_date TEXT;

-- Index for purchase_date queries
CREATE INDEX IF NOT EXISTS idx_reading_progress_purchase_date ON reading_progress(purchase_date);

