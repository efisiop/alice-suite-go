-- Dismiss / accept feedback columns (SQLite ensure* added these locally; Postgres needs explicit ALTER)
ALTER TABLE consultant_prompts ADD COLUMN IF NOT EXISTS dismissed_at TEXT;
ALTER TABLE consultant_prompts ADD COLUMN IF NOT EXISTS accepted_at TEXT;
