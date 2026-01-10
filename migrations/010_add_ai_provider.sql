-- Migration: Add provider field to ai_interactions table
-- This tracks which AI provider (gemini/moonshot) was used for each interaction

-- Add provider column to ai_interactions table
-- Note: SQLite doesn't support IF NOT EXISTS for ALTER TABLE ADD COLUMN
-- We'll handle this gracefully in application code if column already exists
ALTER TABLE ai_interactions ADD COLUMN provider TEXT;

-- Index for provider queries (optional, useful for analytics)
CREATE INDEX IF NOT EXISTS idx_ai_interactions_provider ON ai_interactions(provider);
